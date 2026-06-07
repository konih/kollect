// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package collect

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestEngineTargetsByGVRIndex(t *testing.T) {
	t.Parallel()

	e := &Engine{
		targets:      make(map[string]targetState),
		targetsByGVR: make(map[schema.GroupVersionResource][]string),
	}

	gvrApps := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvrCore := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}

	e.mu.Lock()
	e.targets["team-a/deploys"] = targetState{
		profile: kollectdevv1alpha1.KollectProfile{
			Spec: kollectdevv1alpha1.KollectProfileSpec{
				TargetGVK: kollectdevv1alpha1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"},
			},
		},
	}
	e.indexTargetLocked("team-a/deploys", gvrApps)
	e.targets["team-b/configs"] = targetState{
		profile: kollectdevv1alpha1.KollectProfile{
			Spec: kollectdevv1alpha1.KollectProfileSpec{
				TargetGVK: kollectdevv1alpha1.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"},
			},
		},
	}
	e.indexTargetLocked("team-b/configs", gvrCore)
	e.mu.Unlock()

	if got := len(e.targetsByGVR[gvrApps]); got != 1 {
		t.Fatalf("apps index len = %d, want 1", got)
	}

	e.mu.Lock()
	e.unindexTargetLocked("team-a/deploys", gvrApps)
	e.mu.Unlock()

	if got := len(e.targetsByGVR[gvrApps]); got != 0 {
		t.Fatalf("apps index after remove = %d, want 0", got)
	}
	if got := len(e.targetsByGVR[gvrCore]); got != 1 {
		t.Fatalf("core index len = %d, want 1", got)
	}
}
