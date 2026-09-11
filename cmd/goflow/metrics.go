package main

import "sync"

type metrics struct {
	mu sync.Mutex

	httpRequestsTotal             int
	httpRequestDurationMSLess10   int
	httpRequestDurationMSLess100  int
	httpRequestDurationMSLess1000 int
	httpRequestDurationMSMore1000 int
	jobsSubmittedTotal            int
	jobsCompletedTotal            int
	jobsExecutionFailedTotal      int
	jobsRetriedTotal              int
	jobsDeadLetteredTotal         int
	activeWorkers                 int
	queueDepth                    int
}

type metricsSnapshot struct {
	HTTPRequestsTotal             int `json:"http_requests_total"`
	HTTPRequestDurationMSLess10   int `json:"http_request_duration_ms_le_10"`
	HTTPRequestDurationMSLess100  int `json:"http_request_duration_ms_le_100"`
	HTTPRequestDurationMSLess1000 int `json:"http_request_duration_ms_le_1000"`
	HTTPRequestDurationMSMore1000 int `json:"http_request_duration_ms_gt_1000"`
	JobsSubmittedTotal            int `json:"jobs_submitted_total"`
	JobsCompletedTotal            int `json:"jobs_completed_total"`
	JobsExecutionFailedTotal      int `json:"jobs_execution_failed_total"`
	JobsRetriedTotal              int `json:"jobs_retried_total"`
	JobsDeadLetteredTotal         int `json:"jobs_dead_lettered_total"`
	ActiveWorkers                 int `json:"active_workers"`
	QueueDepth                    int `json:"queue_depth"`
}

func newMetrics() *metrics {
	return &metrics{}
}

func (m *metrics) incrementHTTPRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.httpRequestsTotal++
}

func (m *metrics) observeHTTPRequestDuration(durationMS int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch {
	case durationMS <= 10:
		m.httpRequestDurationMSLess10++
	case durationMS <= 100:
		m.httpRequestDurationMSLess100++
	case durationMS <= 1000:
		m.httpRequestDurationMSLess1000++
	default:
		m.httpRequestDurationMSMore1000++
	}
}

func (m *metrics) incrementJobsSubmitted() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsSubmittedTotal++
}

func (m *metrics) incrementJobsCompleted() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsCompletedTotal++
}

func (m *metrics) incrementJobsExecutionFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsExecutionFailedTotal++
}

func (m *metrics) incrementJobsRetried() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsRetriedTotal++
}

func (m *metrics) incrementJobsDeadLettered() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobsDeadLetteredTotal++
}

func (m *metrics) incrementActiveWorkers() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.activeWorkers++
}

func (m *metrics) decrementActiveWorkers() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.activeWorkers--
}

func (m *metrics) incrementQueueDepth() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queueDepth++
}

func (m *metrics) decrementQueueDepth() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.queueDepth--
}

func (m *metrics) snapshot() metricsSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()

	return metricsSnapshot{
		HTTPRequestsTotal:             m.httpRequestsTotal,
		HTTPRequestDurationMSLess10:   m.httpRequestDurationMSLess10,
		HTTPRequestDurationMSLess100:  m.httpRequestDurationMSLess100,
		HTTPRequestDurationMSLess1000: m.httpRequestDurationMSLess1000,
		HTTPRequestDurationMSMore1000: m.httpRequestDurationMSMore1000,
		JobsSubmittedTotal:            m.jobsSubmittedTotal,
		JobsCompletedTotal:            m.jobsCompletedTotal,
		JobsExecutionFailedTotal:      m.jobsExecutionFailedTotal,
		JobsRetriedTotal:              m.jobsRetriedTotal,
		JobsDeadLetteredTotal:         m.jobsDeadLetteredTotal,
		ActiveWorkers:                 m.activeWorkers,
		QueueDepth:                    m.queueDepth,
	}
}
