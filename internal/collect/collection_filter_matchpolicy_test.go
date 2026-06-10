// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package collect

import (
	"testing"

	"github.com/google/cel-go/cel"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func matchPolicyEnv(t *testing.T) *cel.Env {
	t.Helper()

	env, err := cel.NewEnv(cel.Variable("object", cel.DynType))
	if err != nil {
		t.Fatalf("cel.NewEnv: %v", err)
	}

	return env
}

func TestEvalMatchPolicy_boolResult(t *testing.T) {
	t.Parallel()

	env := matchPolicyEnv(t)
	obj := &unstructured.Unstructured{Object: map[string]any{
		"metadata": map[string]any{"name": "keep"},
	}}

	for expr, want := range map[string]bool{
		`object.metadata.name == "keep"`:  true,
		`object.metadata.name == "other"`: false,
	} {
		prog, err := compileMatchPolicy(env, expr)
		if err != nil {
			t.Fatalf("compileMatchPolicy(%q): %v", expr, err)
		}

		got, err := evalMatchPolicy(prog, obj)
		if err != nil {
			t.Fatalf("evalMatchPolicy(%q): %v", expr, err)
		}
		if got != want {
			t.Fatalf("evalMatchPolicy(%q) = %v, want %v", expr, got, want)
		}
	}
}

// A matchPolicy expression that does not evaluate to a bool must be rejected
// rather than silently coerced — otherwise a misconfigured rule could include
// or drop every object.
func TestEvalMatchPolicy_nonBoolResultErrors(t *testing.T) {
	t.Parallel()

	env := matchPolicyEnv(t)
	obj := &unstructured.Unstructured{Object: map[string]any{}}

	for _, expr := range []string{`1 + 1`, `"truthy"`} {
		prog, err := compileMatchPolicy(env, expr)
		if err != nil {
			t.Fatalf("compileMatchPolicy(%q): %v", expr, err)
		}

		if _, err := evalMatchPolicy(prog, obj); err == nil {
			t.Fatalf("evalMatchPolicy(%q) expected non-bool error", expr)
		}
	}
}
