// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package gitlab

const (
	// TypeName is the KollectSink type value for GitLab git remotes.
	TypeName = "gitlab"
)

func isHTTPSEndpointScheme(scheme string) bool {
	return scheme == "https" || scheme == "http"
}
