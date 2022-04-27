package enhancementservices

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/theparanoids/ashirt-server/backend"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"

	sq "github.com/Masterminds/squirrel"
)

func TestServiceWorker(workerData models.ServiceWorker) (string, bool, error) {
	var basicConfig BasicServiceWorkerConfig
	err := json.Unmarshal([]byte(workerData.Config), &basicConfig)
	if err != nil {
		return "", false, err
	}
	worker, err := findAppropriateWorker(basicConfig)
	if err != nil {
		return "", false, err
	}
	if err = worker.Build(workerData.Name, 0, []byte(workerData.Config)); err != nil {
		return "", false, err
	}

	return worker.Test()
}

func RunAllServiceWorkers(db *database.Connection, evidenceID int64) ([]string, []error) {
	return RunSetOfServiceWorkers(db, []string{}, evidenceID)
}

func RunSetOfServiceWorkers(db *database.Connection, serviceNames []string, evidenceID int64) ([]string, []error) {
	knownWorkers, err := getServiceWorkerList(db)

	if err != nil {
		return []string{}, []error{backend.WrapError("Unable to run service worker", backend.UnauthorizedWriteErr(err))}
	}

	skipWorkerAlignment := false
	// If no services are specified, then run all
	if len(serviceNames) == 0 {
		serviceNames = helpers.Map(knownWorkers, getServiceWorkerName)
		skipWorkerAlignment = true
	}

	// workerErrors tracks errors encountered when running workers
	workerErrors := make(chan error, len(serviceNames))

	var workersToRun []models.ServiceWorker
	// If this is set, then we know that we're running all of the known workers. Otherwise, we might
	// get a list of workers that we can't work on.
	if skipWorkerAlignment {
		workersToRun = knownWorkers
	} else {
		workersToRun = make([]models.ServiceWorker, 0, len(knownWorkers))

		for _, requestedWorker := range serviceNames {
			_, foundWorker := helpers.Find(knownWorkers, func(w models.ServiceWorker) bool {
				return w.Name == requestedWorker
			})

			if foundWorker != nil {
				workersToRun = append(workersToRun, *foundWorker)
			} else {
				workerErrors <- fmt.Errorf("No current worker named %v", requestedWorker)
			}
		}
	}

	payload, err := buildProcessPayload(db, evidenceID)
	if err != nil {
		return []string{}, []error{backend.WrapError("Unable to construct worker message", backend.UnauthorizedWriteErr(err))}
	}

	if err = markWorkStarting(db, evidenceID, helpers.Map(workersToRun, getServiceWorkerName)); err != nil {
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

	if err = handler.Build(worker.Name, evidenceID, []byte(worker.Config)); err != nil {
		return err
	}

	if pendingUpdate, err := handler.Process(payload); err != nil {
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
	return nil, fmt.Errorf("no worker matches the provided configuration")
}

func getServiceWorkerList(db *database.Connection) ([]models.ServiceWorker, error) {
	var knownWorkers []models.ServiceWorker
	err := db.Select(&knownWorkers,
		sq.Select("*").
			From("service_workers").
			Where(sq.Eq{"deleted_at": nil}),
	)
	return knownWorkers, err
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

func markWorkStarting(db *database.Connection, evidenceID int64, sources []string) error {
	now := time.Now()
	return db.BatchInsert("evidence_metadata", len(sources), func(row int) map[string]interface{} {
		return map[string]interface{}{
			"evidence_id":      evidenceID,
			"source":           sources[row],
			"status":           "Processing",
			"work_started_at":  now,
			"last_run_message": nil,
		}
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
