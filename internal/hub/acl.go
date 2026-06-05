// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package hub

import (
	"fmt"
	"strings"
)

// ValidateClusterACL rejects reports from clusters outside the hub registration allowlist.
// An empty allowlist permits any cluster (dev / open queue); production hubs set
// KOLLECT_REMOTE_CLUSTERS from KollectHub.spec.remoteClusters (ADR-0028).
func ValidateClusterACL(cluster string, allowlist []string) error {
	cluster = strings.TrimSpace(cluster)
	if cluster == "" {
		return fmt.Errorf("hub acl: cluster is required")
	}

	if len(allowlist) == 0 {
		return nil
	}

	for _, allowed := range allowlist {
		if cluster == strings.TrimSpace(allowed) {
			return nil
		}
	}

	return fmt.Errorf("hub acl: cluster %q is not registered", cluster)
}
