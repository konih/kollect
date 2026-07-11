// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kollectdevv1alpha1 "github.com/platformrelay/kollect/api/v1alpha1"
)

func TestKollectClusterTargetReconciler_mapFunctions(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := kollectdevv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	profile := &kollectdevv1alpha1.KollectProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "platform-deployments", Namespace: "kollect-system"},
	}
	match := &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "ct-match"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: kollectdevv1alpha1.NamespacedObjectReference{Name: "platform-deployments", Namespace: "kollect-system"},
		},
	}
	other := &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "ct-other"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: kollectdevv1alpha1.NamespacedObjectReference{Name: "other", Namespace: "kollect-system"},
		},
	}

	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(match, other).Build()
	r := &KollectClusterTargetReconciler{Client: cl}

	reqs := r.mapProfileToClusterTargets(context.Background(), profile)
	if len(reqs) != 1 || reqs[0].Name != "ct-match" {
		t.Fatalf("profile map reqs = %#v", reqs)
	}

	reqs = r.mapNamespaceToClusterTargets(context.Background(), nil)
	if len(reqs) != 2 {
		t.Fatalf("namespace map reqs = %#v", reqs)
	}

	if got := r.mapProfileToClusterTargets(context.Background(), match); got != nil {
		t.Fatalf("non-profile object should return nil, got %#v", got)
	}
}
