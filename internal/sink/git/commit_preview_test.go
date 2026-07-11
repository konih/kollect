// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package git

import (
	"context"
	"strings"
	"testing"

	kollectdevv1alpha1 "github.com/platformrelay/kollect/api/v1alpha1"
)

func TestWithCommitContextRoundTrip(t *testing.T) {
	t.Parallel()

	cc := CommitContext{Namespace: "team-a", Name: "apps", Cluster: "cluster-a"}
	ctx := WithCommitContext(context.Background(), cc)
	got, ok := CommitContextFromContext(ctx)
	if !ok {
		t.Fatal("CommitContextFromContext() ok = false")
	}
	if got.Namespace != cc.Namespace || got.Name != cc.Name || got.Cluster != cc.Cluster {
		t.Fatalf("CommitContextFromContext() = %+v, want %+v", got, cc)
	}
}

func TestRenderCommitPreview_UsesSpecOverrides(t *testing.T) {
	t.Parallel()

	spec := kollectdevv1alpha1.KollectSinkSpec{
		Git: &kollectdevv1alpha1.GitSpec{
			CommitMessage:  "sync {{ .Namespace }}/{{ .Name }}",
			CommitBody:     "cluster={{ .Cluster }}",
			CommitTrailers: []string{"X-Test: true"},
		},
	}
	subject, body := RenderCommitPreview(spec, CommitContext{
		Namespace: "team-a",
		Name:      "apps",
		Cluster:   "cluster-a",
	})
	if subject != "sync {{ .Namespace }}/{{ .Name }}" {
		t.Fatalf("subject = %q", subject)
	}
	if !strings.Contains(body, "cluster={{ .Cluster }}") {
		t.Fatalf("body = %q", body)
	}
}
