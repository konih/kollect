// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package scope

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

// ClusterBinding is the active KollectClusterScope, if any.
type ClusterBinding struct {
	Enforced bool
	Scope    *kollectdevv1alpha1.KollectClusterScope
}

// LoadCluster returns the cluster scope binding. When multiple scopes exist, the oldest by name is used.
func LoadCluster(ctx context.Context, c client.Client) (ClusterBinding, error) {
	var list kollectdevv1alpha1.KollectClusterScopeList
	if err := c.List(ctx, &list); err != nil {
		return ClusterBinding{}, fmt.Errorf("list KollectClusterScope: %w", err)
	}

	if len(list.Items) == 0 {
		return ClusterBinding{}, nil
	}

	scope := list.Items[0]
	for i := 1; i < len(list.Items); i++ {
		if list.Items[i].Name < scope.Name {
			scope = list.Items[i]
		}
	}

	return ClusterBinding{Enforced: true, Scope: &scope}, nil
}
