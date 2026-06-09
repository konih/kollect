// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package mongodb

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/konih/kollect/internal/collect"
)

func TestItemDocument_UsesInventoryNamespaceFallback(t *testing.T) {
	t.Parallel()

	scope := exportScope{
		inventoryNamespace: "team-a",
		inventoryName:      "apps",
		cluster:            "prod-a",
	}
	ts := time.Date(2026, time.June, 10, 0, 0, 0, 0, time.UTC)
	doc, err := itemDocument(
		scope,
		collect.Item{
			UID:        "uid-1",
			TargetName: "deployments",
			Namespace:  "",
			Name:       "api",
		},
		ts,
	)
	if err != nil {
		t.Fatalf("itemDocument: %v", err)
	}

	if got := doc["resource_namespace"]; got != "team-a" {
		t.Fatalf("resource_namespace = %v, want team-a", got)
	}
	if got := doc["target_name"]; got != "deployments" {
		t.Fatalf("target_name = %v, want deployments", got)
	}
	payload, ok := doc["payload"].(map[string]any)
	if !ok {
		t.Fatalf("payload type = %T, want map[string]any", doc["payload"])
	}
	if got := payload["name"]; got != "api" {
		t.Fatalf("payload.name = %v, want api", got)
	}
}

func TestStaleDeleteFilter_DeleteAllScopeOnly(t *testing.T) {
	t.Parallel()

	scope := exportScope{
		inventoryNamespace: "team-a",
		inventoryName:      "apps",
		cluster:            "prod-a",
	}
	filter, deleteAll := staleDeleteFilter(scope, nil)
	if !deleteAll {
		t.Fatal("deleteAll = false, want true")
	}
	if _, hasNor := filter["$nor"]; hasNor {
		t.Fatalf("filter unexpectedly contains $nor: %#v", filter)
	}
	if got := filter["inventory_namespace"]; got != "team-a" {
		t.Fatalf("inventory_namespace = %v, want team-a", got)
	}
}

func TestStaleDeleteFilter_ExcludesCurrentSnapshotItems(t *testing.T) {
	t.Parallel()

	scope := exportScope{
		inventoryNamespace: "team-a",
		inventoryName:      "apps",
		cluster:            "prod-a",
	}
	items := []collect.Item{
		{TargetName: "deployments", UID: "uid-1"},
		{TargetName: "pods", UID: "uid-2"},
	}
	filter, deleteAll := staleDeleteFilter(scope, items)
	if deleteAll {
		t.Fatal("deleteAll = true, want false")
	}

	nor, ok := filter["$nor"].([]bson.M)
	if !ok {
		t.Fatalf("$nor type = %T, want []bson.M", filter["$nor"])
	}
	if len(nor) != 2 {
		t.Fatalf("$nor len = %d, want 2", len(nor))
	}
	if nor[0]["target_name"] != "deployments" || nor[0]["source_uid"] != "uid-1" {
		t.Fatalf("unexpected first $nor filter: %#v", nor[0])
	}
}

func TestNewExportScopeAndUpsertFilter(t *testing.T) {
	t.Parallel()

	scope := newExportScope("inventory/team-a/apps.json", "prod-a")
	if scope.inventoryNamespace != "team-a" || scope.inventoryName != "apps" || scope.cluster != "prod-a" {
		t.Fatalf("scope = %#v", scope)
	}

	filter := upsertFilter(scope, collect.Item{TargetName: "deployments", UID: "uid-1"})
	if got := filter["inventory_namespace"]; got != "team-a" {
		t.Fatalf("inventory_namespace = %v, want team-a", got)
	}
	if got := filter["target_name"]; got != "deployments" {
		t.Fatalf("target_name = %v, want deployments", got)
	}
	if got := filter["source_uid"]; got != "uid-1" {
		t.Fatalf("source_uid = %v, want uid-1", got)
	}
}
