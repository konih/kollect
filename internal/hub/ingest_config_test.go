// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package hub

import (
	"testing"
)

func TestIngestConfigFromEnv_customPortAndMode(t *testing.T) {
	t.Setenv(envHubIngestPort, "9090")
	t.Setenv(envHubIngestAuthMode, IngestAuthModeDisabled)

	port, mode := IngestConfigFromEnv()
	if port != 9090 || mode != IngestAuthModeDisabled {
		t.Fatalf("port=%d mode=%q", port, mode)
	}
}

func TestIngestConfigFromEnv_defaults(t *testing.T) {
	t.Setenv(envHubIngestPort, "")
	t.Setenv(envHubIngestAuthMode, "")

	port, mode := IngestConfigFromEnv()
	if port != defaultIngestPort || mode != IngestAuthModeKubernetes {
		t.Fatalf("port=%d mode=%q", port, mode)
	}
}

func TestPlatformNamespaceFromEnv(t *testing.T) {
	t.Setenv(envHubPlatformNamespace, "  kollect-system  ")
	if got := PlatformNamespaceFromEnv(); got != "kollect-system" {
		t.Fatalf("PlatformNamespaceFromEnv() = %q", got)
	}
}
