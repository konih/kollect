// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package operator

import (
	"reflect"
	"testing"
)

func TestParseScrubKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		raw  string
		want []string
	}{
		{raw: "", want: nil},
		{raw: "  ", want: nil},
		{raw: "internalId", want: []string{"internalId"}},
		{raw: "foo, bar ,foo", want: []string{"foo", "bar"}},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			t.Parallel()

			got := ParseScrubKeys(tt.raw)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseScrubKeys(%q) = %#v, want %#v", tt.raw, got, tt.want)
			}
		})
	}
}
