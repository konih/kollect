// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package sink

import "github.com/konih/kollect/internal/sink/cap"

// Capabilities describes sink backend projection behavior (ADR-0401, ADR-0406).
type Capabilities = cap.Capabilities

// SnapshotStoreCapabilities is the default for Git, S3, GCS, and similar backends.
func SnapshotStoreCapabilities() Capabilities { return cap.SnapshotStore() }

// StreamEmitterCapabilities is the default for Kafka and NATS event sinks.
func StreamEmitterCapabilities() Capabilities { return cap.StreamEmitter() }

// RelationalStoreCapabilities is the default for Postgres upsert sinks.
func RelationalStoreCapabilities() Capabilities { return cap.RelationalStore() }

// ExportPayload decides whether to call Backend.Export for the given payload.
func ExportPayload(c Capabilities, payload []byte) (export []byte, skip bool) {
	return cap.ExportPayload(c, payload)
}
