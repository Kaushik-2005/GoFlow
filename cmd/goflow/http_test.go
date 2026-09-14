package main

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"goflow/internal/goflow"
)

type fakeAPIStore struct {
	job       goflow.Job
	jobs      []goflow.Job
	err       error
	created   goflow.Job
	createHit bool
	deletedID string
	deleteHit bool
	pingErr   error
}

func (f *fakeAPIStore) Create(ctx context.Context, job goflow.Job) error {
	if f.err != nil {
		return f.err
	}

	f.created = job
	f.createHit = true
	return nil
}

func (f *fakeAPIStore) Get(ctx context.Context, id string) (goflow.Job, error) {
	if f.err != nil {
		return goflow.Job{}, f.err
	}
	return f.job, nil
}

func (f *fakeAPIStore) List(ctx context.Context) ([]goflow.Job, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.jobs, nil
}

func (f *fakeAPIStore) Update(ctx context.Context, job goflow.Job) error {
	return nil
}

func (f *fakeAPIStore) Ping(ctx context.Context) error {
	return f.pingErr
}

func (f *fakeAPIStore) Delete(ctx context.Context, id string) error {
	if f.err != nil {
		return f.err
	}

	f.deletedID = id
	f.deleteHit = true
	return nil
}

func assertJSONResponse(t *testing.T, response *httptest.ResponseRecorder, wantStatus int, wantContains ...string) {
	t.Helper()

	if response.Code != wantStatus {
		t.Fatalf("expected status %d, got %d", wantStatus, response.Code)
	}

	if got := response.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected content type %q, got %q", "application/json", got)
	}

	for _, want := range wantContains {
		if !strings.Contains(response.Body.String(), want) {
			t.Fatalf("expected response body to contain %q, got %q", want, response.Body.String())
		}
	}
}

func TestLiveHandler(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})

	tests := []struct {
		name         string
		method       string
		wantStatus   int
		wantContains []string
	}{
		{
			name:         "get succeeds",
			method:       http.MethodGet,
			wantStatus:   http.StatusOK,
			wantContains: []string{`"status":"ok"`},
		},
		{
			name:         "post rejected",
			method:       http.MethodPost,
			wantStatus:   http.StatusMethodNotAllowed,
			wantContains: []string{`"code":"METHOD_NOT_ALLOWED"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, "/health/live", nil)
			response := httptest.NewRecorder()

			api.liveHandler(response, request)
			assertJSONResponse(t, response, tt.wantStatus, tt.wantContains...)
		})
	}
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		ready        readinessChecker
		wantStatus   int
		wantContains []string
	}{
		{
			name:         "ready",
			method:       http.MethodGet,
			ready:        &fakeAPIStore{},
			wantStatus:   http.StatusOK,
			wantContains: []string{`"status":"ready"`},
		},
		{
			name:         "database not ready",
			method:       http.MethodGet,
			ready:        &fakeAPIStore{pingErr: errors.New("database unavailable")},
			wantStatus:   http.StatusServiceUnavailable,
			wantContains: []string{`"code":"DATABASE_NOT_READY"`},
		},
		{
			name:         "post rejected",
			method:       http.MethodPost,
			ready:        &fakeAPIStore{},
			wantStatus:   http.StatusMethodNotAllowed,
			wantContains: []string{`"code":"METHOD_NOT_ALLOWED"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPIHandler(&fakeAPIStore{})
			api.ready = tt.ready
			request := httptest.NewRequest(tt.method, "/health/ready", nil)
			response := httptest.NewRecorder()

			api.readyHandler(response, request)
			assertJSONResponse(t, response, tt.wantStatus, tt.wantContains...)
		})
	}
}

func TestMetricsHandler(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})
	api.metrics.incrementHTTPRequests()
	api.metrics.incrementJobsSubmitted()

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	response := httptest.NewRecorder()

	api.metricsHandler(response, request)
	assertJSONResponse(
		t,
		response,
		http.StatusOK,
		`"http_requests_total":1`,
		`"jobs_submitted_total":1`,
		`"active_workers":0`,
		`"queue_depth":0`,
	)
}

func TestMetricsHandlerMethodNotAllowed(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})
	request := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	response := httptest.NewRecorder()

	api.metricsHandler(response, request)
	assertJSONResponse(t, response, http.StatusMethodNotAllowed, `"code":"METHOD_NOT_ALLOWED"`)
}

func TestJobsHandlerMethodNotAllowed(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})
	request := httptest.NewRequest(http.MethodPatch, "/v1/jobs", nil)
	response := httptest.NewRecorder()

	api.jobsHandler(response, request)
	assertJSONResponse(t, response, http.StatusMethodNotAllowed, `"code":"METHOD_NOT_ALLOWED"`)
}

func TestJobByIDHandlerMethodNotAllowed(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})
	request := httptest.NewRequest(http.MethodPatch, "/v1/jobs/job-1", nil)
	response := httptest.NewRecorder()

	api.jobByIDHandler(response, request)
	assertJSONResponse(t, response, http.StatusMethodNotAllowed, `"code":"METHOD_NOT_ALLOWED"`)
}

