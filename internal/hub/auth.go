// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

// +kubebuilder:rbac:groups=authentication.k8s.io,resources=tokenreviews,verbs=create

package hub

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

const (
	// IngestAuthModeKubernetes validates bearer tokens via TokenReview (ADR-0028).
	IngestAuthModeKubernetes = "kubernetes"
	// IngestAuthModeDisabled disables ingest auth (local dev / CI only).
	IngestAuthModeDisabled = "disabled"
)

// IngestAuthConfig configures hub spoke-report HTTP authentication.
type IngestAuthConfig struct {
	Mode   string
	Client kubernetes.Interface
}

// AuthDisabled reports whether ingest auth middleware should be bypassed.
func (a IngestAuthConfig) AuthDisabled() bool {
	switch a.Mode {
	case IngestAuthModeDisabled, "none":
		return true
	default:
		return false
	}
}

// Middleware wraps handlers with Kubernetes token validation when enabled.
func (a IngestAuthConfig) Middleware(next http.Handler) http.Handler {
	if a.AuthDisabled() || a.Client == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := bearerToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}

		if _, err := a.authenticate(r.Context(), token); err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)

			return
		}

		if strings.TrimSpace(r.Header.Get(kollectdevv1alpha1.HeaderClusterID)) == "" {
			http.Error(w, "missing "+kollectdevv1alpha1.HeaderClusterID+" header", http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerToken(header string) (string, error) {
	if header == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", fmt.Errorf("expected Bearer token")
	}

	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}

	return token, nil
}

func (a IngestAuthConfig) authenticate(ctx context.Context, token string) (authenticationv1.UserInfo, error) {
	review := &authenticationv1.TokenReview{
		Spec: authenticationv1.TokenReviewSpec{
			Token: token,
		},
	}

	result, err := a.Client.AuthenticationV1().TokenReviews().Create(ctx, review, metav1.CreateOptions{})
	if err != nil {
		return authenticationv1.UserInfo{}, fmt.Errorf("token review: %w", err)
	}

	if !result.Status.Authenticated {
		return authenticationv1.UserInfo{}, fmt.Errorf("token not authenticated")
	}

	return result.Status.User, nil
}
