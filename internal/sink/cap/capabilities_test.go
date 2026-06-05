// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package cap

import "testing"

func TestExportPayload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		caps       Capabilities
		payload    []byte
		wantSkip   bool
		wantExport string
	}{
		{
			name:       "relational empty snapshot exports",
			caps:       RelationalStore(),
			payload:    []byte("[]"),
			wantSkip:   false,
			wantExport: "[]",
		},
		{
			name:     "snapshot store skips empty",
			caps:     SnapshotStore(),
			payload:  []byte("[]"),
			wantSkip: true,
		},
		{
			name:       "relational nil payload exports empty",
			caps:       RelationalStore(),
			payload:    nil,
			wantSkip:   false,
			wantExport: "[]",
		},
		{
			name:       "non-empty passes through",
			caps:       RelationalStore(),
			payload:    []byte(`[{"uid":"a"}]`),
			wantSkip:   false,
			wantExport: `[{"uid":"a"}]`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			export, skip := ExportPayload(tc.caps, tc.payload)
			if skip != tc.wantSkip {
				t.Fatalf("skip = %v, want %v", skip, tc.wantSkip)
			}

			if !tc.wantSkip && string(export) != tc.wantExport {
				t.Fatalf("export = %q, want %q", export, tc.wantExport)
			}
		})
	}
}
