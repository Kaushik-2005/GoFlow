package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const requestIDContextKey contextKey = "request_id"

type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func requestLoggingMiddleware(logger *slog.Logger, metrics *metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &statusResponseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			durationMS := time.Since(start).Milliseconds()
			if metrics != nil {
				metrics.incrementHTTPRequests()
				metrics.observeHTTPRequestDuration(durationMS)
			}

			requestID, _ := r.Context().Value(requestIDContextKey).(string)
			logger.InfoContext(
				r.Context(),
				"http request completed",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", wrapped.status,
				"duration_ms", durationMS,
			)
		})
	}
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				writeJSONError(
					w,
					http.StatusInternalServerError,
					"INTERNAL_SERVER_ERROR",
					"An internal server error occurred",
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), requestIDContextKey, requestID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func requestBodyLimitMiddleware(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

func requireJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.Header.Get("Content-Type") != "application/json" {
			writeJSONError(
				w,
				http.StatusUnsupportedMediaType,
				"UNSUPPORTED_MEDIA_TYPE",
				"Content-Type must be application/json",
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
