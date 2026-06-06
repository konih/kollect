// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package operator

import "strings"

// ParseScrubKeys splits a comma-separated scrub key list from a CLI flag or Helm values.
func ParseScrubKeys(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	keys := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		normalized := strings.ToLower(part)
		if _, ok := seen[normalized]; ok {
			continue
		}

		seen[normalized] = struct{}{}
		keys = append(keys, part)
	}

	return keys
}
