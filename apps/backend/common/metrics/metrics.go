package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// S3 Metrics
	S3OperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "decm_external_s3_operations_total",
			Help: "The total number of S3 operations",
		},
		[]string{"operation", "status"},
	)
	S3OperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "decm_external_s3_operation_duration_seconds",
			Help: "The duration of S3 operations",
		},
		[]string{"operation", "status"},
	)

	// Google OAuth Metrics
	GoogleOAuthOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "decm_external_google_oauth_operations_total",
			Help: "The total number of Google OAuth operations",
		},
		[]string{"operation", "status"},
	)
	GoogleOAuthOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "decm_external_google_oauth_operation_duration_seconds",
			Help: "The duration of Google OAuth operations",
		},
		[]string{"operation", "status"},
	)

	// Blockchain Metrics
	BlockchainOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "decm_external_blockchain_operations_total",
			Help: "The total number of Blockchain operations",
		},
		[]string{"operation", "status"},
	)
	BlockchainOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "decm_external_blockchain_operation_duration_seconds",
			Help: "The duration of Blockchain operations",
		},
		[]string{"operation", "status"},
	)
	BlockchainGasPrice = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "decm_external_blockchain_gas_price_gwei",
			Help: "The current gas price in Gwei",
		},
		[]string{"type"}, // base_fee, priority_fee, max_fee
	)
)