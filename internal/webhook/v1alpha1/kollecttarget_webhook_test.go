// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestKollectTargetValidator_validateWatchMode(t *testing.T) {
	t.Parallel()

	v := &kollectTargetValidator{}

	if err := v.validate(&kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{
			ProfileRef: "deployment-images",
			WatchMode:  "",
		},
	}); err != nil {
		t.Fatalf("empty watchMode: %v", err)
	}

	if err := v.validate(&kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{
			ProfileRef: "deployment-images",
			WatchMode:  kollectdevv1alpha1.WatchModeAll,
		},
	}); err != nil {
		t.Fatalf("All watchMode: %v", err)
	}

	if err := v.validate(&kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{
			ProfileRef: "deployment-images",
			WatchMode:  kollectdevv1alpha1.WatchModeOptIn,
		},
	}); err != nil {
		t.Fatalf("OptIn watchMode: %v", err)
	}

	if err := v.validate(&kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{
			ProfileRef: "deployment-images",
			WatchMode:  "Maybe",
		},
	}); err == nil {
		t.Fatal("expected error for invalid watchMode")
	}

	if err := v.validate(&kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{ProfileRef: "team-a/deployment-images"},
	}); err == nil {
		t.Fatal("expected error for cross-namespace profileRef")
	}
}

func TestKollectTargetValidator_ValidateLifecycle(t *testing.T) {
	t.Parallel()

	v := &kollectTargetValidator{}
	target := &kollectdevv1alpha1.KollectTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "deployments", Namespace: "team-a"},
		Spec: kollectdevv1alpha1.KollectTargetSpec{
			ProfileRef: "deployment-images",
			WatchMode:  kollectdevv1alpha1.WatchModeAll,
		},
	}

	if _, err := v.ValidateCreate(context.Background(), target); err != nil {
		t.Fatalf("create: %v", err)
	}

	if _, err := v.ValidateUpdate(context.Background(), target, target); err != nil {
		t.Fatalf("update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), target); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
