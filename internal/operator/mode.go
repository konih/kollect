// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package operator

import (
	"os"
	"strings"
)

const (
	// ModeCluster is the default single-cluster operator (collect + export).
	ModeCluster = "cluster"

	envMode = "KOLLECT_MODE"
)

// ResolveMode returns the effective operator mode from flag and KOLLECT_MODE env.
// Only cluster/single mode is supported; unknown values fall back to cluster.
func ResolveMode(flagValue string) string {
	if v := strings.TrimSpace(flagValue); v != "" {
		return normalizeMode(v)
	}

	if v := strings.TrimSpace(os.Getenv(envMode)); v != "" {
		return normalizeMode(v)
	}

	return ModeCluster
}

func normalizeMode(_ string) string {
	// Keep normalization behind this seam so future modes can be added without
	// changing flag and environment precedence in ResolveMode.
	return ModeCluster
}
