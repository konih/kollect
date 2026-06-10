// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package bigquery

import (
	"errors"
	"strings"
	"testing"

	"github.com/konih/kollect/internal/collect"
)

func TestExport_MergeErrorStopsDeleteStale(t *testing.T) {
	t.Parallel()

	exec := &fakeQueryExecutor{errs: []error{errors.New("merge failed")}}
	b := &Backend{
		cfg:      Config{Project: "proj", Dataset: "inventory", Table: "items", Cluster: "cluster-a"},
		executor: exec,
	}

	payload := mustEnvelope(t, []collect.Item{{
		TargetName: "deployments",
		UID:        "uid-1",
		Namespace:  "team-a",
	}})
	err := b.Export(t.Context(), payload, "inventory/team-a/apps.json")
	if err == nil || !strings.Contains(err.Error(), "bigquery merge upsert") {
		t.Fatalf("Export() error = %v, want wrapped merge error", err)
	}
	if len(exec.calls) != 1 {
		t.Fatalf("calls = %d, want 1 (merge only)", len(exec.calls))
	}
}

func TestExport_DeleteStaleErrorIsWrapped(t *testing.T) {
	t.Parallel()

	exec := &fakeQueryExecutor{errs: []error{nil, errors.New("delete stale failed")}}
	b := &Backend{
		cfg:      Config{Project: "proj", Dataset: "inventory", Table: "items", Cluster: "cluster-a"},
		executor: exec,
	}

	payload := mustEnvelope(t, []collect.Item{{
		TargetName: "deployments",
		UID:        "uid-1",
		Namespace:  "team-a",
	}})
	err := b.Export(t.Context(), payload, "inventory/team-a/apps.json")
	if err == nil || !strings.Contains(err.Error(), "bigquery delete stale") {
		t.Fatalf("Export() error = %v, want wrapped delete stale error", err)
	}
}

func TestExport_EmulatorInsertErrorIsWrapped(t *testing.T) {
	t.Parallel()

	exec := &fakeQueryExecutor{errs: []error{nil, errors.New("insert failed")}}
	b := &Backend{
		cfg:      Config{Project: "proj", Dataset: "inventory", Table: "items", Cluster: "cluster-a", UseEmulator: true},
		executor: exec,
	}

	payload := mustEnvelope(t, []collect.Item{{
		TargetName: "deployments",
		UID:        "uid-1",
		Namespace:  "team-a",
	}})
	err := b.Export(t.Context(), payload, "inventory/team-a/apps.json")
	if err == nil || !strings.Contains(err.Error(), "bigquery emulator insert") {
		t.Fatalf("Export() error = %v, want wrapped emulator insert error", err)
	}
}
