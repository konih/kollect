// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package spoke

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kollectdevv1alpha1 "github.com/konih/kollect/api/v1alpha1"
	"github.com/konih/kollect/internal/collect"
	"github.com/konih/kollect/internal/hub"
)

func TestHTTPPublisherSendsBearerAndClusterID(t *testing.T) {
	resetPublisherCache()
	resetPublishState()
	t.Cleanup(func() {
		resetPublisherCache()
		resetPublishState()
	})

	var gotAuth, gotCluster string
	var gotBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotCluster = r.Header.Get(kollectdevv1alpha1.HeaderClusterID)
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	t.Setenv("KOLLECT_SPOKE_CLUSTER", "spoke-a")
	t.Setenv("KOLLECT_SPOKE_TOKEN", "test-token")
	t.Setenv("KOLLECT_TRANSPORT_TYPE", "http")
	t.Setenv("KOLLECT_HUB_URL", srv.URL+hub.IngestReportsPath())

	store := collect.NewStore()
	inv := &kollectdevv1alpha1.KollectInventory{
		ObjectMeta: metav1.ObjectMeta{Namespace: "team-a", Name: "inv", Generation: 1},
	}
	store.Upsert(collect.Item{
		TargetNamespace: "team-a",
		TargetName:      "t",
		Namespace:       "apps",
		Name:            "demo",
		UID:             "uid-1",
		Version:         "v1",
		Kind:            "Deployment",
	})

	if err := TryPublishReport(context.Background(), store, inv); err != nil {
		t.Fatalf("publish: %v", err)
	}

	if gotAuth != "Bearer test-token" {
		t.Fatalf("Authorization = %q", gotAuth)
	}

	if gotCluster != "spoke-a" {
		t.Fatalf("cluster header = %q", gotCluster)
	}

	var report hub.SpokeReport
	if err := json.Unmarshal(gotBody, &report); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if report.Cluster != "spoke-a" {
		t.Fatalf("report cluster = %q", report.Cluster)
	}
}
