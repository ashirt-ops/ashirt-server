// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package enhancementservices

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/logging"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/servicetypes/evidencemetadata"

	sq "github.com/Masterminds/squirrel"
)

// TestServiceWorker contacts the indicated worker to verify that it's running
func TestServiceWorker(workerData models.ServiceWorker) ServiceTestResult {
	var basicConfig BasicServiceWorkerConfig
	err := json.Unmarshal([]byte(workerData.Config), &basicConfig)
	if err != nil {
		return ErrorTestResultWithMessage(err, "Unable to parse worker configuration")
	}
	worker, err := findAppropriateWorker(basicConfig)
	if err != nil {
		return ErrorTestResultWithMessage(err, "Unable to find matching worker")
	}
	if err = worker.Build(workerData.Name, []byte(workerData.Config)); err != nil {
		return ErrorTestResultWithMessage(err, "Unable to prep worker for test")
	}

	return worker.Test()
}

func RunServiceWorkerMatrix(ctx context.Context, db *database.Connection, operationID int64, evidenceUUIDs []string, workerNames []string) error {
	var workersToRun []models.ServiceWorker
	var evidence []models.Evidence
	var expandedPayloads []ExpandedPayload

	err := db.WithTx(ctx, func(tx *database.Transactable) {
		workersToRun, _ = filterWorkers(tx, workerNames)
		if len(evidenceUUIDs) == 0 {
			evidence, _ = getAllEvidenceForOperation(tx, operationID)
		} else {
			evidence, _ = filterEvidenceByUUID(tx, operationID, evidenceUUIDs)
		}
		expandedPayloads, _ = batchBuildProcessPayload(tx, helpers.Map(evidence, func(i models.Evidence) int64 {
			return i.ID
		}))

		evidenceIDs := helpers.Map(evidence, func(e models.Evidence) int64 { return e.ID })

		markWorkStarting(tx, evidenceIDs, helpers.Map(workersToRun, getServiceWorkerName))
	})
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, ePayload := range expandedPayloads {
		evidenceID := ePayload.EvidenceID
		payload := ePayload.Payload
		for _, worker := range workersToRun {
			wg.Add(1)
			workerCopy := worker
			go func() {
				defer wg.Done()
				err := runWorker(db, workerCopy, evidenceID, &payload)
				logger := logging.With(logging.ReqLogger(ctx),
					"worker", workerCopy.Name,
					"evidenceID", evidenceID,
				)

				if err != nil {
					logger.Log("msg", "Unable to run worker", "error", err)
				} else {
					logger.Log("msg", "Worker completed")
				}
			}()
		}
	}
	wg.Wait()
	return nil
}

// RunAllServiceWorkers starts _all_ of the currents service workers
func RunAllServiceWorkers(db *database.Connection, evidenceID int64) ([]string, []error) {
	return RunSetOfServiceWorkers(db, []string{}, evidenceID)
}

// RunSetOfServiceWorkers starts the indicated service workers (by name)
func RunSetOfServiceWorkers(db *database.Connection, serviceNames []string, evidenceID int64) ([]string, []error) {
	knownWorkers, err := getServiceWorkerList(db)
	if err != nil {
		return []string{}, []error{backend.WrapError("Unable to run service worker", backend.UnauthorizedWriteErr(err))}
	}

	workersToRun, workerErrors := alignWorkers(serviceNames, knownWorkers)

	payload, err := buildProcessPayload(db, evidenceID)
	if err != nil {
		return []string{}, []error{backend.WrapError("Unable to construct worker message", backend.UnauthorizedWriteErr(err))}
	}

	if err = markWorkStarting(db, []int64{evidenceID}, helpers.Map(workersToRun, getServiceWorkerName)); err != nil {
		return []string{}, []error{backend.WrapError("Unable to run workers", err)}
	}

	var wg sync.WaitGroup
	completedWorkersChan := make(chan string, len(workersToRun))
	for _, worker := range workersToRun {
		wg.Add(1)
		workerCopy := worker
		go func() {
			defer wg.Done()
			err := runWorker(db, workerCopy, evidenceID, payload)
			if err != nil {
				workerErrors <- err
			} else {
				completedWorkersChan <- workerCopy.Name
			}
		}()
	}
	wg.Wait()

	completedWorkers := helpers.ChanToSlice(&completedWorkersChan)
	errors := helpers.ChanToSlice(&workerErrors)
	return completedWorkers, errors
}

