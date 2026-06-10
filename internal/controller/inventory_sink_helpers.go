// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/sink"
)

func loadResolvedSink(
	ctx context.Context,
	c client.Client,
	namespace string,
	binding kollectdevv1alpha1.InventorySinkBinding,
) (*sink.ResolvedSink, error) {
	opts := sink.ResolveOptionsForBinding(namespace, binding)
	resolved, err := sink.ResolveSink(ctx, c, opts)
	if err != nil {
		kind := familySinkKind(binding.Family)
		return nil, fmt.Errorf("get %s %q in %q: %w", kind, binding.Name, namespace, err)
	}

	return resolved, nil
}

func familySinkKind(family string) string {
	switch family {
	case kollectdevv1alpha1.SinkFamilySnapshot:
		return "KollectSnapshotSink"
	case kollectdevv1alpha1.SinkFamilyDatabase:
		return "KollectDatabaseSink"
	case kollectdevv1alpha1.SinkFamilyEvent:
		return "KollectEventSink"
	default:
		return "KollectSink"
	}
}

// sinkBindingNamespace resolves the namespace for a cluster inventory sink ref: the per-ref
// namespace when set, otherwise the inventory default sink namespace (ADR-0208).
func sinkBindingNamespace(binding kollectdevv1alpha1.InventorySinkBinding, defaultNS string) string {
	if ns := strings.TrimSpace(binding.Ref.Namespace); ns != "" {
		return binding.Ref.Namespace
	}

	return defaultNS
}

func inventorySinkBindings(inv *kollectdevv1alpha1.KollectInventory) []kollectdevv1alpha1.InventorySinkBinding {
	if inv == nil {
		return nil
	}

	return kollectdevv1alpha1.CollectInventorySinkBindings(&inv.Spec)
}

func clusterInventorySinkBindings(inv *kollectdevv1alpha1.KollectClusterInventory) []kollectdevv1alpha1.InventorySinkBinding {
	if inv == nil {
		return nil
	}

	return kollectdevv1alpha1.CollectClusterInventorySinkBindings(&inv.Spec)
}

func totalInventorySinkRefs(inv *kollectdevv1alpha1.KollectInventory) int {
	return kollectdevv1alpha1.TotalInventorySinkRefCount(&inv.Spec)
}

func totalClusterInventorySinkRefs(inv *kollectdevv1alpha1.KollectClusterInventory) int {
	return kollectdevv1alpha1.TotalClusterInventorySinkRefCount(&inv.Spec)
}

func sinkExportKey(binding kollectdevv1alpha1.InventorySinkBinding) string {
	return binding.Family + "/" + binding.Name
}

func sinkExportMinInterval(resolved *sink.ResolvedSink) *metav1.Duration {
	if resolved == nil || resolved.ExportMinInterval == nil {
		return nil
	}
	return resolved.ExportMinInterval.ExportMinInterval
}

// loadClusterInventorySink resolves a cluster inventory sink ref in its effective namespace:
// the per-ref namespace when set, otherwise the inventory default sink namespace. No cluster
// fallback resolution is attempted (ADR-0208).
func loadClusterInventorySink(
	ctx context.Context,
	c client.Client,
	sinkNS string,
	binding kollectdevv1alpha1.InventorySinkBinding,
) (*sink.ResolvedSink, error) {
	return loadResolvedSink(ctx, c, sinkBindingNamespace(binding, sinkNS), binding)
}
