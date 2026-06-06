// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package inventory

import (
	"testing"

	"github.com/konih/kollect/internal/collect"
)

func TestPaginateItems_maxLimitCap(t *testing.T) {
	t.Parallel()

	items := make([]collect.Item, 6000)
	for i := range items {
		items[i] = collect.Item{UID: "uid"}
	}

	page, meta := paginateItems(items, 10000, 0)
	if meta.Limit != maxPageLimit {
		t.Fatalf("limit = %d, want %d", meta.Limit, maxPageLimit)
	}
	if len(page) != maxPageLimit {
		t.Fatalf("page len = %d, want %d", len(page), maxPageLimit)
	}
}
