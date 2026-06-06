// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package hub

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	authenticationv1 "k8s.io/api/authentication/v1"
)

const defaultIngestAuthCacheTTL = 30 * time.Second

type ingestAuthCacheEntry struct {
	user      authenticationv1.UserInfo
	allowed   bool
	expiresAt time.Time
}

type ingestAuthCache struct {
	ttl   time.Duration
	mu    sync.Mutex
	items map[string]ingestAuthCacheEntry
}

func newIngestAuthCache(ttl time.Duration) *ingestAuthCache {
	if ttl <= 0 {
		return nil
	}

	return &ingestAuthCache{ttl: ttl, items: make(map[string]ingestAuthCacheEntry)}
}

func (c *ingestAuthCache) getAllowed(key string) (allowed, ok bool) {
	if c == nil {
		return false, false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	entry, hit := c.items[key]
	if !hit || time.Now().After(entry.expiresAt) {
		return false, false
	}

	return entry.allowed, true
}

func (c *ingestAuthCache) set(key string, user authenticationv1.UserInfo, allowed bool) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = ingestAuthCacheEntry{
		user:      user,
		allowed:   allowed,
		expiresAt: time.Now().Add(c.ttl),
	}
}

func ingestTokenHash(token string) string {
	sum := sha256.Sum256([]byte(token))

	return hex.EncodeToString(sum[:])
}

func ingestAuthCacheKey(tokenHash, clusterID string) string {
	return tokenHash + "|" + clusterID
}
