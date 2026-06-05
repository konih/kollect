// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package validation

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

// ValidateClusterTargetSpec checks KollectClusterTarget spec fields.
func ValidateClusterTargetSpec(spec *kollectdevv1alpha1.KollectClusterTargetSpec) field.ErrorList {
	if spec == nil {
		return nil
	}

	var allErrs field.ErrorList
	profilePath := field.NewPath("spec").Child("profileRef")

	if strings.TrimSpace(spec.ProfileRef) == "" {
		allErrs = append(allErrs, field.Required(profilePath, "profileRef is required"))
	} else if strings.Contains(spec.ProfileRef, "/") {
		allErrs = append(allErrs, field.Invalid(profilePath, spec.ProfileRef,
			"profileRef must be a cluster profile name, not namespace/name"))
	}

	return allErrs
}

// ClusterTargetInvalid formats a validation failure for admission.
func ClusterTargetInvalid(name string, errs field.ErrorList) error {
	return fmt.Errorf("KollectClusterTarget %q is invalid: %s", name, formatErrors(errs))
}
