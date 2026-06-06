// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/sink"
)

func testNameSuffix() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func newPostgresRecordingRegistry(recorder sink.Backend) *sink.Registry {
	reg := sink.NewRegistry()
	reg.Register(kollectdevv1alpha1.SinkTypePostgres, func(
		_ kollectdevv1alpha1.KollectSinkSpec, _ sink.BuildContext,
	) (sink.Backend, error) {
		return recorder, nil
	})

	return reg
}

func newPostgresRecordingRegistryWithSelector(
	selectBackend func(kollectdevv1alpha1.KollectSinkSpec) (sink.Backend, error),
) *sink.Registry {
	reg := sink.NewRegistry()
	reg.Register(kollectdevv1alpha1.SinkTypePostgres, func(
		spec kollectdevv1alpha1.KollectSinkSpec, _ sink.BuildContext,
	) (sink.Backend, error) {
		return selectBackend(spec)
	})

	return reg
}

func createPostgresSinkFixtures(sinkName, secretName, ns string) (*kollectdevv1alpha1.KollectSink, *corev1.Secret) {
	sinkObj := &kollectdevv1alpha1.KollectSink{
		ObjectMeta: metav1.ObjectMeta{Name: sinkName, Namespace: ns},
		Spec: kollectdevv1alpha1.KollectSinkSpec{
			Type: kollectdevv1alpha1.SinkTypePostgres,
			Postgres: &kollectdevv1alpha1.PostgresSpec{
				DatabaseRef: &kollectdevv1alpha1.SecretReference{Name: secretName},
				Table:       "inventory_items",
			},
		},
	}

	pgSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: ns},
		Data:       map[string][]byte{"dsn": []byte("postgres://example")},
	}

	return sinkObj, pgSecret
}

// listFailInventoryClient injects List errors for KollectInventoryList (map-handler tests).
type listFailInventoryClient struct {
	client.Client
	listErr error
}

func (c *listFailInventoryClient) List(
	ctx context.Context,
	list client.ObjectList,
	opts ...client.ListOption,
) error {
	if c.listErr != nil {
		if _, ok := list.(*kollectdevv1alpha1.KollectInventoryList); ok {
			return c.listErr
		}
	}

	return c.Client.List(ctx, list, opts...)
}

func newListFailInventoryClient(base client.Client, listErr error) client.Client {
	if listErr == nil {
		listErr = errors.New("simulated inventory list failure")
	}

	return &listFailInventoryClient{Client: base, listErr: listErr}
}

func createRemoteClusterWithRequiredStatus(ctx context.Context, rc *kollectdevv1alpha1.KollectRemoteCluster) error {
	return k8sClient.Create(ctx, rc)
}
