package enhancementservices

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ashirt-ops/ashirt-server/backend"
	"github.com/ashirt-ops/ashirt-server/backend/database"
	"github.com/ashirt-ops/ashirt-server/backend/helpers"
	"github.com/ashirt-ops/ashirt-server/backend/logging"
	"github.com/ashirt-ops/ashirt-server/backend/models"
	"github.com/ashirt-ops/ashirt-server/backend/servicetypes/evidencemetadata"

	sq "github.com/Masterminds/squirrel"
)

// TestServiceWorker contacts the indicated worker to verify that it's running
func TestServiceWorker(workerData models.ServiceWorker) ServiceTestResult {
	var basicConfig BasicServiceWorkerConfig
	if err := json.Unmarshal([]byte(workerData.Config), &basicConfig); err != nil {
		return errorTestResultWithMessage(err, "Unable to parse worker configuration")
	}
	worker, err := findAppropriateWorker(basicConfig)
	if err != nil {
		return errorTestResultWithMessage(err, "Unable to find matching worker")
	}
	if err = worker.Build(workerData.Name, []byte(workerData.Config)); err != nil {
		return errorTestResultWithMessage(err, "Unable to prep worker for test")
	}

	return worker.Test()
}

type SendServiceWorkerEventInput struct {
	Logger      logging.Logger
	WorkerNames []string
	Builder     func(db database.ConnectionProxy) ([]interface{}, error)
	EventType   string
}

func SendServiceWorkerEvent(db *database.Connection, input SendServiceWorkerEventInput) {
	var workersToRun []models.ServiceWorker
	var payloads []interface{}
	workerContext := context.Background()

	go func() {
		err := db.WithTx(workerContext, func(tx *database.Transactable) {
			workersToRun, _ = filterWorkers(tx, input.WorkerNames)
			payloads, _ = input.Builder(tx)
		})
		if err != nil {
			input.Logger.Log("msg", "Unable to execute service workers", "error", err.Error())
			return
		}

		var wg sync.WaitGroup
		for _, payload := range payloads {
			// set these aside so that they are preserved in the closure
			payloadCopy := payload
			for _, worker := range workersToRun {
				wg.Add(1)
				workerCopy := worker
				go func() {
					defer wg.Done()
					err := runProcessEvent(db, workerCopy, &payloadCopy)
					logger := logging.With(input.Logger,
						"worker", workerCopy.Name,
						"eventType", input.EventType,
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
		if notifyWorkersRunForTest != nil {
			notifyWorkersRunForTest <- true
		}
	}()
}

// SendEvidenceCreatedEvent starts a specified set of workers for a specified set of evidenceUUIDs
// Note that this process kicks off a number of goroutines.
func SendEvidenceCreatedEvent(db *database.Connection, reqLogger logging.Logger, operationID int64, evidenceUUIDs []string, workerNames []string) error {
	var workersToRun []models.ServiceWorker
	var expandedPayloads []ExpandedNewEvidencePayload
	workerContext := context.Background()

	go func() {
		err := db.WithTx(workerContext, func(tx *database.Transactable) {
			workersToRun, _ = filterWorkers(tx, workerNames)
			expandedPayloads, _ = BatchBuildNewEvidencePayload(workerContext, tx, operationID, evidenceUUIDs)

			markWorkStarting(tx,
				helpers.Map(expandedPayloads, getExpandedPayloadID),
				helpers.Map(workersToRun, getServiceWorkerName))
		})
		if err != nil {
			reqLogger.Log("msg", "Unable to execute service workers", "error", err.Error())
			return
		}

		var wg sync.WaitGroup
		for _, ePayload := range expandedPayloads {
			// set these aside so that they are preserved in the closure
			evidenceID := ePayload.EvidenceID
			payload := ePayload.NewEvidencePayload
			for _, worker := range workersToRun {
				wg.Add(1)
				workerCopy := worker
				go func() {
					defer wg.Done()
					err := runProcessMetadata(db, workerCopy, evidenceID, &payload)
					logger := logging.With(reqLogger,
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
		if notifyWorkersRunForTest != nil {
			notifyWorkersRunForTest <- true
		}
	}()

	return nil
}

func runProcessEvent(db *database.Connection, worker models.ServiceWorker, payload interface{}) error {
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

	return handler.ProcessEvent(payload)
}

func runProcessMetadata(db *database.Connection, worker models.ServiceWorker, evidenceID int64, payload *NewEvidencePayload) error {
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

	if pendingUpdate, err := handler.ProcessMetadata(evidenceID, payload); err != nil {
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

func markWorkStarting(db database.ConnectionProxy, evidenceIDs []int64, sources []string) error {
	type entry struct {
		source string
		id     int64
	}

	// create a set of (source/id) pairs.
	entries := make([]entry, len(sources)*len(evidenceIDs))
	numEvidenceIDs := len(evidenceIDs)
	for i, v := range sources {
		rowOffset := i * numEvidenceIDs
		for j, w := range evidenceIDs {
			entries[rowOffset+j] = entry{
				source: v,
				id:     w,
			}
		}
	}

	now := time.Now()
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

	return evidence, err
}

func getAllEvidenceForOperation(db database.ConnectionProxy, operationID int64) ([]models.Evidence, error) {
	var evidence []models.Evidence

	err := db.Select(&evidence, sq.Select("*").From("evidence").Where(sq.Eq{
		"operation_id": operationID,
	}))

	return evidence, err
}

var notifyWorkersRunForTest chan<- bool = nil

func SetNotifyWorkersRunForTest(notifier chan<- bool) {
	if cap(notifier) == 0 {
		fmt.Println("Capacity for notifier channel is 0. This could easily cause deadlocks")
	}
	notifyWorkersRunForTest = notifier
}
