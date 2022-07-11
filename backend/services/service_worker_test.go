// Copyright 2022, Yahoo Inc.
// Licensed under the terms of the MIT. See LICENSE file in project root for terms.

package services_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/theparanoids/ashirt-server/backend/database"
	"github.com/theparanoids/ashirt-server/backend/dtos"
	"github.com/theparanoids/ashirt-server/backend/enhancementservices"
	"github.com/theparanoids/ashirt-server/backend/helpers"
	"github.com/theparanoids/ashirt-server/backend/models"
	"github.com/theparanoids/ashirt-server/backend/services"
)

func TestListServiceWorkers(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		tryList := func(u models.User) ([]*dtos.ServiceWorker, error) {
			ctx := contextForUser(u, db)
			return services.ListServiceWorker(ctx, db)
		}

		// verify permissions
		_, err := tryList(UserRon)
		require.Error(t, err) // non-admin

		// verify result
		list, err := tryList(UserDumbledore)
		require.NoError(t, err)
		for _, worker := range list {
			_, match := helpers.Find(seed.ServiceWorkers, func(w models.ServiceWorker) bool {
				return w.ID == worker.ID
			})
			require.NotNil(t, match)
			require.Equal(t, match.Name, worker.Name)
			require.Equal(t, match.DeletedAt != nil, worker.Deleted)
			require.JSONEq(t, match.Config, worker.Config)
		}
	})
}

func TestCreateServiceWorker(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		tryCreate := func(u models.User, input services.CreateServiceWorkerInput) error {
			ctx := contextForUser(u, db)
			return services.CreateServiceWorker(ctx, db, input)
		}

		input := services.CreateServiceWorkerInput{
			Name:   "IsAHorcrux",
			Config: `{"type": "web", "version": 1, "url": "http://test:1234"}`,
		}

		// verify permissions
		require.Error(t, tryCreate(UserRon, input)) // non-admin

		// verify result
		require.NoError(t, tryCreate(UserDumbledore, input))
		svc := getServiceWorkerByName(t, db, input.Name)
		require.JSONEq(t, svc.Config, input.Config)
	})
}

func TestUpdateServiceWorker(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		tryUpdate := func(u models.User, input services.UpdateServiceWorkerInput) error {
			ctx := contextForUser(u, db)
			return services.UpdateServiceWorker(ctx, db, input)
		}

		input := services.UpdateServiceWorkerInput{
			ID:     DemoServiceWorker.ID,
			Name:   "BetterWorker",
			Config: `{"type": "web", "version": 1, "url": "http://test:1234"}`,
		}

		// verify permissions
		require.Error(t, tryUpdate(UserRon, input)) // non-admin

		// verify result
		require.NoError(t, tryUpdate(UserDumbledore, input))
		svc := getServiceWorkerByID(t, db, input.ID)
		require.JSONEq(t, svc.Config, input.Config)
	})
}

func TestDeleteServiceWorker(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		tryDelete := func(u models.User, input services.DeleteServiceWorkerInput) error {
			ctx := contextForUser(u, db)
			return services.DeleteServiceWorker(ctx, db, input)
		}

		input := services.DeleteServiceWorkerInput{
			ID:       DemoServiceWorker.ID,
			DoDelete: true,
		}

		checkNilness := func(shouldBeNil bool) {
			svc := getServiceWorkerByID(t, db, input.ID)
			require.Equal(t, svc.DeletedAt == nil, shouldBeNil)
		}

		// verify pre-conditions
		checkNilness(true)

		// verify permissions
		require.Error(t, tryDelete(UserRon, input)) // non-admin
		checkNilness(true)

		// verify delete result
		require.NoError(t, tryDelete(UserDumbledore, input))
		checkNilness(false)

		// verify restore result
		input.DoDelete = false
		require.NoError(t, tryDelete(UserDumbledore, input))
		checkNilness(true)
	})
}

func TestTestServiceWorker(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, _ TestSeedData) {
		worker := DemoServiceWorker
		wasCalled := false

		testSuccess := buildRequestMock(func(w *httptest.ResponseRecorder) {
			wasCalled = true
			w.WriteHeader(http.StatusOK)
			w.Write(mustJSONMarshal(t, enhancementservices.TestResp{Status: "ok"}))
		})

		testFailure := buildRequestMock(func(w *httptest.ResponseRecorder) {
			wasCalled = true
			w.WriteHeader(http.StatusOK)
			w.Write(mustJSONMarshal(t, enhancementservices.TestResp{Status: "error", Message: helpers.Ptr("oops")}))
		})

		tryTest := func(u models.User, id int64) (*dtos.ServiceWorkerTestOutput, error) {
			ctx := contextForUser(u, db)
			return services.TestServiceWorker(ctx, db, id)
		}

		// verify permissions
		_, err := tryTest(UserRon, worker.ID)
		require.Error(t, err) // non-admin

		// verify result (success)
		enhancementservices.SetWebRequestFunctionForWorker(worker.Name, &testSuccess)
		out, err := tryTest(UserDumbledore, worker.ID)
		require.NoError(t, err)
		require.True(t, wasCalled)
		require.Equal(t, worker.ID, out.ID)
		require.Equal(t, worker.Name, out.Name)
		require.True(t, out.Live)

		// verify result (failure)
		wasCalled = false // reset to un-called
		enhancementservices.SetWebRequestFunctionForWorker(worker.Name, &testFailure)
		out, err = tryTest(UserDumbledore, worker.ID)
		require.NoError(t, err)
		require.True(t, wasCalled)
		require.Equal(t, worker.ID, out.ID)
		require.Equal(t, worker.Name, out.Name)
		require.False(t, out.Live)
	})
}

