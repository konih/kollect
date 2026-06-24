// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/scope"
)

func (r *KollectClusterInventoryReconciler) enforceClusterScopePolicy(
	ctx context.Context,
	inv *kollectdevv1alpha1.KollectClusterInventory,
	sinkNS string,
	bindings []kollectdevv1alpha1.InventorySinkBinding,
) (ctrl.Result, error) {
	clusterBinding, loadErr := scope.LoadCluster(ctx, r.Client)
	if loadErr != nil {
		return ctrl.Result{}, loadErr
	}

	if !clusterBinding.Enforced {
		return ctrl.Result{}, nil
	}

	for _, ns := range scope.ClusterInventoryStaticRefNamespaces(&inv.Spec, sinkNS) {
		if scopeErr := scope.ValidateClusterScopeStaticRefNamespace(clusterBinding.Scope, ns); scopeErr != nil {
			recordWarning(r.Recorder, inv, reasonSinkNamespaceDenied, scopeErr.Error())
			return r.setDegraded(ctx, inv, reasonSinkNamespaceDenied, scopeErr.Error())
		}
	}

	if scopeErr := scope.ValidateClusterInventoryClusterScopeSinkRefs(clusterBinding.Scope, bindings); scopeErr != nil {
		recordWarning(r.Recorder, inv, scopeReasonSinkDenied, scopeErr.Error())
		return r.setDegraded(ctx, inv, scopeReasonSinkDenied, scopeErr.Error())
	}

	return ctrl.Result{}, nil
}

func (r *KollectClusterTargetReconciler) loadClusterScopeBinding(
	ctx context.Context,
	ct *kollectdevv1alpha1.KollectClusterTarget,
) (scope.ClusterBinding, ctrl.Result, error) {
	clusterBinding, loadErr := scope.LoadCluster(ctx, r.Client)
	if loadErr != nil {
		return scope.ClusterBinding{}, ctrl.Result{}, loadErr
	}

	if clusterBinding.Enforced {
		if scopeErr := scope.ValidateClusterScopeStaticRefNamespace(clusterBinding.Scope, ct.Spec.ProfileRef.Namespace); scopeErr != nil {
			r.unregisterAll(ct)
			recordWarning(r.Recorder, ct, scopeReasonNSDenied, scopeErr.Error())
			if degErr := r.setDegraded(ctx, ct, scopeReasonNSDenied, scopeErr.Error()); degErr != nil {
				return clusterBinding, ctrl.Result{}, degErr
			}
			return clusterBinding, ctrl.Result{}, nil
		}
	}

	return clusterBinding, ctrl.Result{}, nil
}
