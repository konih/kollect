// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package hub

import (
	"testing"
	"time"

	authenticationv1 "k8s.io/api/authentication/v1"
)

func TestIngestAuthCache_hitMiss(t *testing.T) {
	t.Parallel()

	cache := newIngestAuthCache(30 * time.Second)
	user := authenticationv1.UserInfo{Username: "spoke"}

	cache.set("key", user, true)
	if allowed, ok := cache.getAllowed("key"); !ok || !allowed {
		t.Fatal("expected cache hit")
	}

	cache.set("deny", user, false)
	if allowed, ok := cache.getAllowed("deny"); !ok || allowed {
		t.Fatal("expected cached denial")
	}
}
