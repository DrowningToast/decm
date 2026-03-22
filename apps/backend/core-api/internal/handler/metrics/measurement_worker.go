package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var workerMeasurement = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "decm_worker_duration_milliseconds",
	// Worker operations are slower (blockchain tx, DB polling) — buckets skew higher.
	Buckets: []float64{50, 100, 200, 400, 800, 2000, 5000, 10000, 30000, 60000},
}, []string{"worker", "status"})

// MeasureWorkerOperation starts a timer and returns a finish func to call
// after the operation completes.  Pass the operation error (or nil) so the
// correct status label ("success" / "failure") is recorded.
//
// Usage:
//
//	finish := metrics.MeasureWorkerOperation("blockchainSubmission.ProcessEventJoin")
//	err := doWork(ctx)
//	finish(err)
func MeasureWorkerOperation(operation string) func(err error) time.Duration {
	startTime := time.Now()

	return func(err error) time.Duration {
		elapsed := time.Since(startTime)

		status := "success"
		if err != nil {
			status = "failure"
		}

		workerMeasurement.WithLabelValues(operation, status).Observe(float64(elapsed.Milliseconds()))

		return elapsed
	}
}
