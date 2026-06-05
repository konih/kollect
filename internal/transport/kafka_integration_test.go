//go:build integration

// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package transport

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"

	"github.com/konih/kollect/internal/integrationtest"
)

func TestKafkaTransportRoundTrip(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := redpanda.Run(ctx, "docker.redpanda.com/redpandadata/redpanda:v24.2.4")
	if err != nil {
		if integrationtest.IsDockerUnavailable(err) {
			t.Skipf("docker not available: %v", err)
		}

		t.Fatalf("start redpanda: %v", err)
	}

	t.Cleanup(func() {
		_ = container.Terminate(context.Background())
	})

	broker, err := container.KafkaSeedBroker(ctx)
	if err != nil {
		t.Fatal(err)
	}

	const topic = "kollect.test"
	if err := createKafkaTopic(ctx, broker, topic); err != nil {
		t.Fatalf("create topic: %v", err)
	}

	pub, sub, err := NewTransport(Config{
		Type: TypeKafka,
		Kafka: KafkaConfig{
			Brokers: []string{broker},
			Topic:   topic,
			Group:   "kollect-test",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = Close(pub)
		_ = Close(sub)
	}()

	rtCtx, rtCancel := context.WithTimeout(ctx, 45*time.Second)
	defer rtCancel()

	if err := RoundTrip(rtCtx, struct {
		Publisher
		Subscriber
	}{pub, sub}, "inventory/default", []byte(`{"ok":true}`)); err != nil {
		t.Fatalf("round trip: %v", err)
	}
}

func createKafkaTopic(ctx context.Context, broker, topic string) error {
	conn, err := kafka.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}
