// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/collect"
)

func TestKollectRemoteClusterReconciler_deleteCleansHubStoreAndRemovesFinalizer(t *testing.T) {
	t.Parallel()

	store := collect.NewStore()
	store.Upsert(collect.Item{
		TargetNamespace: "spoke-a",
		TargetName:      "inv",
		UID:             "uid-1",
		Namespace:       "apps",
		Name:            "demo",
		Version:         "v1",
		Kind:            "Deployment",
	})

	scheme := runtime.NewScheme()
	if err := kollectdevv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme: %v", err)
	}

	now := metav1.Now()
	rc := &kollectdevv1alpha1.KollectRemoteCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "spoke-a",
			Namespace:         "kollect-system",
			Finalizers:        []string{remoteClusterCleanupFinalizer},
			DeletionTimestamp: &now,
		},
		Spec: kollectdevv1alpha1.KollectRemoteClusterSpec{ClusterName: "spoke-a"},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(rc).
		WithStatusSubresource(rc).
		Build()

	rec := &KollectRemoteClusterReconciler{
		Client: cl,
		Scheme: scheme,
		Store:  store,
	}

	_, err := rec.Reconcile(context.Background(), reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "spoke-a", Namespace: "kollect-system"},
	})
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}

	if store.TotalCount() != 0 {
		t.Fatalf("store count = %d, want 0", store.TotalCount())
	}

	var got kollectdevv1alpha1.KollectRemoteCluster
	err = cl.Get(context.Background(), types.NamespacedName{Name: "spoke-a", Namespace: "kollect-system"}, &got)
	if err == nil && containsFinalizer(got.Finalizers, remoteClusterCleanupFinalizer) {
		t.Fatalf("finalizer still present: %v", got.Finalizers)
	}
}

func TestKollectRemoteClusterReconciler_addsCleanupFinalizer(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := kollectdevv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("AddToScheme: %v", err)
	}

	rc := &kollectdevv1alpha1.KollectRemoteCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "spoke-a", Namespace: "kollect-system"},
		Spec:       kollectdevv1alpha1.KollectRemoteClusterSpec{ClusterName: "spoke-a"},
	}

	cl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(rc).
		WithStatusSubresource(rc).
		Build()

	rec := &KollectRemoteClusterReconciler{
		Client: cl,
		Scheme: scheme,
	}

	_, err := rec.Reconcile(context.Background(), reconcile.Request{
		NamespacedName: types.NamespacedName{Name: "spoke-a", Namespace: "kollect-system"},
	})
	if err != nil {
		t.Fatalf("Reconcile: %v", err)
	}

	var got kollectdevv1alpha1.KollectRemoteCluster
	if err := cl.Get(context.Background(), types.NamespacedName{Name: "spoke-a", Namespace: "kollect-system"}, &got); err != nil {
		t.Fatalf("Get remote cluster: %v", err)
	}

	if !containsFinalizer(got.Finalizers, remoteClusterCleanupFinalizer) {
		t.Fatalf("finalizers = %v, want %q", got.Finalizers, remoteClusterCleanupFinalizer)
	}
}