func TestRunServiceWorker(t *testing.T) {
	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		singleWorker := DemoServiceWorker

		// pre-test: create another worker
		makeAWorker(t, db, "IsAHorcrux")
		knownWorkersList := requireNumberOfWorkers(t, db, 2)

		// set up test helpers
		// basically, we hook into the run process to know when all workers have run.
		// once we know that number, we set up a test to verify that we reach that number
		numWorkers := len(knownWorkersList)
		calledCh := make(chan bool, numWorkers*2) // making extra room, in case there are extras (which would be a failure)
		mockHandler := buildRequestMock(
			func(w *httptest.ResponseRecorder) {
				calledCh <- true
				w.WriteHeader(http.StatusOK)
				w.Write(mustJSONMarshal(t, enhancementservices.TestResp{Status: "ok"}))
			})

		for _, v := range knownWorkersList {
			enhancementservices.SetWebRequestFunctionForWorker(v.Name, &mockHandler)
		}

		tryRun := func(u models.User, input services.RunServiceWorkerInput) error {
			ctx := contextForUser(u, db)
			return services.RunServiceWorker(ctx, db, input)
		}
		allWorkersCalled := makeNotifierChannel()

		evi := EviDobby
		op := seed.OperationForEvidence(evi)
		input := services.RunServiceWorkerInput{
			OperationSlug: op.Slug,
			EvidenceUUID:  evi.UUID,
			WorkerName:    singleWorker.Name,
		}

		// verify permissions
		require.Error(t, tryRun(UserDraco, input)) // no access

		// verify result (single worker)
		require.NoError(t, tryRun(UserRon, input))
		<-allWorkersCalled // wait for the work to complete
		require.Equal(t, 1, len(calledCh))
		<-calledCh // empty the channel

		// verify result (all workers)
		input.WorkerName = "" // set the trigger for executing all workers
		require.NoError(t, tryRun(UserRon, input))
		<-allWorkersCalled // wait for the work to complete
		require.Equal(t, numWorkers, len(calledCh))
	})
}

func TestBatchRunServiceWorker(t *testing.T) {
	// This test is complicated. Basically, we need to hook into many systems and replace their
	// default operation with a mock version. Once we do all of that, then we can verify that
	// a worker was called.

	RunResettableDBTest(t, func(db *database.Connection, seed TestSeedData) {
		singleWorker := DemoServiceWorker

		// pre-test: create another worker
		altWorker := makeAWorker(t, db, "IsAHorcrux")
		knownWorkersList := []models.ServiceWorker{
			singleWorker,
			altWorker,
		}
		numWorkers := len(knownWorkersList)

		op := OpChamberOfSecrets
		setOfEvidence := seed.EvidenceForOperation(op.ID)
		expectedCalls := len(setOfEvidence) * numWorkers
		calledCh := make(chan bool, expectedCalls*2) // making extra room, in case there are extras (which would be a failure)
		allWorkersCalled := makeNotifierChannel()

		tryBatchRun := func(u models.User, input services.BatchRunServiceWorkerInput) error {
			ctx := contextForUser(u, db)
			return services.BatchRunServiceWorker(ctx, db, input)
		}
		mockHandler := buildRequestMock(
			func(w *httptest.ResponseRecorder) {
				calledCh <- true
				w.WriteHeader(http.StatusOK)
				// we actually don't really care about the data -- other tests establish that the data
				// needs to match a specific shape to make sense and produce the correct result.
				w.Write(mustJSONMarshal(t, enhancementservices.TestResp{Status: "ok"}))
			})
		for _, v := range knownWorkersList {
			enhancementservices.SetWebRequestFunctionForWorker(v.Name, &mockHandler)
		}

		input := services.BatchRunServiceWorkerInput{
			OperationSlug: op.Slug,
			EvidenceUUIDs: helpers.Map(setOfEvidence, func(e models.Evidence) string { return e.UUID }),
			WorkerNames:   helpers.Map(knownWorkersList, func(w models.ServiceWorker) string { return w.Name }),
		}

		// verify permissions
		require.Error(t, tryBatchRun(UserDraco, input)) // no access

		// verify result
		require.NoError(t, tryBatchRun(UserRon, input))
		<-allWorkersCalled // wait for the work to complete
		require.Equal(t, expectedCalls, len(calledCh))
	})
}

func buildRequestMock(writeResponse func(*httptest.ResponseRecorder)) enhancementservices.RequestFn {
	return func(method, url string, body io.Reader, updateRequest helpers.ModifyReqFunc) (*http.Response, error) {
		w := httptest.NewRecorder()

		writeResponse(w)
		return w.Result(), nil
	}
}

func mustJSONMarshal(t *testing.T, obj any) []byte {
	rtn, err := json.Marshal(obj)
	require.NoError(t, err)
	return rtn
}

func makeAWorker(t *testing.T, db *database.Connection, workerName string) models.ServiceWorker {
	ctx := contextForUser(UserDumbledore, db)
	require.NoError(t, services.CreateServiceWorker(ctx, db, services.CreateServiceWorkerInput{
		Name:   workerName,
		Config: `{"type": "web", "version": 1, "url": "http://test:1234"}`,
	}))
	return getServiceWorkerByName(t, db, workerName)
}

func requireNumberOfWorkers(t *testing.T, db *database.Connection, minWorkers int) []models.ServiceWorker {
	knownWorkersList := listServiceWorkers(t, db)
	numWorkers := len(knownWorkersList)
	require.GreaterOrEqual(t, minWorkers, numWorkers)
	return knownWorkersList
}

func makeNotifierChannel() chan bool {
	allWorkersCalled := make(chan bool, 1)
	enhancementservices.SetNotifyWorkersRunForTest(allWorkersCalled)
	return allWorkersCalled
}
