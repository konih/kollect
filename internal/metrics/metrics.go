// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

const (
	ResultSuccess = "success"
	ResultFailure = "failure"

	ErrorClassTransient = "transient"
	ErrorClassTerminal  = "terminal"
	ErrorClassForbidden = "forbidden"
)

var (
	InventoryItemsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kollect_inventory_items_total",
			Help: "Number of inventory items in the last aggregated snapshot.",
		},
	)

	CollectItemsTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kollect_collect_items_total",
			Help: "Number of items currently held in the in-memory collection store.",
		},
	)

	CollectedObjects = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kollect_collected_objects",
			Help: "Collected objects by profile and GVK.",
		},
		[]string{"profile", "gvk"},
	)

	ReconcileTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kollect_reconcile_total",
			Help: "Reconcile attempts by controller and result.",
		},
		[]string{"controller", "result"},
	)

	ReconcileErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kollect_reconcile_errors_total",
			Help: "Reconcile errors by kind and error class.",
		},
		[]string{"kind", "error_class"},
	)

	ExportDurationSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kollect_export_duration_seconds",
			Help:    "Sink export duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"sink_type"},
	)

	SinkConnectionTestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kollect_sink_connection_test_total",
			Help: "Git/TLS sink connection tests by sink type and result.",
		},
		[]string{"type", "result"},
	)
)

// Register adds kollect custom metrics to the controller-runtime registry.
func Register() {
	metrics.Registry.MustRegister(
		InventoryItemsTotal,
		CollectItemsTotal,
		CollectedObjects,
		ReconcileTotal,
		ReconcileErrorsTotal,
		ExportDurationSeconds,
		SinkConnectionTestTotal,
	)
}
