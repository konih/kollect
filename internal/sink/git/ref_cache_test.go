// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package git

import (
	"errors"
	"testing"
	"time"
)

func TestRefCache_ttl(t *testing.T) {
	t.Parallel()

	cache := newRefCache(30 * time.Millisecond)
	key := "test-key"
	errSentinel := errors.New("probe failed")

	cache.set(key, errSentinel)
	if ok, got := cache.get(key); !ok || got != errSentinel {
		t.Fatalf("expected cached error, got %v ok=%v", got, ok)
	}

	time.Sleep(40 * time.Millisecond)
	if ok, _ := cache.get(key); ok {
		t.Fatal("expected cache expiry")
	}
}

func TestNewRefCache_nonPositiveTTLDefaults(t *testing.T) {
	t.Parallel()

	cache := newRefCache(0)
	if cache.ttl != defaultRefCacheTTL {
		t.Fatalf("ttl = %v, want default %v", cache.ttl, defaultRefCacheTTL)
	}
}

func TestRefCacheKey_stable(t *testing.T) {
	t.Parallel()

	a := refCacheKey("https://example.com/r.git", Auth{Username: "u", Token: "t"})
	b := refCacheKey("https://example.com/r.git", Auth{Username: "u", Token: "t"})
	if a != b {
		t.Fatalf("keys differ: %q vs %q", a, b)
	}
}

// TestRefCacheKey_distinguishesCredentialMaterial guards the collision-avoidance
// contract: same endpoint + username but different credential material (token vs
// password vs ssh key) must produce distinct cache keys, so a probe result for one
// auth mode is never served from another mode's cached entry.
func TestRefCacheKey_distinguishesCredentialMaterial(t *testing.T) {
	t.Parallel()

	const endpoint = "https://example.com/r.git"
	tokenKey := refCacheKey(endpoint, Auth{Username: "u", Token: "t"})
	passwordKey := refCacheKey(endpoint, Auth{Username: "u", Password: "p"})
	sshKey := refCacheKey(endpoint, Auth{Username: "u", SSHPrivateKey: []byte("KEY")})
	noneKey := refCacheKey(endpoint, Auth{Username: "u"})

	keys := map[string]string{
		"token":    tokenKey,
		"password": passwordKey,
		"ssh":      sshKey,
		"none":     noneKey,
	}

	seen := make(map[string]string, len(keys))
	for name, k := range keys {
		if other, dup := seen[k]; dup {
			t.Fatalf("cache key collision: %q and %q hash to the same key %q", name, other, k)
		}
		seen[k] = name
	}
}

// TestRefCacheKey_variesByAuthType ensures the AuthType field participates in the
// key so token vs ssh admission modes to the same endpoint don't share cache state.
func TestRefCacheKey_variesByAuthType(t *testing.T) {
	t.Parallel()

	const endpoint = "https://example.com/r.git"
	tokenAuth := refCacheKey(endpoint, Auth{AuthType: AuthTypeToken, Token: "t"})
	sshAuth := refCacheKey(endpoint, Auth{AuthType: AuthTypeSSH, Token: "t"})
	if tokenAuth == sshAuth {
		t.Fatalf("expected distinct keys for differing AuthType, both = %q", tokenAuth)
	}
}
