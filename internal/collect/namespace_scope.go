// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package collect

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

// namespaceMatchesSelector reports whether namespace labels satisfy the target selector.
// A nil selector matches all namespaces.
func namespaceMatchesSelector(selector *metav1.LabelSelector, nsLabels labels.Set) bool {
	if selector == nil {
		return true
	}

	sel, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return false
	}

	return sel.Matches(nsLabels)
}

// MatchedNamespacesForTarget returns workload namespaces matched by the target selector.
func (e *Engine) MatchedNamespacesForTarget(targetNamespace, targetName string) []string {
	e.mu.RLock()
	st, ok := e.targets[targetKey(targetNamespace, targetName)]
	e.mu.RUnlock()
	if !ok {
		return nil
	}

	return e.matchedNamespacesForTarget(&st.target)
}

func (e *Engine) matchedNamespacesForTarget(target *kollectdevv1alpha1.KollectTarget) []string {
	e.nsMu.RLock()
	nsMeta := e.nsMeta
	defaults := e.defaults
	e.nsMu.RUnlock()

	matched := MatchIntentNamespaces(
		target.Spec.CollectionFilterSpec,
		target.Spec.NamespaceSelector,
		namespaceMetaMapToFilter(nsMeta),
		defaults,
	)

	key := targetKey(target.Namespace, target.Name)
	e.mu.RLock()
	st, ok := e.targets[key]
	e.mu.RUnlock()
	if ok && len(st.effectiveNamespaces) > 0 {
		out := make([]string, 0, len(st.effectiveNamespaces))
		for ns := range st.effectiveNamespaces {
			out = append(out, ns)
		}

		return sortedUniqueStrings(out)
	}

	return matched
}

// watchNamespaceForGVR returns a single namespace for a scoped informer when every active
// target for the GVR agrees on exactly one watched namespace; otherwise metav1.NamespaceAll.
func (e *Engine) watchNamespaceForGVR(gvr schema.GroupVersionResource) string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var single string
	seen := false

	for _, st := range e.targets {
		if gvrFromProfile(st.profile.Spec.TargetGVK) != gvr {
			continue
		}

		namespaces := e.matchedNamespacesForTarget(&st.target)
		if len(namespaces) != 1 {
			return metav1.NamespaceAll
		}

		if !seen {
			single = namespaces[0]
			seen = true

			continue
		}

		if namespaces[0] != single {
			return metav1.NamespaceAll
		}
	}

	if !seen {
		return metav1.NamespaceAll
	}

	return single
}