func runWorker(db *database.Connection, worker models.ServiceWorker, evidenceID int64, payload *Payload) error {
	var err error
	var basicConfig BasicServiceWorkerConfig
	if err = json.Unmarshal([]byte(worker.Config), &basicConfig); err != nil {
		return err
	}

	var handler ServiceWorker
	if handler, err = findAppropriateWorker(basicConfig); err != nil {
		return err
	}

	if err = handler.Build(worker.Name, []byte(worker.Config)); err != nil {
		return err
	}

	if pendingUpdate, err := handler.Process(evidenceID, payload); err != nil {
		return err
	} else if pendingUpdate != nil { // should always be not-nil
		_, err := upsertWorkerCompleteData(db, *pendingUpdate)
		return err
	}

	return nil
}

func findAppropriateWorker(config BasicServiceWorkerConfig) (ServiceWorker, error) {
	if config.Type == "web" {
		if config.Version == 1 {
			return &webConfigV1Worker{}, nil
		}
	}
	if config.Type == "aws" {
		if config.Version == 1 {
			return &awsConfigV1Worker{}, nil
		}
	}
	return nil, fmt.Errorf("no worker matches the provided configuration")
}

func getServiceWorkerList(db database.ConnectionProxy) ([]models.ServiceWorker, error) {
	var knownWorkers []models.ServiceWorker
	err := db.Select(&knownWorkers,
		sq.Select("*").
			From("service_workers").
			Where(sq.Eq{"deleted_at": nil}),
	)
	return knownWorkers, err
}

// batchBuildProcessPayload builds a payload by getting all of the necessary details in bulk.
// Note: this relies on the ordering of evidenceIDs. No particular order is required as input,
// but the result is ordered by evidenceID, in ASC order.
func batchBuildProcessPayload(db database.ConnectionProxy, evidenceIDs []int64) ([]ExpandedPayload, error) {
	var payloads []ExpandedPayload

	err := db.Select(&payloads, sq.Select(
		"e.id AS id",
		"e.uuid AS uuid",
		"e.content_type",
		"slug AS operation_slug",
		"'process' AS type", // hardcode in the type so we don't have to edit each entry manually
	).
		From("evidence e").
		LeftJoin("operations o ON e.operation_id = o.id").
		Where(sq.Eq{"e.id": evidenceIDs}).
		OrderBy(`e.id`),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to gather evidence data for worker")
	}

	return payloads, nil
}

func buildProcessPayload(db *database.Connection, evidenceID int64) (*Payload, error) {
	var payload Payload

	err := db.Get(&payload, sq.Select(
		"e.uuid AS uuid",
		"e.content_type",
		"slug AS operation_slug",
	).
		From("evidence e").
		LeftJoin("operations o ON e.operation_id = o.id").
		Where(sq.Eq{"e.id": evidenceID}),
	)

	payload.Type = "process"

	if err != nil {
		return nil, fmt.Errorf("unable to gather evidence data for worker")
	}

	return &payload, nil
}

