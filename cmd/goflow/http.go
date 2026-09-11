package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goflow/internal/goflow"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type readinessChecker interface {
	Ping(ctx context.Context) error
}

type apiHandler struct {
	store   goflow.JobStore
	ready   readinessChecker
	logger  *slog.Logger
	metrics *metrics
}

func newAPIHandler(store goflow.JobStore) *apiHandler {
	return &apiHandler{
		store:   store,
		logger:  slog.New(slog.DiscardHandler),
		metrics: newMetrics(),
	}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(errorResponse{
		Error: apiError{
			Code:    code,
			Message: message,
		},
	})
}

func (h *apiHandler) liveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *apiHandler) readyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	if h.ready == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "SERVICE_NOT_READY", "The service is not ready")
		return
	}

	if err := h.ready.Ping(r.Context()); err != nil {
		writeJSONError(w, http.StatusServiceUnavailable, "DATABASE_NOT_READY", "The database is not ready")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
	})
}

func (h *apiHandler) metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(h.metrics.snapshot())
}

func (h *apiHandler) jobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listJobsHandler(w, r)
	case http.MethodPost:
		h.createJobHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
	}
}

func (h *apiHandler) jobByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getJobHandler(w, r)
	case http.MethodDelete:
		h.deleteJobHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
	}
}

func (h *apiHandler) listJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	jobs, err := h.store.List(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "JOB_LIST_FAILED", "Failed to list jobs")
		return
	}

	statusFilter := r.URL.Query().Get("status")
	if statusFilter != "" {
		filtered := make([]goflow.Job, 0, len(jobs))
		for _, job := range jobs {
			if string(job.Status) == statusFilter {
				filtered = append(filtered, job)
			}
		}
		jobs = filtered
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(jobs)
}

func (h *apiHandler) getJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/jobs/")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "JOB_ID_REQUIRED", "A job ID is required")
		return
	}

	job, err := h.store.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, goflow.ErrJobNotFound) {
			writeJSONError(w, http.StatusNotFound, "JOB_NOT_FOUND", "The requested job does not exist")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "JOB_GET_FAILED", "Failed to get job")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(job)
}

func (h *apiHandler) deleteJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/jobs/")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "JOB_ID_REQUIRED", "A job ID is required")
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		if errors.Is(err, goflow.ErrJobNotFound) {
			writeJSONError(w, http.StatusNotFound, "JOB_NOT_FOUND", "The requested job does not exist")
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "JOB_DELETE_FAILED", "Failed to delete job")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type createJobRequest struct {
	Type string `json:"type"`
}

func (h *apiHandler) createJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeJSONError(
				w,
				http.StatusRequestEntityTooLarge,
				"REQUEST_BODY_TOO_LARGE",
				"The request body exceeds the allowed size",
			)
			return
		}

		writeJSONError(
			w,
			http.StatusBadRequest,
			"INVALID_REQUEST_BODY",
			"The request body is invalid or malformed",
		)
		return
	}

	if req.Type == "" {
		writeJSONError(w, http.StatusBadRequest, "JOB_TYPE_REQUIRED", "Job type is required")
		return
	}

	jobID, err := newJobID()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "JOB_ID_GENERATION_FAILED", "Failed to generate job ID")
		return
	}

	job := goflow.Job{
		ID:          jobID,
		Type:        req.Type,
		Status:      goflow.StatusPending,
		Attempts:    0,
		MaxAttempts: 3,
		AvailableAt: time.Now(),
	}

	if err := h.store.Create(r.Context(), job); err != nil {
		if errors.Is(err, goflow.ErrJobAlreadyExists) {
			writeJSONError(w, http.StatusConflict, "JOB_ALREADY_EXISTS", fmt.Sprintf("Job %s already exists", job.ID))
			return
		}

		writeJSONError(w, http.StatusInternalServerError, "JOB_CREATE_FAILED", "Failed to create job")
		return
	}

	requestID, _ := r.Context().Value(requestIDContextKey).(string)
	h.logger.InfoContext(
		r.Context(),
		"job created",
		"request_id", requestID,
		"job_id", job.ID,
		"job_type", job.Type,
	)

	if h.metrics != nil {
		h.metrics.incrementJobsSubmitted()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(job)
}
