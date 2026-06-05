// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package gitlab

import (
	"fmt"
	"net/url"
	"strings"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/sink/git"
)

// Config holds resolved GitLab sink settings (HTTPS git remote).
type Config struct {
	Endpoint string
	TLS      git.TLSConfig
	CABundle []byte
}

// ConfigFromSpec validates and resolves a KollectSink GitLab spec.
func ConfigFromSpec(spec kollectdevv1alpha1.KollectSinkSpec, caPEM []byte) (Config, error) {
	if spec.Type != TypeName {
		return Config{}, fmt.Errorf("expected gitlab sink, got %q", spec.Type)
	}

	if err := git.ValidateTLSSpec(spec.TLS); err != nil {
		return Config{}, err
	}

	endpoint := strings.TrimSpace(spec.Endpoint)
	if endpoint == "" {
		return Config{}, fmt.Errorf("gitlab sink requires spec.endpoint")
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return Config{}, fmt.Errorf("parse endpoint: %w", err)
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return Config{}, fmt.Errorf("gitlab endpoint must use https or http, got %q", u.Scheme)
	}

	tlsCfg, err := git.TLSConfigFromSpec(spec.TLS, caPEM)
	if err != nil {
		return Config{}, err
	}

	pem := caPEM
	if len(pem) == 0 && spec.TLS != nil {
		pem = spec.TLS.CABundle
	}

	return Config{
		Endpoint: endpoint,
		TLS:      tlsCfg,
		CABundle: pem,
	}, nil
}

// GitConfig converts GitLab settings to the shared git export config.
func (c Config) GitConfig() git.Config {
	return git.Config{
		Endpoint: c.Endpoint,
		TLS:      c.TLS,
		CABundle: c.CABundle,
	}
}
