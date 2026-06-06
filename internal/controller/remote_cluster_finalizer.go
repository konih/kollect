// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/hub"
)

const remoteClusterCleanupFinalizer = "kollect.dev/remote-cluster-cleanup"

func (r *KollectRemoteClusterReconciler) ensureRemoteClusterFinalizer(
	ctx context.Context,
	rc *kollectdevv1alpha1.KollectRemoteCluster,
) error {
	return ensureFinalizer(ctx, r.Client, rc, remoteClusterCleanupFinalizer)
}

func (r *KollectRemoteClusterReconciler) finalizeRemoteClusterDeletion(
	ctx context.Context,
	rc *kollectdevv1alpha1.KollectRemoteCluster,
) (ctrl.Result, error) {
	if !containsFinalizer(rc.Finalizers, remoteClusterCleanupFinalizer) {
		return ctrl.Result{}, nil
	}

	clusterName := strings.TrimSpace(rc.Spec.ClusterName)
	if clusterName != "" && r.Store != nil {
		hub.CleanupCluster(r.Store, clusterName)
	}

	return removeFinalizerAndUpdate(ctx, r.Client, rc, remoteClusterCleanupFinalizer)
}
