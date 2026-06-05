// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KollectClusterScopeSpec defines a cluster-scoped tenancy ceiling (ADR-0207).
type KollectClusterScopeSpec struct {
	ScopeCeilingSpec `json:",inline"`

	// sinkRefs lists cluster-scoped KollectSink names permitted for export from this scope.
	// +listType=set
	// +optional
	SinkRefs []string `json:"sinkRefs,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=kcscope

// KollectClusterScope is a cluster governance boundary for cluster targets and inventories.
// Static config only — no controller or status subresource (ADR-0202).
type KollectClusterScope struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitzero"`

	Spec KollectClusterScopeSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KollectClusterScopeList contains a list of KollectClusterScope.
type KollectClusterScopeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KollectClusterScope `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KollectClusterScope{}, &KollectClusterScopeList{})
}
