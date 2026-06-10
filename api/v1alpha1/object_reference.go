// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package v1alpha1

import "strings"

// NamespacedObjectReference points at a namespaced object by name and namespace (ADR-0208).
//
// Cluster-scoped reconciled kinds use this to reference namespaced static config
// (KollectProfile / family sinks). The namespace is required on cluster kinds at admission;
// callers that allow a documented default resolve it via EffectiveNamespace.
type NamespacedObjectReference struct {
	// name is the referenced object's name.
	// +required
	Name string `json:"name"`

	// namespace is the referenced object's namespace.
	// Required on cluster-scoped kinds; optional where a default namespace applies.
	// +optional
	Namespace string `json:"namespace,omitempty"`
}

// EffectiveNamespace returns the reference namespace, falling back to defaultNamespace when unset.
func (r NamespacedObjectReference) EffectiveNamespace(defaultNamespace string) string {
	if strings.TrimSpace(r.Namespace) != "" {
		return r.Namespace
	}

	return defaultNamespace
}

// IsZero reports whether the reference is empty (no name set).
func (r NamespacedObjectReference) IsZero() bool {
	return strings.TrimSpace(r.Name) == "" && strings.TrimSpace(r.Namespace) == ""
}
