// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package collect

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestEffectiveNamespaces_intersectionAndDeny(t *testing.T) {
	t.Parallel()

	nsMeta := map[string]NamespaceMeta{
		"team-a":      {Labels: labels.Set{"team": "a"}},
		"team-b":      {Labels: labels.Set{"team": "b"}},
		"kube-system": {},
	}

	filter := kollectdevv1alpha1.CollectionFilterSpec{
		IncludedNamespaces: []string{"team-a", "team-b", "kube-system"},
	}
	matched := MatchIntentNamespaces(filter, nil, nsMeta, NamespaceDefaults{})
	ceiling := ScopeCeiling{
		AllowedNamespaces: []string{"team-a", "team-b"},
		DeniedNamespaces:  []string{"kube-system"},
	}

	effective := EffectiveNamespaces(matched, ceiling, filter, NamespaceDefaults{})

	if len(effective) != 2 || effective[0] != "team-a" || effective[1] != "team-b" {
		t.Fatalf("effective = %v", effective)
	}
}

func TestResourceRules_unionAndFallback(t *testing.T) {
	t.Parallel()

	gvr := schema.GroupVersionResource{
		Group: "aquasecurity.github.io", Version: "v1alpha1", Resource: "vulnerabilityreports",
	}
	profile := &kollectdevv1alpha1.KollectProfile{
		Spec: kollectdevv1alpha1.KollectProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{
				Group: "aquasecurity.github.io", Version: "v1alpha1", Kind: "VulnerabilityReport",
			},
		},
	}

	ext, err := NewExtractor()
	if err != nil {
		t.Fatal(err)
	}

	rules, err := CompileResourceRules([]kollectdevv1alpha1.ResourceRule{
		{
			GVK: profile.Spec.TargetGVK,
			MatchLabels: map[string]string{
				"trivy-operator.resource.criticality": "high",
			},
		},
		{
			GVK:         profile.Spec.TargetGVK,
			MatchPolicy: "has(object.status.summary) && object.status.summary.criticalCount > 0",
		},
	}, ext.celEnv)
	if err != nil {
		t.Fatal(err)
	}

	target := &kollectdevv1alpha1.KollectTarget{}
	nsMeta := map[string]NamespaceMeta{"team-a": {}}

	labelObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"name": "r1", "namespace": "team-a",
			"labels": map[string]any{"trivy-operator.resource.criticality": "high"},
		},
	}}
	if !ResourceMatchesRules(labelObj, gvr, target, profile, rules, nsMeta) {
		t.Fatal("expected label rule match")
	}

	celObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"name": "r2", "namespace": "team-a"},
		"status":   map[string]any{"summary": map[string]any{"criticalCount": int64(2)}},
	}}
	if !ResourceMatchesRules(celObj, gvr, target, profile, rules, nsMeta) {
		t.Fatal("expected CEL rule match")
	}

	mediumObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"name": "r3", "namespace": "team-a"},
		"status":   map[string]any{"summary": map[string]any{"criticalCount": int64(0)}},
	}}
	if ResourceMatchesRules(mediumObj, gvr, target, profile, rules, nsMeta) {
		t.Fatal("expected no match")
	}

	legacyTarget := &kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "billing"},
			},
		},
	}
	legacyObj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{
			"name": "deploy", "namespace": "team-a",
			"labels": map[string]any{"app": "billing"},
		},
	}}
	if !ResourceMatchesRules(legacyObj, gvr, legacyTarget, profile, nil, nsMeta) {
		t.Fatal("expected legacy fallback match")
	}
}

func TestShouldCollect_evaluatedAfterNamespaceGates(t *testing.T) {
	t.Parallel()

	target := &kollectdevv1alpha1.KollectTarget{
		Spec: kollectdevv1alpha1.KollectTargetSpec{WatchMode: kollectdevv1alpha1.WatchModeAll},
	}
	ns := namespaceMeta{Labels: labels.Set{kollectdevv1alpha1.LabelWatch: kollectdevv1alpha1.WatchValueDisabled}}

	if ShouldCollect(labels.Set{}, ns, target) {
		t.Fatal("watch label opt-out should apply last")
	}
}
