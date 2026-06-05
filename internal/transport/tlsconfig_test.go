// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package transport

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTLSSettingsFromEnv(t *testing.T) {
	t.Setenv("KOLLECT_TRANSPORT_TLS_INSECURE_SKIP_VERIFY", "true")
	t.Setenv("KOLLECT_TRANSPORT_TLS_CA_FILE", "/tmp/ca.pem")

	got := TLSSettingsFromEnv()
	if !got.InsecureSkipVerify {
		t.Fatal("expected insecure skip verify")
	}

	if got.CAFile != "/tmp/ca.pem" {
		t.Fatalf("ca file = %q", got.CAFile)
	}
}

func TestTLSSettingsClientConfig_insecure(t *testing.T) {
	cfg, err := TLSSettings{InsecureSkipVerify: true}.ClientConfig()
	if err != nil {
		t.Fatalf("ClientConfig: %v", err)
	}

	if cfg == nil || !cfg.InsecureSkipVerify {
		t.Fatal("expected insecure tls config")
	}
}

func TestTLSSettingsClientConfig_requiresBothClientFiles(t *testing.T) {
	dir := t.TempDir()
	cert := filepath.Join(dir, "client.crt")
	if err := os.WriteFile(cert, []byte("not-a-cert"), 0o600); err != nil {
		t.Fatal(err)
	}

	_, err := TLSSettings{ClientCertFile: cert}.ClientConfig()
	if err == nil {
		t.Fatal("expected error when client key missing")
	}
}

func TestTLSSettingsDisabled(t *testing.T) {
	cfg, err := TLSSettings{}.ClientConfig()
	if err != nil {
		t.Fatalf("ClientConfig: %v", err)
	}

	if cfg != nil {
		t.Fatalf("expected nil config, got %#v", cfg)
	}
}
