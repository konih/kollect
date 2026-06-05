//go:build integration

// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package transport

import "strings"

func isDockerUnavailable(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "cannot connect to the docker daemon") ||
		strings.Contains(msg, "docker.sock") ||
		strings.Contains(msg, "executable file not found") ||
		strings.Contains(msg, "permission denied")
}
