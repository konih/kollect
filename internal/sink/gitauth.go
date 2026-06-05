// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package sink

import (
	"github.com/konih/kollect/internal/sink/git"
)

const (
	secretKeyPassword = "password"
	secretKeyToken    = "token"
)

// GitAuthFromSecretData maps standard secret keys to git sink credentials.
func GitAuthFromSecretData(data map[string][]byte) git.Auth {
	auth := git.Auth{}
	if data == nil {
		return auth
	}

	if v, ok := data["username"]; ok {
		auth.Username = string(v)
	}

	if v, ok := data[secretKeyPassword]; ok {
		auth.Password = string(v)
	}

	if v, ok := data[secretKeyToken]; ok {
		auth.Token = string(v)
	}

	return auth
}
