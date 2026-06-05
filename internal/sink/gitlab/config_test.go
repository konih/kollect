// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package gitlab

import (
	"testing"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestConfigFromSpec(t *testing.T) {
	t.Parallel()

	t.Run("valid https endpoint", func(t *testing.T) {
		t.Parallel()

		cfg, err := ConfigFromSpec(kollectdevv1alpha1.KollectSinkSpec{
			Type:     TypeName,
			Endpoint: "https://gitlab.example.com/platform/inventory.git",
		}, nil)
		if err != nil {
			t.Fatalf("ConfigFromSpec() error = %v", err)
		}

		if cfg.Endpoint != "https://gitlab.example.com/platform/inventory.git" {
			t.Fatalf("Endpoint = %q", cfg.Endpoint)
		}
	})

	t.Run("wrong type", func(t *testing.T) {
		t.Parallel()

		if _, err := ConfigFromSpec(kollectdevv1alpha1.KollectSinkSpec{
			Type:     "git",
			Endpoint: "https://gitlab.example.com/platform/inventory.git",
		}, nil); err == nil {
			t.Fatal("expected error for git type")
		}
	})

	t.Run("missing endpoint", func(t *testing.T) {
		t.Parallel()

		if _, err := ConfigFromSpec(kollectdevv1alpha1.KollectSinkSpec{Type: TypeName}, nil); err == nil {
			t.Fatal("expected error for missing endpoint")
		}
	})

	t.Run("unsupported scheme", func(t *testing.T) {
		t.Parallel()

		if _, err := ConfigFromSpec(kollectdevv1alpha1.KollectSinkSpec{
			Type:     TypeName,
			Endpoint: "ssh://git@gitlab.example.com/platform/inventory.git",
		}, nil); err == nil {
			t.Fatal("expected error for ssh scheme")
		}
	})
}
