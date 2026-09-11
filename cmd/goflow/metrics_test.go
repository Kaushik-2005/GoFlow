package main

import "testing"

func TestMetricsSnapshotIncludesCountersAndGauges(t *testing.T) {
	metrics := newMetrics()

	metrics.incrementHTTPRequests()
	metrics.incrementJobsSubmitted()
	metrics.incrementJobsCompleted()
	metrics.incrementJobsExecutionFailed()
	metrics.incrementJobsRetried()
	metrics.incrementJobsDeadLettered()
	metrics.incrementActiveWorkers()
	metrics.incrementQueueDepth()
	metrics.decrementActiveWorkers()
	metrics.decrementQueueDepth()

	snapshot := metrics.snapshot()
	if snapshot.HTTPRequestsTotal != 1 {
		t.Fatalf("expected http_requests_total 1, got %d", snapshot.HTTPRequestsTotal)
	}
	if snapshot.JobsSubmittedTotal != 1 {
		t.Fatalf("expected jobs_submitted_total 1, got %d", snapshot.JobsSubmittedTotal)
	}
	if snapshot.JobsCompletedTotal != 1 {
		t.Fatalf("expected jobs_completed_total 1, got %d", snapshot.JobsCompletedTotal)
	}
	if snapshot.JobsExecutionFailedTotal != 1 {
		t.Fatalf("expected jobs_execution_failed_total 1, got %d", snapshot.JobsExecutionFailedTotal)
	}
	if snapshot.JobsRetriedTotal != 1 {
		t.Fatalf("expected jobs_retried_total 1, got %d", snapshot.JobsRetriedTotal)
	}
	if snapshot.JobsDeadLetteredTotal != 1 {
		t.Fatalf("expected jobs_dead_lettered_total 1, got %d", snapshot.JobsDeadLetteredTotal)
	}
	if snapshot.ActiveWorkers != 0 {
		t.Fatalf("expected active_workers 0, got %d", snapshot.ActiveWorkers)
	}
	if snapshot.QueueDepth != 0 {
		t.Fatalf("expected queue_depth 0, got %d", snapshot.QueueDepth)
	}
}

func TestObserveHTTPRequestDurationBuckets(t *testing.T) {
	metrics := newMetrics()

	metrics.observeHTTPRequestDuration(10)
	metrics.observeHTTPRequestDuration(100)
	metrics.observeHTTPRequestDuration(1000)
	metrics.observeHTTPRequestDuration(1001)

	snapshot := metrics.snapshot()
	if snapshot.HTTPRequestDurationMSLess10 != 1 {
		t.Fatalf("expected le_10 bucket 1, got %d", snapshot.HTTPRequestDurationMSLess10)
	}
	if snapshot.HTTPRequestDurationMSLess100 != 1 {
		t.Fatalf("expected le_100 bucket 1, got %d", snapshot.HTTPRequestDurationMSLess100)
	}
	if snapshot.HTTPRequestDurationMSLess1000 != 1 {
		t.Fatalf("expected le_1000 bucket 1, got %d", snapshot.HTTPRequestDurationMSLess1000)
	}
	if snapshot.HTTPRequestDurationMSMore1000 != 1 {
		t.Fatalf("expected gt_1000 bucket 1, got %d", snapshot.HTTPRequestDurationMSMore1000)
	}
}
