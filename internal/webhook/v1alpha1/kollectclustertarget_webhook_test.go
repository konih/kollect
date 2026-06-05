// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestKollectClusterTargetValidator_ValidateCreate(t *testing.T) {
	t.Parallel()

	v := &kollectClusterTargetValidator{}

	_, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "bad"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: "platform-deployments",
		},
	})
	if err == nil {
		t.Fatal("expected validation error for missing namespaceSelector")
	}

	_, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "ok"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: "platform-deployments",
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"team": "platform"},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected valid cluster target: %v", err)
	}
}

func TestKollectClusterTargetValidator_ValidateUpdateDeletion(t *testing.T) {
	t.Parallel()

	v := &kollectClusterTargetValidator{}
	now := metav1.Now()
	target := &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "ct"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: "platform-deployments",
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"team": "platform"},
			},
		},
	}
	deleting := target.DeepCopy()
	deleting.DeletionTimestamp = &now

	if _, err := v.ValidateUpdate(context.Background(), target, deleting); err != nil {
		t.Fatalf("deletion update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), target); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
