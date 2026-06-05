// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestKollectClusterProfileValidator_ValidateCreate(t *testing.T) {
	t.Parallel()

	v := &kollectClusterProfileValidator{}

	warnings, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "bad"},
		Spec: kollectdevv1alpha1.KollectClusterProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{Version: "v1", Kind: ""},
		},
	})
	if err == nil {
		t.Fatal("expected validation error for missing kind")
	}
	if len(warnings) != 0 {
		t.Fatalf("expected no warnings on invalid profile, got %v", warnings)
	}

	warnings, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "ok"},
		Spec: kollectdevv1alpha1.KollectClusterProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{
				Group:   "apps",
				Version: "v1",
				Kind:    "Deployment",
			},
			Attributes: []kollectdevv1alpha1.AttributeSpec{
				{Name: "image", Path: "$.spec.template.spec.containers[0].image"},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected valid cluster profile: %v", err)
	}
	_ = warnings
}

func TestKollectClusterProfileValidator_ValidateUpdateDeletion(t *testing.T) {
	t.Parallel()

	v := &kollectClusterProfileValidator{}
	now := metav1.Now()
	profile := &kollectdevv1alpha1.KollectClusterProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "deployments"},
		Spec: kollectdevv1alpha1.KollectClusterProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{Version: "v1", Kind: "Deployment"},
			Attributes: []kollectdevv1alpha1.AttributeSpec{
				{Name: "name", Path: "$.metadata.name"},
			},
		},
	}
	deleting := profile.DeepCopy()
	deleting.DeletionTimestamp = &now

	if _, err := v.ValidateUpdate(context.Background(), profile, deleting); err != nil {
		t.Fatalf("deletion update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), profile); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