func TestListJobsHandler(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		store        *fakeAPIStore
		wantStatus   int
		wantContains []string
		wantMissing  []string
	}{
		{
			name: "success",
			path: "/v1/jobs",
			store: &fakeAPIStore{jobs: []goflow.Job{
				{ID: "job-1", Type: "email", Status: goflow.StatusPending},
				{ID: "job-2", Type: "report", Status: goflow.StatusRunning},
			}},
			wantStatus:   http.StatusOK,
			wantContains: []string{`"id":"job-1"`, `"id":"job-2"`, `"status":"running"`},
		},
		{
			name: "status filter",
			path: "/v1/jobs?status=running",
			store: &fakeAPIStore{jobs: []goflow.Job{
				{ID: "job-1", Type: "email", Status: goflow.StatusPending},
				{ID: "job-2", Type: "report", Status: goflow.StatusRunning},
			}},
			wantStatus:   http.StatusOK,
			wantContains: []string{`"id":"job-2"`, `"status":"running"`},
			wantMissing:  []string{`"id":"job-1"`},
		},
		{
			name:         "store failure",
			path:         "/v1/jobs",
			store:        &fakeAPIStore{err: goflow.ErrJobNotFound},
			wantStatus:   http.StatusInternalServerError,
			wantContains: []string{`"code":"JOB_LIST_FAILED"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPIHandler(tt.store)
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			api.listJobsHandler(response, request)
			assertJSONResponse(t, response, tt.wantStatus, tt.wantContains...)

			for _, unwanted := range tt.wantMissing {
				if strings.Contains(response.Body.String(), unwanted) {
					t.Fatalf("expected response body not to contain %q, got %q", unwanted, response.Body.String())
				}
			}
		})
	}
}

func TestGetJobHandler(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		store        *fakeAPIStore
		wantStatus   int
		wantContains []string
	}{
		{
			name: "success",
			path: "/v1/jobs/job-1111111111111111",
			store: &fakeAPIStore{job: goflow.Job{
				ID:     "job-1111111111111111",
				Type:   "email",
				Status: goflow.StatusPending,
			}},
			wantStatus:   http.StatusOK,
			wantContains: []string{`"id":"job-1111111111111111"`, `"type":"email"`},
		},
		{
			name:         "not found",
			path:         "/v1/jobs/job-9999999999999999",
			store:        &fakeAPIStore{err: goflow.ErrJobNotFound},
			wantStatus:   http.StatusNotFound,
			wantContains: []string{`"code":"JOB_NOT_FOUND"`},
		},
		{
			name:         "store failure",
			path:         "/v1/jobs/job-5005005005005005",
			store:        &fakeAPIStore{err: context.DeadlineExceeded},
			wantStatus:   http.StatusInternalServerError,
			wantContains: []string{`"code":"JOB_GET_FAILED"`},
		},
		{
			name:         "missing id",
			path:         "/v1/jobs/",
			store:        &fakeAPIStore{},
			wantStatus:   http.StatusBadRequest,
			wantContains: []string{`"code":"JOB_ID_REQUIRED"`},
		},
		{
			name:         "invalid id",
			path:         "/v1/jobs/abc",
			store:        &fakeAPIStore{},
			wantStatus:   http.StatusBadRequest,
			wantContains: []string{`"code":"INVALID_JOB_ID"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPIHandler(tt.store)
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()

			api.getJobHandler(response, request)
			assertJSONResponse(t, response, tt.wantStatus, tt.wantContains...)
		})
	}
}

