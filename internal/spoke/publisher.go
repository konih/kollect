// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package spoke

import (
	"fmt"
	"os"
	"sync"

	"github.com/konih/kollect/internal/transport"
)

var (
	publisherMu   sync.Mutex
	cachedPub     transport.Publisher
	cachedKey     string
	testPublisher transport.Publisher
)

func publisherFor(cfg transport.Config) (transport.Publisher, error) {
	if testPublisher != nil {
		return testPublisher, nil
	}
	key, err := cfgCacheKey(cfg)
	if err != nil {
		return nil, err
	}

	publisherMu.Lock()
	defer publisherMu.Unlock()

	if cachedPub != nil && cachedKey == key {
		return cachedPub, nil
	}

	var pub transport.Publisher

	switch cfg.Type {
	case transport.TypeHTTP:
		if cfg.HTTP.URL == "" {
			return nil, fmt.Errorf("spoke publish: KOLLECT_HUB_URL is required when transport type is http")
		}

		cluster := os.Getenv("KOLLECT_SPOKE_CLUSTER")
		if cluster == "" {
			return nil, fmt.Errorf("spoke publish: KOLLECT_SPOKE_CLUSTER is required for http transport")
		}

		pub = NewHTTPPublisher(cfg.HTTP.URL, cluster)
	default:
		var sub transport.Subscriber
		pub, sub, err = transport.NewTransport(cfg)
		_ = sub
		if err != nil {
			return nil, err
		}
	}

	cachedPub = pub
	cachedKey = key

	return pub, nil
}

// resetPublisherCache clears cached transport (tests only).
func resetPublisherCache() {
	publisherMu.Lock()
	cachedPub = nil
	cachedKey = ""
	testPublisher = nil
	publisherMu.Unlock()
}

func cfgCacheKey(cfg transport.Config) (string, error) {
	switch cfg.Type {
	case "", transport.TypeInProcess:
		return string(transport.TypeInProcess), nil
	case transport.TypeHTTP:
		return fmt.Sprintf("http|%s", cfg.HTTP.URL), nil
	case transport.TypeRedis:
		return fmt.Sprintf("redis|%s|%s|%s", cfg.Redis.URL, cfg.Stream, cfg.Group), nil
	case transport.TypeKafka:
		return fmt.Sprintf("kafka|%v|%s|%s", cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.Group), nil
	case transport.TypeNATS:
		return fmt.Sprintf("nats|%s|%s|%s", cfg.NATS.URL, cfg.Stream, cfg.Group), nil
	default:
		return "", fmt.Errorf("spoke publish: unknown transport %q", cfg.Type)
	}
}
