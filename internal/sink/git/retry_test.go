// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package git

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWithTransportRetry_succeedsFirstAttempt(t *testing.T) {
	t.Parallel()

	calls := 0
	err := withTransportRetry(t.Context(), defaultTransportRetry(), func() error {
		calls++

		return nil
	})
	if err != nil || calls != 1 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}

func TestWithTransportRetry_retriesTransient(t *testing.T) {
	t.Parallel()

	cfg := defaultTransportRetry()
	cfg.initialDelay = time.Millisecond
	cfg.maxDelay = time.Millisecond

	calls := 0
	err := withTransportRetry(t.Context(), cfg, func() error {
		calls++
		if calls < 3 {
			return errors.New("connection reset by peer")
		}

		return nil
	})
	if err != nil || calls != 3 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}

func TestWithTransportRetry_stopsOnTerminal(t *testing.T) {
	t.Parallel()

	cfg := defaultTransportRetry()
	cfg.initialDelay = time.Millisecond

	calls := 0
	err := withTransportRetry(t.Context(), cfg, func() error {
		calls++

		return errors.New("authentication required")
	})
	if err == nil || calls != 1 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}

func TestWithTransportRetry_respectsContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	cfg := defaultTransportRetry()
	cfg.initialDelay = time.Second

	err := withTransportRetry(ctx, cfg, func() error {
		return errors.New("connection reset")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v", err)
	}
}

// TestWithTransportRetry_normalizesInvalidConfig covers the two defensive defaults:
// maxAttempts < 1 is clamped to a single attempt, and a nil shouldRetry falls back to
// isTransientTransportError (so a transient error is still retried under a zero-value cfg).
func TestWithTransportRetry_zeroAttemptsRunsOnce(t *testing.T) {
	t.Parallel()

	calls := 0
	err := withTransportRetry(t.Context(), retryConfig{maxAttempts: 0}, func() error {
		calls++

		return errors.New("terminal boom")
	})
	if err == nil || calls != 1 {
		t.Fatalf("calls=%d err=%v, want exactly one attempt with error", calls, err)
	}
}

func TestWithTransportRetry_nilShouldRetryUsesDefault(t *testing.T) {
	t.Parallel()

	// factor 0 / initialDelay 0 also exercise retryBackoff's defaults, but keep the
	// backoff tiny by relying on the config defaults being large -> cancel via attempts.
	cfg := retryConfig{maxAttempts: 2, initialDelay: time.Millisecond, maxDelay: time.Millisecond}
	// shouldRetry left nil on purpose.

	calls := 0
	err := withTransportRetry(t.Context(), cfg, func() error {
		calls++
		// "connection reset" is transient per the default classifier, so the nil
		// shouldRetry must default to it and drive a second attempt.
		return errors.New("connection reset by peer")
	})
	if calls != 2 {
		t.Fatalf("calls=%d, want 2 (nil shouldRetry should default to transient classifier)", calls)
	}
	if err == nil {
		t.Fatal("expected final error after exhausting attempts")
	}
}

// TestRetryBackoff_defaultsOnNonPositiveConfig covers the initialDelay<=0 and factor<=0
// fallbacks: with a zero-value config the backoff must still be a positive duration
// derived from the package defaults rather than zero.
func TestRetryBackoff_defaultsOnNonPositiveConfig(t *testing.T) {
	t.Parallel()

	got := retryBackoff(retryConfig{}, 0)
	if got != defaultRetryInitial {
		t.Fatalf("backoff = %v, want default initial %v for attempt 0", got, defaultRetryInitial)
	}

	// attempt 1 with defaulted factor (2.0) doubles the default initial delay.
	if got := retryBackoff(retryConfig{}, 1); got != 2*defaultRetryInitial {
		t.Fatalf("backoff attempt 1 = %v, want %v", got, 2*defaultRetryInitial)
	}
}

func TestSleepContext_nonPositiveReturnsImmediately(t *testing.T) {
	t.Parallel()

	if err := sleepContext(t.Context(), 0); err != nil {
		t.Fatalf("sleepContext(0) = %v, want nil", err)
	}
	if err := sleepContext(t.Context(), -time.Second); err != nil {
		t.Fatalf("sleepContext(negative) = %v, want nil", err)
	}
}

func TestSleepContext_cancelInterruptsWait(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	if err := sleepContext(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Fatalf("sleepContext with cancelled ctx = %v, want context.Canceled", err)
	}
}

func TestRetryBackoff_capsAtMax(t *testing.T) {
	t.Parallel()

	cfg := retryConfig{initialDelay: time.Second, maxDelay: 5 * time.Second, factor: 10}
	got := retryBackoff(cfg, 3)
	if got != 5*time.Second {
		t.Fatalf("backoff = %v", got)
	}
}

func TestIsTransientTransportError(t *testing.T) {
	t.Parallel()

	cases := []struct {
		err  string
		want bool
	}{
		{"connection reset by peer", true},
		{"authentication required", false},
		{"503 service unavailable", true},
		{"pre-receive hook declined", false},
	}

	for _, tc := range cases {
		if got := isTransientTransportError(errors.New(tc.err)); got != tc.want {
			t.Fatalf("%q: got %v want %v", tc.err, got, tc.want)
		}
	}
}