func TestDeleteJobHandler(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		store         *fakeAPIStore
		wantStatus    int
		wantDeleteID  string
		wantDeleteHit bool
		wantContains  []string
	}{
		{
			name:          "success",
			path:          "/v1/jobs/job-1111111111111111",
			store:         &fakeAPIStore{},
			wantStatus:    http.StatusNoContent,
			wantDeleteID:  "job-1111111111111111",
			wantDeleteHit: true,
		},
		{
			name:          "not found",
			path:          "/v1/jobs/job-9999999999999999",
			store:         &fakeAPIStore{err: goflow.ErrJobNotFound},
			wantStatus:    http.StatusNotFound,
			wantDeleteHit: false,
			wantContains:  []string{`"code":"JOB_NOT_FOUND"`},
		},
		{
			name:          "store failure",
			path:          "/v1/jobs/job-5005005005005005",
			store:         &fakeAPIStore{err: context.DeadlineExceeded},
			wantStatus:    http.StatusInternalServerError,
			wantDeleteHit: false,
			wantContains:  []string{`"code":"JOB_DELETE_FAILED"`},
		},
		{
			name:          "missing id",
			path:          "/v1/jobs/",
			store:         &fakeAPIStore{},
			wantStatus:    http.StatusBadRequest,
			wantDeleteHit: false,
			wantContains:  []string{`"code":"JOB_ID_REQUIRED"`},
		},
		{
			name:          "invalid id",
			path:          "/v1/jobs/abc",
			store:         &fakeAPIStore{},
			wantStatus:    http.StatusBadRequest,
			wantDeleteHit: false,
			wantContains:  []string{`"code":"INVALID_JOB_ID"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPIHandler(tt.store)
			request := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			response := httptest.NewRecorder()

			api.deleteJobHandler(response, request)

			if tt.wantStatus == http.StatusNoContent {
				if response.Code != tt.wantStatus {
					t.Fatalf("expected status %d, got %d", tt.wantStatus, response.Code)
				}
				if response.Body.Len() != 0 {
					t.Fatalf("expected empty body, got %q", response.Body.String())
				}
			} else {
				assertJSONResponse(t, response, tt.wantStatus, tt.wantContains...)
			}

			if tt.store.deleteHit != tt.wantDeleteHit {
				t.Fatalf("expected deleteHit %t, got %t", tt.wantDeleteHit, tt.store.deleteHit)
			}

			if tt.wantDeleteID != "" && tt.store.deletedID != tt.wantDeleteID {
				t.Fatalf("expected deleted id %q, got %q", tt.wantDeleteID, tt.store.deletedID)
			}
		})
	}
}

func TestCreateJobHandler(t *testing.T) {
	tests := []struct {
		name            string
		body            string
		store           *fakeAPIStore
		wantStatus      int
		wantCreateHit   bool
		wantCreatedType string
		wantContains    []string
	}{
		{
			name:            "success",
			body:            `{"type":"email"}`,
			store:           &fakeAPIStore{},
			wantStatus:      http.StatusCreated,
			wantCreateHit:   true,
			wantCreatedType: "email",
			wantContains:    []string{`"type":"email"`},
		},
		{
			name:          "missing type",
			body:          `{}`,
			store:         &fakeAPIStore{},
			wantStatus:    http.StatusBadRequest,
			wantCreateHit: false,
			wantContains:  []string{`"code":"JOB_TYPE_REQUIRED"`},
		},
		{
			name:          "unsupported type",
			body:          `{"type":"temporary-fail"}`,
			store:         &fakeAPIStore{},
			wantStatus:    http.StatusBadRequest,
			wantCreateHit: false,
			wantContains:  []string{`"code":"INVALID_JOB_TYPE"`},
		},
		{
			name:          "invalid json",
			body:          `{"type":`,
			store:         &fakeAPIStore{},
			wantStatus:    http.StatusBadRequest,
			wantCreateHit: false,
			wantContains:  []string{`"code":"INVALID_REQUEST_BODY"`},
		},
		{
			name:          "duplicate job",
			body:          `{"type":"email"}`,
			store:         &fakeAPIStore{err: goflow.ErrJobAlreadyExists},
			wantStatus:    http.StatusConflict,
			wantCreateHit: false,
			wantContains:  []string{`"code":"JOB_ALREADY_EXISTS"`},
		},
		{
			name:          "store failure",
			body:          `{"type":"email"}`,
			store:         &fakeAPIStore{err: context.DeadlineExceeded},
			wantStatus:    http.StatusInternalServerError,
			wantCreateHit: false,
			wantContains:  []string{`"code":"JOB_CREATE_FAILED"`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPIHandler(tt.store)
			request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(tt.body))
			response := httptest.NewRecorder()

			api.createJobHandler(response, request)
			assertJSONResponse(t, response, tt.wantStatus, tt.wantContains...)

			if tt.store.createHit != tt.wantCreateHit {
				t.Fatalf("expected createHit %t, got %t", tt.wantCreateHit, tt.store.createHit)
			}

			if tt.wantCreatedType != "" && tt.store.created.Type != tt.wantCreatedType {
				t.Fatalf("expected created job type %q, got %q", tt.wantCreatedType, tt.store.created.Type)
			}

			if tt.wantCreateHit && tt.store.created.Status != goflow.StatusPending {
				t.Fatalf("expected created job status %q, got %q", goflow.StatusPending, tt.store.created.Status)
			}
		})
	}
}

func TestCreateJobHandlerIncrementsSubmittedMetric(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})
	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(`{"type":"email"}`))
	response := httptest.NewRecorder()

	api.createJobHandler(response, request)
	assertJSONResponse(t, response, http.StatusCreated, `"type":"email"`)

	snapshot := api.metrics.snapshot()
	if snapshot.JobsSubmittedTotal != 1 {
		t.Fatalf("expected jobs_submitted_total 1, got %d", snapshot.JobsSubmittedTotal)
	}
}

func TestRequireJSONMiddlewareAcceptsJSONWithCharset(t *testing.T) {
	api := newAPIHandler(&fakeAPIStore{})

	handler := requireJSONMiddleware(http.HandlerFunc(api.createJobHandler))

	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(`{"type":"email"}`))
	request.Header.Set("Content-Type", "application/json; charset=8")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	assertJSONResponse(t, response, http.StatusCreated, `"type":"email"`)
}
