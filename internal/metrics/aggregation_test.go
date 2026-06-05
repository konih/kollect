// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package metrics

import "testing"

func TestMetricPathSpecStub(t *testing.T) {
	t.Parallel()

	spec := MetricPathSpec{
		Name:   "not_after_days",
		Path:   `status.notAfter`,
		Labels: []string{"namespace", "name"},
	}
	if spec.Name == "" || spec.Path == "" || len(spec.Labels) != 2 {
		t.Fatal("metric path spec should carry name, path, and optional labels")
	}
}
