// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package export

import (
	"cmp"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/konih/kollect/internal/collect"
)

// EnvelopePartition is one bounded export envelope slice.
type EnvelopePartition struct {
	Index     int
	Total     int
	ItemCount int
	Checksum  string
	Envelope  []byte
}

// PartitionEnvelopes splits items into bounded envelope parts.
func PartitionEnvelopes(items []collect.Item, meta Metadata, maxBytes int64) ([]EnvelopePartition, error) {
	sorted, err := stableItems(items)
	if err != nil {
		return nil, err
	}

	full, err := MarshalEnvelope(sorted, meta)
	if err != nil {
		return nil, err
	}
	if maxBytes <= 0 || int64(len(full)) <= maxBytes {
		return finalizePartitions([]EnvelopePartition{{
			ItemCount: len(sorted),
			Checksum:  EnvelopeMetaFromPayload(full).Checksum,
			Envelope:  full,
		}}), nil
	}

	if len(sorted) == 0 {
		empty, err := MarshalEnvelope([]collect.Item{}, meta)
		if err != nil {
			return nil, err
		}

		return finalizePartitions([]EnvelopePartition{{
			ItemCount: 0,
			Checksum:  EnvelopeMetaFromPayload(empty).Checksum,
			Envelope:  empty,
		}}), nil
	}

	parts := make([]EnvelopePartition, 0, 4)
	current := make([]collect.Item, 0, len(sorted))

	flushCurrent := func() error {
		payload, marshalErr := MarshalEnvelope(current, meta)
		if marshalErr != nil {
			return marshalErr
		}
		parts = append(parts, EnvelopePartition{
			ItemCount: len(current),
			Checksum:  EnvelopeMetaFromPayload(payload).Checksum,
			Envelope:  payload,
		})
		current = current[:0]

		return nil
	}

	for i := range sorted {
		candidate := append(current, sorted[i])
		payload, marshalErr := MarshalEnvelope(candidate, meta)
		if marshalErr != nil {
			return nil, marshalErr
		}
		if int64(len(payload)) <= maxBytes {
			current = candidate
			continue
		}

		if len(current) == 0 {
			return nil, fmt.Errorf("single item export envelope exceeds maxExportBytes (%d)", maxBytes)
		}
		if err := flushCurrent(); err != nil {
			return nil, err
		}

		current = append(current, sorted[i])
		singlePayload, marshalErr := MarshalEnvelope(current, meta)
		if marshalErr != nil {
			return nil, marshalErr
		}
		if int64(len(singlePayload)) > maxBytes {
			return nil, fmt.Errorf("single item export envelope exceeds maxExportBytes (%d)", maxBytes)
		}
	}

	if len(current) > 0 {
		if err := flushCurrent(); err != nil {
			return nil, err
		}
	}

	return finalizePartitions(parts), nil
}

// PartitionsChecksum returns a stable digest over part checksums.
func PartitionsChecksum(parts []EnvelopePartition) string {
	sum := sha256.New()
	for i := range parts {
		_, _ = sum.Write([]byte(parts[i].Checksum))
		_, _ = sum.Write([]byte{'\n'})
	}

	return hex.EncodeToString(sum.Sum(nil))
}

// PartitionObjectPath appends a deterministic part suffix for multipart exports.
func PartitionObjectPath(baseObjectPath string, index, total int) string {
	if total <= 1 {
		return baseObjectPath
	}

	dir, file := path.Split(baseObjectPath)
	dot := strings.LastIndex(file, ".")
	if dot <= 0 {
		return path.Join(dir, fmt.Sprintf("%s.part-%04d-of-%04d", file, index, total))
	}

	name := file[:dot]
	ext := file[dot:]

	return path.Join(dir, fmt.Sprintf("%s.part-%04d-of-%04d%s", name, index, total, ext))
}

func finalizePartitions(parts []EnvelopePartition) []EnvelopePartition {
	total := len(parts)
	for i := range parts {
		parts[i].Index = i + 1
		parts[i].Total = total
	}

	return parts
}

func stableItems(items []collect.Item) ([]collect.Item, error) {
	if len(items) == 0 {
		return []collect.Item{}, nil
	}

	type keyedItem struct {
		item collect.Item
		key  string
	}

	keyed := make([]keyedItem, 0, len(items))
	for i := range items {
		raw, err := json.Marshal(items[i])
		if err != nil {
			return nil, fmt.Errorf("marshal item sort key: %w", err)
		}
		keyed = append(keyed, keyedItem{item: items[i], key: string(raw)})
	}

	slices.SortFunc(keyed, func(a, b keyedItem) int {
		return cmp.Compare(a.key, b.key)
	})

	out := make([]collect.Item, 0, len(keyed))
	for i := range keyed {
		out = append(out, keyed[i].item)
	}

	return out, nil
}
