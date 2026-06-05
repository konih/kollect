// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestKollectScopeValidator_ValidateCreate(t *testing.T) {
	t.Parallel()

	v := &kollectScopeValidator{}

	_, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectScope{
		ObjectMeta: metav1.ObjectMeta{Name: "bad", Namespace: "team-a"},
		Spec: kollectdevv1alpha1.KollectScopeSpec{
			ScopeCeilingSpec: kollectdevv1alpha1.ScopeCeilingSpec{
				AllowedGVKs: []kollectdevv1alpha1.GroupVersionKind{
					{Version: "v1", Kind: ""},
				},
			},
		},
	})
	if err == nil {
		t.Fatal("expected validation error for missing kind")
	}

	_, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectScope{
		ObjectMeta: metav1.ObjectMeta{Name: "dup-sinks", Namespace: "team-a"},
		Spec: kollectdevv1alpha1.KollectScopeSpec{
			SinkRefs: []string{"git-a", "git-a"},
		},
	})
	if err == nil {
		t.Fatal("expected validation error for duplicate sinkRefs")
	}

	_, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectScope{
		ObjectMeta: metav1.ObjectMeta{Name: "ok", Namespace: "team-a"},
		Spec: kollectdevv1alpha1.KollectScopeSpec{
			ScopeCeilingSpec: kollectdevv1alpha1.ScopeCeilingSpec{
				AllowedGVKs: []kollectdevv1alpha1.GroupVersionKind{
					{Group: "apps", Version: "v1", Kind: "Deployment"},
				},
				AllowedNamespaces: []string{"team-a"},
			},
			SinkRefs: []string{"git-inventory-demo"},
		},
	})
	if err != nil {
		t.Fatalf("expected valid scope: %v", err)
	}
}

func TestKollectScopeValidator_ValidateUpdateDelete(t *testing.T) {
	t.Parallel()

	v := &kollectScopeValidator{}
	scope := &kollectdevv1alpha1.KollectScope{
		ObjectMeta: metav1.ObjectMeta{Name: "ok", Namespace: "team-a"},
		Spec: kollectdevv1alpha1.KollectScopeSpec{
			ScopeCeilingSpec: kollectdevv1alpha1.ScopeCeilingSpec{
				AllowedGVKs: []kollectdevv1alpha1.GroupVersionKind{
					{Group: "apps", Version: "v1", Kind: "Deployment"},
				},
			},
		},
	}

	if _, err := v.ValidateUpdate(context.Background(), scope, scope); err != nil {
		t.Fatalf("update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), scope); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestValidateUniqueNonEmptyStrings(t *testing.T) {
	t.Parallel()

	if err := validateUniqueNonEmptyStrings([]string{"a", "b"}); err != nil {
		t.Fatalf("unique values: %v", err)
	}

	if err := validateUniqueNonEmptyStrings([]string{""}); err == nil {
		t.Fatal("expected empty string error")
	}

	if err := validateUniqueNonEmptyStrings([]string{"dup", "dup"}); err == nil {
		t.Fatal("expected duplicate error")
	}
}
