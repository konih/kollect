// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package prometheus

import (
	"context"
	"fmt"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

const typeName = "prometheus"

// Backend is a stub export sink; operator metrics cover collection counts (ADR-0012).
type Backend struct{}

// NewBackend returns a registered stub that rejects export until Phase 4 ships.
func NewBackend(spec kollectdevv1alpha1.KollectSinkSpec) (*Backend, error) {
	if spec.Type != typeName {
		return nil, fmt.Errorf("expected prometheus sink, got %q", spec.Type)
	}

	return &Backend{}, nil
}

// Type returns the sink type identifier.
func (b *Backend) Type() string {
	return typeName
}

// Export is not implemented; use operator /metrics gauges instead.
func (b *Backend) Export(_ context.Context, _ []byte, _ string) error {
	return fmt.Errorf("prometheus sink export is not implemented; use operator metrics (kollect_collected_objects)")
}
