// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package v1alpha1

import "testing"

func TestConnectionTestEnabled(t *testing.T) {
	t.Parallel()
	falseVal := false
	trueVal := true
	tests := []struct {
		name string
		spec *KollectSinkSpec
		want bool
	}{
		{name: "nil spec", spec: nil, want: true},
		{name: "unset", spec: &KollectSinkSpec{}, want: true},
		{name: "explicit true", spec: &KollectSinkSpec{ConnectionTest: &trueVal}, want: true},
		{name: "explicit false", spec: &KollectSinkSpec{ConnectionTest: &falseVal}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ConnectionTestEnabled(tt.spec); got != tt.want {
				t.Fatalf("ConnectionTestEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
