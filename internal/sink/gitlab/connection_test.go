// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package gitlab

import (
	"context"
	"testing"
)

func TestTestConnection_invalidURL(t *testing.T) {
	t.Parallel()

	err := TestConnection(context.Background(), Config{Endpoint: "://bad"})
	if err == nil {
		t.Fatal("expected error for invalid endpoint")
	}
}
