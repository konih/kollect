// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"
	"strings"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestResolveClusterTargetProfile_namespacedProfile(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := kollectdevv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("scheme: %v", err)
	}

	profile := &kollectdevv1alpha1.KollectProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "deployments", Namespace: "kollect-system"},
		Spec: kollectdevv1alpha1.KollectProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{Version: "v1", Kind: "Deployment"},
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(profile).Build()
	ref := kollectdevv1alpha1.NamespacedObjectReference{Name: "deployments", Namespace: "kollect-system"}
	got, err := resolveClusterTargetProfile(context.Background(), c, ref)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}

	if got.Spec.TargetGVK.Kind != "Deployment" {
		t.Fatalf("kind = %q", got.Spec.TargetGVK.Kind)
	}
}

func TestResolveClusterTargetProfile_notFoundNoFallback(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := kollectdevv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatalf("scheme: %v", err)
	}

	// Profile exists in a different namespace — there is no implicit fallback (ADR-0208).
	profile := &kollectdevv1alpha1.KollectProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "argo", Namespace: "kollect-system"},
		Spec: kollectdevv1alpha1.KollectProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{Version: "v1alpha1", Kind: "Application"},
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(profile).Build()
	ref := kollectdevv1alpha1.NamespacedObjectReference{Name: "argo", Namespace: "team-a"}
	_, err := resolveClusterTargetProfile(context.Background(), c, ref)
	if err == nil {
		t.Fatalf("expected not-found error for cross-namespace ref")
	}
	if !strings.Contains(err.Error(), "team-a") {
		t.Fatalf("error should name the declared namespace, got %q", err.Error())
	}
}
