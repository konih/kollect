// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package gitlab

import (
	"context"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/sink/git"
)

// Backend exports inventory payloads to a GitLab git remote.
type Backend struct {
	cfg  Config
	auth git.Auth
}

// NewBackend constructs a GitLab sink backend from spec, optional resolved CA PEM, and credentials.
func NewBackend(
	spec kollectdevv1alpha1.KollectSinkSpec,
	caPEM []byte,
	auth git.Auth,
) (*Backend, error) {
	cfg, err := ConfigFromSpec(spec, caPEM)
	if err != nil {
		return nil, err
	}

	return &Backend{cfg: cfg, auth: auth}, nil
}

// Type returns the sink type identifier.
func (b *Backend) Type() string {
	return TypeName
}

// Config exposes the resolved configuration (for connection tests).
func (b *Backend) Config() Config {
	return b.cfg
}

// Export writes payload at objectPath and pushes to the configured GitLab remote.
func (b *Backend) Export(ctx context.Context, payload []byte, objectPath string) error {
	return git.Export(ctx, b.cfg.GitConfig(), b.auth, payload, objectPath)
}