func markWorkStarting(db database.ConnectionProxy, evidenceIDs []int64, sources []string) error {
	now := time.Now()

	type entry struct {
		source string
		id     int64
	}

	entries := make([]entry, len(sources)*len(evidenceIDs))

	numEvidenceIDs := len(evidenceIDs)
	for i, v := range sources {
		for j, w := range evidenceIDs {
			entries[i*numEvidenceIDs+j] = entry{
				source: v,
				id:     w,
			}
		}
	}

	return db.BatchInsert("evidence_metadata", len(entries), func(row int) map[string]interface{} {
		return map[string]interface{}{
			"body":             "",
			"evidence_id":      entries[row].id,
			"source":           entries[row].source,
			"status":           evidencemetadata.StatusProcessing,
			"work_started_at":  now,
			"last_run_message": nil,
		}
		// Note that ON DUPLICATE does not update the body. This helps preserve the last body
		// until the work is complete.
	}, "ON DUPLICATE KEY UPDATE "+
		"status=VALUES(status),"+
		"work_started_at=VALUES(work_started_at),"+
		"last_run_message=VALUES(last_run_message)",
	)
}

func getServiceWorkerName(w models.ServiceWorker) string {
	return w.Name
}

func upsertWorkerCompleteData(db *database.Connection, data models.EvidenceMetadata) (int64, error) {
	return db.Insert("evidence_metadata", map[string]interface{}{
		"evidence_id":      data.EvidenceID,
		"source":           data.Source,
		"body":             data.Body,
		"status":           data.Status,
		"last_run_message": data.LastRunMessage,
		"can_process":      data.CanProcess,
	}, "ON DUPLICATE KEY UPDATE "+
		"body=VALUES(body),"+
		"status=VALUES(status),"+
		"last_run_message=VALUES(last_run_message),"+
		"can_process=VALUES(can_process)",
	)
}

// alignWorkers matches the names of the provided services with the currently active services.
// This will return a list of the found workers, and a channel with any errors that occurrd finding these workers.
// Note: if no serviceNames are provided, then _all_ services are returned
func alignWorkers(serviceNames []string, knownServices []models.ServiceWorker) ([]models.ServiceWorker, chan error) {
	// If no services are specified, then run all
	if len(serviceNames) == 0 {
		return knownServices, make(chan error, len(knownServices))
	}

	// workerErrors tracks errors encountered when running workers
	workerErrors := make(chan error, len(serviceNames))
	workersToRun := make([]models.ServiceWorker, 0, len(knownServices))

	for _, requestedWorker := range serviceNames {
		_, foundWorker := helpers.Find(knownServices, func(w models.ServiceWorker) bool {
			return w.Name == requestedWorker
		})

		if foundWorker != nil {
			workersToRun = append(workersToRun, *foundWorker)
		} else {
			workerErrors <- fmt.Errorf("no current worker named %v", requestedWorker)
		}
	}

	return workersToRun, workerErrors
}

// filterWorkers retrives a list of all workers, compares that with a list of workers and returns
// the intersection of those workers. Note that this ignores any errrors encountered when trying to
// match up the workers. For example, if requesting FastWorker and MediumWorker, and only MediumWorker
// is available, then only MediumWorker (and no error) will be returned.
func filterWorkers(db database.ConnectionProxy, serviceNames []string) ([]models.ServiceWorker, error) {
	knownWorkers, err := getServiceWorkerList(db)
	if err != nil {
		return []models.ServiceWorker{}, backend.WrapError("Unable to find service workers", backend.UnauthorizedWriteErr(err))
	}

	workersToRun, _ := alignWorkers(serviceNames, knownWorkers)
	return workersToRun, nil
}

// filterEvidenceByUUID returns all matching evidence given an operation ID and a list of evidence uuids.
// This ignores any errors regarding mismatched evidence UUIDs between what's present for an operation
// and what's requested.
func filterEvidenceByUUID(db database.ConnectionProxy, operationID int64, evidenceUUIDs []string) ([]models.Evidence, error) {
	var evidence []models.Evidence

	err := db.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{
		"operation_id": operationID,
		"uuid":         evidenceUUIDs,
	}))

	if err != nil {
		return []models.Evidence{}, err
	}
	return evidence, nil
}

func getAllEvidenceForOperation(db database.ConnectionProxy, operationID int64) ([]models.Evidence, error) {
	var evidence []models.Evidence

	err := db.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{
		"operation_id": operationID,
	}))

	if err != nil {
		return []models.Evidence{}, err
	}
	return evidence, nil
}
