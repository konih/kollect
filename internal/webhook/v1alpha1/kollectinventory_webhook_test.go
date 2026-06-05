// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestKollectInventoryValidator_ValidateCreate(t *testing.T) {
	t.Parallel()

	v := &kollectInventoryValidator{}

	_, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectInventory{
		ObjectMeta: metav1.ObjectMeta{Name: "bad", Namespace: "team-a"},
		Spec: kollectdevv1alpha1.KollectInventorySpec{
			SinkRefs: []string{"other-ns/sink"},
		},
	})
	if err == nil {
		t.Fatal("expected validation error for cross-namespace sinkRef")
	}

	_, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectInventory{
		ObjectMeta: metav1.ObjectMeta{Name: "ok", Namespace: "team-a"},
		Spec:       kollectdevv1alpha1.KollectInventorySpec{},
	})
	if err != nil {
		t.Fatalf("expected valid inventory: %v", err)
	}
}

func TestKollectInventoryValidator_ValidateUpdateDeletion(t *testing.T) {
	t.Parallel()

	v := &kollectInventoryValidator{}
	now := metav1.Now()
	inv := &kollectdevv1alpha1.KollectInventory{
		ObjectMeta: metav1.ObjectMeta{Name: "inv", Namespace: "team-a"},
		Spec:       kollectdevv1alpha1.KollectInventorySpec{},
	}
	deleting := inv.DeepCopy()
	deleting.DeletionTimestamp = &now

	if _, err := v.ValidateUpdate(context.Background(), inv, deleting); err != nil {
		t.Fatalf("deletion update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), inv); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
