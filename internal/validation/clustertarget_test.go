// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package validation

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestValidateClusterTargetSpec_namespaceSelectorRequired(t *testing.T) {
	t.Parallel()

	spec := kollectdevv1alpha1.KollectClusterTargetSpec{
		ProfileRef: kollectdevv1alpha1.NamespacedObjectReference{
			Name:      "platform-deployments",
			Namespace: "kollect-system",
		},
	}
	errs := ValidateClusterTargetSpec(&spec)
	if len(errs) == 0 {
		t.Fatal("expected error for missing namespaceSelector")
	}

	spec.NamespaceSelector = &metav1.LabelSelector{
		MatchLabels: map[string]string{"team": "platform"},
	}
	errs = ValidateClusterTargetSpec(&spec)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestValidateClusterTargetSpec_profileNamespaceRequired(t *testing.T) {
	t.Parallel()

	spec := kollectdevv1alpha1.KollectClusterTargetSpec{
		ProfileRef: kollectdevv1alpha1.NamespacedObjectReference{Name: "platform-deployments"},
		NamespaceSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{"team": "platform"},
		},
	}
	errs := ValidateClusterTargetSpec(&spec)
	if len(errs) == 0 {
		t.Fatal("expected error for missing profileRef.namespace")
	}
}

func TestValidateClusterTargetSpecNil(t *testing.T) {
	t.Parallel()

	if errs := ValidateClusterTargetSpec(nil); len(errs) != 0 {
		t.Fatalf("nil spec errs = %v", errs)
	}
}
