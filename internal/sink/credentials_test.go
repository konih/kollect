// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package sink

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kollectdevv1alpha1 "github.com/platformrelay/kollect/api/v1alpha1"
)

func TestResolveSecret_nilRef(t *testing.T) {
	t.Parallel()

	creds, err := ResolveSecret(context.Background(), nil, nil, "kollect-system")
	if err != nil || creds.Data != nil {
		t.Fatalf("nil ref: creds=%+v err=%v", creds, err)
	}
}

func TestResolveSecret_loadsCredentials(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "git-creds", Namespace: "kollect-system"},
		Data: map[string][]byte{
			"username": []byte("bot"),
			"password": []byte("secret-token"),
		},
	}
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()

	ref := &kollectdevv1alpha1.SecretReference{Name: "git-creds"}
	creds, err := ResolveSecret(context.Background(), cl, ref, "kollect-system")
	if err != nil {
		t.Fatalf("ResolveSecret: %v", err)
	}

	if creds.Username != "bot" || creds.Password != "secret-token" || creds.Token != "secret-token" {
		t.Fatalf("creds = %+v", creds)
	}
}

func TestResolveSecret_missingSecret(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	_, err := ResolveSecret(context.Background(), fake.NewClientBuilder().WithScheme(scheme).Build(),
		&kollectdevv1alpha1.SecretReference{Name: "missing"}, "kollect-system")
	if err == nil {
		t.Fatal("expected error for missing secret")
	}
}
