// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package validation

import (
	"testing"

	kollectdevv1alpha1 "github.com/platformrelay/kollect/api/v1alpha1"
)

func TestValidateProfileMetrics(t *testing.T) {
	t.Parallel()

	base := kollectdevv1alpha1.KollectProfileSpec{
		TargetGVK: kollectdevv1alpha1.GroupVersionKind{Version: "v1", Kind: "Pod"},
		Attributes: []kollectdevv1alpha1.AttributeSpec{
			{Name: "ready", Path: "$.status.ready", Type: "int"},
		},
	}

	tests := []struct {
		name    string
		metrics []kollectdevv1alpha1.MetricSpec
		wantErr bool
	}{
		{
			name: "valid metric path references attribute",
			metrics: []kollectdevv1alpha1.MetricSpec{
				{Name: "ready_total", Path: "ready"},
			},
		},
		{
			name: "unknown attribute path rejected",
			metrics: []kollectdevv1alpha1.MetricSpec{
				{Name: "bad", Path: "missing"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			spec := base
			spec.Metrics = tt.metrics

			errs := ValidateProfileSpec(&spec)
			if tt.wantErr && len(errs) == 0 {
				t.Fatal("expected validation error, got none")
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Fatalf("unexpected validation errors: %v", errs)
			}
		})
	}
}
