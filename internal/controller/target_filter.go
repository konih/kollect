// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	kollectdevv1alpha1 "github.com/platformrelay/kollect/api/v1alpha1"
)

func updateTargetFilterStatus(
	target *kollectdevv1alpha1.KollectTarget,
	matched, effective []string,
	activeRules int,
) {
	target.Status.MatchedNamespaces = matched
	target.Status.EffectiveNamespaces = effective
	target.Status.ActiveResourceRules = activeRules
}

func updateClusterTargetFilterStatus(
	target *kollectdevv1alpha1.KollectClusterTarget,
	matched, effective []string,
	activeRules int,
) {
	target.Status.MatchedNamespaces = matched
	target.Status.EffectiveNamespaces = effective
	target.Status.ActiveResourceRules = activeRules
}
