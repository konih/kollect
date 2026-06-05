// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

// Package cap holds sink capability types shared by the registry and backends
// without import cycles.
package cap

import "strings"

// Capabilities describes how a sink backend projects inventory snapshots
// (ADR-0401, ADR-0406).
type Capabilities struct {
	// Stream is true for event emitters (Kafka, NATS); false for state stores.
	Stream bool
	// SupportsDelete is true when the backend reconciles removed resources
	// (Postgres diff-and-delete). Snapshot stores replace whole files and
	// implicit deletes do not require this flag.
	SupportsDelete bool
}

// SnapshotStore is the default for Git, S3, GCS, and similar backends.
func SnapshotStore() Capabilities {
	return Capabilities{Stream: false, SupportsDelete: false}
}

// StreamEmitter is the default for Kafka and NATS event sinks.
func StreamEmitter() Capabilities {
	return Capabilities{Stream: true, SupportsDelete: false}
}

// RelationalStore is the default for Postgres upsert sinks.
func RelationalStore() Capabilities {
	return Capabilities{Stream: false, SupportsDelete: true}
}

// ExportPayload decides whether to call Backend.Export and normalizes empty snapshots.
// Backends with SupportsDelete still receive "[]" so stale relational rows are pruned.
func ExportPayload(c Capabilities, payload []byte) (export []byte, skip bool) {
	if len(payload) == 0 {
		payload = []byte("[]")
	}

	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "[]" || trimmed == "null" || trimmed == "" {
		if c.SupportsDelete {
			return []byte("[]"), false
		}

		return nil, true
	}

	return payload, false
}
