// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package hub

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authenticationv1 "k8s.io/api/authentication/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func TestIngestAuthMiddlewareKubernetesSuccess(t *testing.T) {
	t.Parallel()

	client := fake.NewSimpleClientset() //nolint:staticcheck // SimpleClientset sufficient for auth unit test
	client.PrependReactor("create", "tokenreviews", func(action k8stesting.Action) (bool, runtime.Object, error) {
		review := action.(k8stesting.CreateAction).GetObject().(*authenticationv1.TokenReview)
		review.Status = authenticationv1.TokenReviewStatus{
			Authenticated: true,
			User: authenticationv1.UserInfo{
				Username: "system:serviceaccount:spoke-a:kollect-spoke",
			},
		}

		return true, review, nil
	})

	called := false
	handler := IngestAuthConfig{
		Mode:   IngestAuthModeKubernetes,
		Client: client,
	}.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, ingestReportsPath, nil)
	req.Header.Set("Authorization", "Bearer good-token")
	req.Header.Set(kollectdevv1alpha1.HeaderClusterID, "spoke-a")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	if !called {
		t.Fatal("handler not called")
	}
}

func TestIngestAuthMiddlewareMissingClusterID(t *testing.T) {
	t.Parallel()

	client := fake.NewSimpleClientset() //nolint:staticcheck // SimpleClientset sufficient for auth unit test
	client.PrependReactor("create", "tokenreviews", func(action k8stesting.Action) (bool, runtime.Object, error) {
		review := action.(k8stesting.CreateAction).GetObject().(*authenticationv1.TokenReview)
		review.Status = authenticationv1.TokenReviewStatus{Authenticated: true}

		return true, review, nil
	})

	handler := IngestAuthConfig{
		Mode:   IngestAuthModeKubernetes,
		Client: client,
	}.Middleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, ingestReportsPath, nil)
	req.Header.Set("Authorization", "Bearer good-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestIngestAuthMiddlewareDisabledBypasses(t *testing.T) {
	t.Parallel()

	handler := IngestAuthConfig{Mode: IngestAuthModeDisabled}.Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)

	req := httptest.NewRequest(http.MethodPost, ingestReportsPath, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
}
