// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package validation

import (
	"strings"

	k8svalidation "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

// ValidateNamespacedObjectRef checks a NamespacedObjectReference. When requireNamespace is true,
// the namespace field must be set; otherwise it is validated only when present (ADR-0208).
func ValidateNamespacedObjectRef(
	ref kollectdevv1alpha1.NamespacedObjectReference,
	fldPath *field.Path,
	requireNamespace bool,
) field.ErrorList {
	var errs field.ErrorList

	namePath := fldPath.Child("name")
	if strings.TrimSpace(ref.Name) == "" {
		errs = append(errs, field.Required(namePath, "name is required"))
	} else {
		if strings.Contains(ref.Name, "/") {
			errs = append(errs, field.Invalid(namePath, ref.Name, "name must be a single object name, not namespace/name"))
		} else if msgs := k8svalidation.IsDNS1123Subdomain(ref.Name); len(msgs) > 0 {
			errs = append(errs, field.Invalid(namePath, ref.Name, strings.Join(msgs, "; ")))
		}
	}

	nsPath := fldPath.Child("namespace")
	switch {
	case strings.TrimSpace(ref.Namespace) == "":
		if requireNamespace {
			errs = append(errs, field.Required(nsPath,
				"namespace is required on cluster-scoped references — no implicit default"))
		}
	default:
		if msgs := k8svalidation.IsDNS1123Label(ref.Namespace); len(msgs) > 0 {
			errs = append(errs, field.Invalid(nsPath, ref.Namespace, strings.Join(msgs, "; ")))
		}
	}

	return errs
}

func validateNameOnlyRef(ref string, fldPath *field.Path, kind string) field.ErrorList {
	var errs field.ErrorList

	if strings.TrimSpace(ref) == "" {
		errs = append(errs, field.Required(fldPath, kind+" must be a non-empty name"))

		return errs
	}

	if strings.Contains(ref, "/") {
		errs = append(errs, field.Invalid(fldPath, ref,
			kind+" must be a name only, not namespace/name"))
	}

	return errs
}

func validateSameNamespaceRef(ref string, fldPath *field.Path, kind string) field.ErrorList {
	var errs field.ErrorList

	if strings.TrimSpace(ref) == "" {
		errs = append(errs, field.Required(fldPath, kind+" ref is required"))

		return errs
	}

	if strings.Contains(ref, "/") {
		errs = append(errs, field.Invalid(fldPath, ref,
			kind+" ref must be a name in the same namespace, not namespace/name"))
	}

	return errs
}
