// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package main

import (
	"testing"

	"github.com/konih/kollect/internal/collect"
)

func TestCollectCmd_helpExitsZero(t *testing.T) {
	t.Parallel()

	cmd := newCollectCmd()
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute(--help) error = %v", err)
	}
}

func TestCollectCmd_missingConfigFlagReturnsError(t *testing.T) {
	t.Parallel()

	cmd := newCollectCmd()
	cmd.SetArgs([]string{})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --config, got nil")
	}
}

func TestCollectCmd_invalidLogLevelReturnsExitFatal(t *testing.T) {
	t.Parallel()

	cmd := newCollectCmd()
	cmd.SetArgs([]string{"--config", t.TempDir(), "--log-level", "banana"})
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for invalid --log-level, got nil")
	}
}

func TestMapResultToExit_allSucceeded(t *testing.T) {
	t.Parallel()

	got := mapResultToExit(collect.RunResult{ItemCount: 5})
	if got != ExitSuccess {
		t.Errorf("got %d, want ExitSuccess", got)
	}
}

func TestMapResultToExit_someSkipped(t *testing.T) {
	t.Parallel()

	got := mapResultToExit(collect.RunResult{
		ItemCount:      3,
		SkippedTargets: []collect.SkippedTarget{{Name: "default/t1", Reason: "forbidden"}},
	})
	if got != ExitPartialFailure {
		t.Errorf("got %d, want ExitPartialFailure", got)
	}
}

func TestMapResultToExit_allFailed(t *testing.T) {
	t.Parallel()

	got := mapResultToExit(collect.RunResult{
		ItemCount:      0,
		SkippedTargets: []collect.SkippedTarget{{Name: "default/t1", Reason: "forbidden"}},
	})
	if got != ExitFatalError {
		t.Errorf("got %d, want ExitFatalError", got)
	}
}

func TestMapResultToExit_fatalErrorsAlwaysExitFatal(t *testing.T) {
	t.Parallel()

	got := mapResultToExit(collect.RunResult{
		ItemCount: 5,
		Errors:    []error{errFixture{}},
	})
	if got != ExitFatalError {
		t.Errorf("got %d, want ExitFatalError", got)
	}
}

type errFixture struct{}

func (errFixture) Error() string { return "fixture error" }
