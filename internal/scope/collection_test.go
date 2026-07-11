// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package scope

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/platformrelay/kollect/api/v1alpha1"
)

func TestValidateDeniedNamespaces(t *testing.T) {
	t.Parallel()

	scope := &kollectdevv1alpha1.KollectScope{
		ObjectMeta: metav1.ObjectMeta{Name: "platform"},
		Spec: kollectdevv1alpha1.KollectScopeSpec{
			ScopeCeilingSpec: kollectdevv1alpha1.ScopeCeilingSpec{
				DeniedNamespaces: []string{"kube-system"},
			},
		},
	}

	if err := ValidateDeniedNamespaces(scope, []string{"team-a"}); err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	if err := ValidateDeniedNamespaces(scope, []string{"kube-system"}); err == nil {
		t.Fatal("expected denied namespace violation")
	}
}

func TestValidateTargetIncludedNamespaces(t *testing.T) {
	t.Parallel()

	scope := &kollectdevv1alpha1.KollectScope{
		Spec: kollectdevv1alpha1.KollectScopeSpec{
			ScopeCeilingSpec: kollectdevv1alpha1.ScopeCeilingSpec{
				AllowedNamespaces: []string{"team-a", "team-b"},
				DeniedNamespaces:  []string{"kube-system"},
			},
		},
	}

	if err := ValidateTargetIncludedNamespaces(scope, []string{"team-a"}); err != nil {
		t.Fatalf("expected allow: %v", err)
	}

	if err := ValidateTargetIncludedNamespaces(scope, []string{"team-c"}); err == nil {
		t.Fatal("expected allowlist violation")
	}

	if err := ValidateTargetIncludedNamespaces(scope, []string{"kube-system"}); err == nil {
		t.Fatal("expected deny violation")
	}
}

func TestCollectRuleGVKs_fallback(t *testing.T) {
	t.Parallel()

	profileGVK := kollectdevv1alpha1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	gvks := CollectRuleGVKs(kollectdevv1alpha1.CollectionFilterSpec{}, profileGVK)
	if len(gvks) != 1 || gvks[0] != profileGVK {
		t.Fatalf("gvks = %v", gvks)
	}
}
