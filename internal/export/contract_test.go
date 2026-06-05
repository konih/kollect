// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package export

import "testing"

func TestNormalizeSchemaVersion(t *testing.T) {
	t.Parallel()

	if got := NormalizeSchemaVersion(""); got != SchemaVersion {
		t.Fatalf("empty = %q, want %q", got, SchemaVersion)
	}

	if got := NormalizeSchemaVersion("custom"); got != "custom" {
		t.Fatalf("custom = %q", got)
	}
}

func TestValidateSchemaVersion(t *testing.T) {
	t.Parallel()

	if err := ValidateSchemaVersion(""); err != nil {
		t.Fatalf("empty default: %v", err)
	}

	if err := ValidateSchemaVersion(SchemaVersion); err != nil {
		t.Fatalf("current: %v", err)
	}

	if err := ValidateSchemaVersion("kollect.dev/v99"); err == nil {
		t.Fatal("expected error for unsupported version")
	}
}
