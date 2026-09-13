package main

import "testing"

func BenchmarkMetricsSnapshot(b *testing.B) {
	metrics := newMetrics()

	metrics.incrementHTTPRequests()
	metrics.incrementJobsSubmitted()
	metrics.incrementJobsCompleted()
	metrics.incrementJobsExecutionFailed()
	metrics.incrementJobsRetried()
	metrics.incrementJobsDeadLettered()
	metrics.incrementActiveWorkers()
	metrics.incrementQueueDepth()

	for i := 0; i < b.N; i++ {
		_ = metrics.snapshot()
	}
}

func BenchmarkMetricsSnapshotParallel(b *testing.B) {
	metrics := newMetrics()

	metrics.incrementHTTPRequests()
	metrics.incrementJobsSubmitted()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = metrics.snapshot()
		}
	})
}
