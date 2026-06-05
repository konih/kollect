// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestKollectRemoteClusterValidator_ValidateCreate(t *testing.T) {
	t.Parallel()

	v := &kollectRemoteClusterValidator{}

	_, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectRemoteCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "bad"},
		Spec:       kollectdevv1alpha1.KollectRemoteClusterSpec{},
	})
	if err == nil {
		t.Fatal("expected validation error for missing clusterName")
	}

	_, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectRemoteCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "spoke-a"},
		Spec:       kollectdevv1alpha1.KollectRemoteClusterSpec{ClusterName: "spoke-a"},
	})
	if err != nil {
		t.Fatalf("expected valid remote cluster: %v", err)
	}
}

func TestKollectRemoteClusterValidator_ValidateUpdateDelete(t *testing.T) {
	t.Parallel()

	v := &kollectRemoteClusterValidator{}
	rc := &kollectdevv1alpha1.KollectRemoteCluster{
		ObjectMeta: metav1.ObjectMeta{Name: "spoke-a"},
		Spec:       kollectdevv1alpha1.KollectRemoteClusterSpec{ClusterName: "spoke-a"},
	}

	if _, err := v.ValidateUpdate(context.Background(), rc, rc); err != nil {
		t.Fatalf("update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), rc); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
