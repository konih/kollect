// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package scope

import (
	"fmt"
	"strings"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

const (
	ReasonScopeDeniedNamespace = "ScopeDeniedNamespace"
)

// ValidateDeniedNamespaces returns a violation when any namespace is in scope deniedNamespaces.
func ValidateDeniedNamespaces(scope *kollectdevv1alpha1.KollectScope, namespaces []string) error {
	if scope == nil || len(scope.Spec.DeniedNamespaces) == 0 {
		return nil
	}

	denied := make(map[string]struct{}, len(scope.Spec.DeniedNamespaces))
	for _, ns := range scope.Spec.DeniedNamespaces {
		denied[ns] = struct{}{}
	}

	for _, ns := range namespaces {
		if _, ok := denied[ns]; ok {
			return fmt.Errorf(
				"workload namespace %q is in KollectScope %q deniedNamespaces",
				ns, scope.Name,
			)
		}
	}

	return nil
}

// ValidateTargetIncludedNamespaces ensures Target intent namespaces are within Scope allow and not denied.
func ValidateTargetIncludedNamespaces(scope *kollectdevv1alpha1.KollectScope, namespaces []string) error {
	if scope == nil {
		return nil
	}

	if err := ValidateDeniedNamespaces(scope, namespaces); err != nil {
		return err
	}

	return ValidateWorkloadNamespaces(scope, namespaces)
}

// ValidateResourceRuleGVKs returns a violation when allowedGVKs is non-empty and any GVK is not listed.
func ValidateResourceRuleGVKs(
	scope *kollectdevv1alpha1.KollectScope,
	gvks []kollectdevv1alpha1.GroupVersionKind,
) error {
	if scope == nil || len(scope.Spec.AllowedGVKs) == 0 {
		return nil
	}

	for _, gvk := range gvks {
		if err := ValidateTargetGVK(scope, gvk); err != nil {
			return err
		}
	}

	return nil
}

// ValidateClusterScopeDeniedNamespaces checks namespaces against cluster scope deny list.
func ValidateClusterScopeDeniedNamespaces(scope *kollectdevv1alpha1.KollectClusterScope, namespaces []string) error {
	if scope == nil || len(scope.Spec.DeniedNamespaces) == 0 {
		return nil
	}

	denied := make(map[string]struct{}, len(scope.Spec.DeniedNamespaces))
	for _, ns := range scope.Spec.DeniedNamespaces {
		denied[ns] = struct{}{}
	}

	for _, ns := range namespaces {
		if _, ok := denied[ns]; ok {
			return fmt.Errorf(
				"workload namespace %q is in KollectClusterScope %q deniedNamespaces",
				ns, scope.Name,
			)
		}
	}

	return nil
}

// ValidateClusterScopeNamespaces ensures namespaces are within cluster scope allow and not denied.
func ValidateClusterScopeNamespaces(scope *kollectdevv1alpha1.KollectClusterScope, namespaces []string) error {
	if scope == nil {
		return nil
	}

	if err := ValidateClusterScopeDeniedNamespaces(scope, namespaces); err != nil {
		return err
	}

	if len(scope.Spec.AllowedNamespaces) == 0 {
		return nil
	}

	allowed := make(map[string]struct{}, len(scope.Spec.AllowedNamespaces))
	for _, ns := range scope.Spec.AllowedNamespaces {
		allowed[ns] = struct{}{}
	}

	for _, ns := range namespaces {
		if _, ok := allowed[ns]; !ok {
			return fmt.Errorf(
				"workload namespace %q is not in KollectClusterScope %q allowedNamespaces",
				ns, scope.Name,
			)
		}
	}

	return nil
}

// ValidateClusterScopeGVKs validates GVKs against cluster scope allowedGVKs.
func ValidateClusterScopeGVKs(
	scope *kollectdevv1alpha1.KollectClusterScope,
	gvk kollectdevv1alpha1.GroupVersionKind,
) error {
	if scope == nil || len(scope.Spec.AllowedGVKs) == 0 {
		return nil
	}

	for _, allowed := range scope.Spec.AllowedGVKs {
		if gvkMatches(allowed, gvk) {
			return nil
		}
	}

	return fmt.Errorf(
		"target GVK %s/%s/%s is not in KollectClusterScope %q allowedGVKs",
		gvk.Group, gvk.Version, gvk.Kind, scope.Name,
	)
}

// CollectRuleGVKs returns GVKs declared in resourceRules, or profile GVK when rules are empty.
func CollectRuleGVKs(
	filter kollectdevv1alpha1.CollectionFilterSpec,
	profileGVK kollectdevv1alpha1.GroupVersionKind,
) []kollectdevv1alpha1.GroupVersionKind {
	if len(filter.ResourceRules) == 0 {
		return []kollectdevv1alpha1.GroupVersionKind{profileGVK}
	}

	gvks := make([]kollectdevv1alpha1.GroupVersionKind, 0, len(filter.ResourceRules))
	for _, rule := range filter.ResourceRules {
		gvks = append(gvks, rule.GVK)
	}

	return gvks
}

// NormalizeNamespaceList trims and deduplicates namespace strings for validation.
func NormalizeNamespaceList(values []string) []string {
	if len(values) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}

		if _, ok := seen[v]; ok {
			continue
		}

		seen[v] = struct{}{}
		out = append(out, v)
	}

	return out
}
