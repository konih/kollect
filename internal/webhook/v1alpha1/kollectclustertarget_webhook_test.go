// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package webhookv1alpha1

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
)

func platformProfileRef() kollectdevv1alpha1.NamespacedObjectReference {
	return kollectdevv1alpha1.NamespacedObjectReference{
		Name:      "platform-deployments",
		Namespace: "kollect-system",
	}
}

func TestKollectClusterTargetValidator_ValidateCreate(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	_ = kollectdevv1alpha1.AddToScheme(scheme)
	v := &kollectClusterTargetValidator{client: fake.NewClientBuilder().WithScheme(scheme).Build()}

	_, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "bad"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: platformProfileRef(),
		},
	})
	if err == nil {
		t.Fatal("expected validation error for missing namespaceSelector")
	}

	_, err = v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "ok"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: platformProfileRef(),
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"team": "platform"},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected valid cluster target: %v", err)
	}
}

func TestKollectClusterTargetValidator_ValidateCreate_requiresProfileNamespace(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	_ = kollectdevv1alpha1.AddToScheme(scheme)
	v := &kollectClusterTargetValidator{client: fake.NewClientBuilder().WithScheme(scheme).Build()}

	_, err := v.ValidateCreate(context.Background(), &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "no-ns"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: kollectdevv1alpha1.NamespacedObjectReference{Name: "platform-deployments"},
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"team": "platform"},
			},
		},
	})
	if err == nil {
		t.Fatal("expected validation error for missing profileRef.namespace")
	}
}

func TestKollectClusterTargetValidator_ValidateUpdateDeletion(t *testing.T) {
	t.Parallel()

	v := &kollectClusterTargetValidator{client: fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build()}
	now := metav1.Now()
	target := &kollectdevv1alpha1.KollectClusterTarget{
		ObjectMeta: metav1.ObjectMeta{Name: "ct"},
		Spec: kollectdevv1alpha1.KollectClusterTargetSpec{
			ProfileRef: platformProfileRef(),
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"team": "platform"},
			},
		},
	}
	deleting := target.DeepCopy()
	deleting.DeletionTimestamp = &now

	if _, err := v.ValidateUpdate(context.Background(), target, deleting); err != nil {
		t.Fatalf("deletion update: %v", err)
	}

	if _, err := v.ValidateDelete(context.Background(), target); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestResolveClusterTargetProfileForWebhook(t *testing.T) {
	t.Parallel()

	scheme := runtime.NewScheme()
	if err := kollectdevv1alpha1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}

	namespacedProfile := &kollectdevv1alpha1.KollectProfile{
		ObjectMeta: metav1.ObjectMeta{Name: "team-deployments", Namespace: "kollect-system"},
		Spec: kollectdevv1alpha1.KollectProfileSpec{
			TargetGVK: kollectdevv1alpha1.GroupVersionKind{Group: "apps", Version: "v1", Kind: "StatefulSet"},
		},
	}

	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(namespacedProfile).Build()

	ref := kollectdevv1alpha1.NamespacedObjectReference{Name: "team-deployments", Namespace: "kollect-system"}
	got, err := resolveClusterTargetProfileForWebhook(context.Background(), cl, ref)
	if err != nil || got.Spec.TargetGVK.Kind != "StatefulSet" {
		t.Fatalf("namespaced profile = %#v err=%v", got, err)
	}

	missing := kollectdevv1alpha1.NamespacedObjectReference{Name: "missing", Namespace: "kollect-system"}
	if _, err := resolveClusterTargetProfileForWebhook(context.Background(), cl, missing); err == nil {
		t.Fatal("expected not found for missing profile")
	}
}
