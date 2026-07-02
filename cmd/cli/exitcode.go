// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package main

import "github.com/konih/kollect/internal/collect"

// Process exit codes (ADR-0801).
const (
	// ExitSuccess: all targets collected and written.
	ExitSuccess = 0
	// ExitPartialFailure: at least one target was skipped (forbidden/transient/gvk-not-found)
	// but at least one item was still collected.
	ExitPartialFailure = 1
	// ExitFatalError: config invalid, cluster unreachable, or every target failed.
	ExitFatalError = 2
)

// mapResultToExit maps a single-context RunResult to a process exit code.
func mapResultToExit(r collect.RunResult) int {
	switch {
	case len(r.Errors) > 0:
		return ExitFatalError
	case len(r.SkippedTargets) > 0 && r.ItemCount == 0:
		return ExitFatalError
	case len(r.SkippedTargets) > 0:
		return ExitPartialFailure
	default:
		return ExitSuccess
	}
}
