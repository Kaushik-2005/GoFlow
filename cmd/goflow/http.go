package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
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

func liveHandler(w http.ResponseWriter, r *http.Request) {
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

func jobsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listJobsHandler(w, r)
	case http.MethodPost:
		createJobHandler(w, r)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
	}
}

func listJobsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	store, err := LoadStore(jobsFile)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "STORE_LOAD_FAILED", "Failed to load jobs")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(store.List())
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/jobs/")
	if id == "" {
		writeJSONError(w, http.StatusBadRequest, "JOB_ID_REQUIRED", "A job ID is required")
		return
	}

	store, err := LoadStore(jobsFile)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "STORE_LOAD_FAILED", "Failed to load jobs")
		return
	}

	job, err := store.Get(id)
	if err != nil {
		if errors.Is(err, ErrJobNotFound) {
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

type createJobRequest struct {
	Type string `json:"type"`
}

func createJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "The requested method is not allowed for this endpoint")
		return
	}

	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "INVALID_REQUEST_BODY", "The request body must be valid JSON")
		return
	}

	if req.Type == "" {
		writeJSONError(w, http.StatusBadRequest, "JOB_TYPE_REQUIRED", "Job type is required")
		return
	}

	store, err := LoadStore(jobsFile)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "STORE_LOAD_FAILED", "Failed to load jobs")
		return
	}

	job := Job{
		ID:          fmt.Sprintf("job-%d", len(store.List())+1),
		Type:        req.Type,
		Status:      StatusPending,
		Attempts:    0,
		MaxAttempts: 3,
	}

	if err := store.Create(job); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "JOB_CREATE_FAILED", "Failed to create job")
		return
	}

	if err := store.Save(jobsFile); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "JOB_SAVE_FAILED", "Failed to save job")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(job)
}
