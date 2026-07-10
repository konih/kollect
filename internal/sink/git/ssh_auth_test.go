// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package git

import (
	"crypto/ed25519"
	"crypto/rand"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func TestSSHAuthMethod_insecure(t *testing.T) {
	t.Parallel()

	key := testEd25519PrivateKeyPEM(t)
	auth, err := sshAuthMethod("git", key, SSHConfig{InsecureSkipVerify: true})
	if err != nil {
		t.Fatal(err)
	}
	if auth == nil || auth.HostKeyCallback == nil {
		t.Fatal("expected host key callback")
	}
}

func TestSSHAuthMethod_requiresKnownHosts(t *testing.T) {
	t.Parallel()

	key := testEd25519PrivateKeyPEM(t)
	_, err := sshAuthMethod("git", key, SSHConfig{})
	if err == nil {
		t.Fatal("expected error without known_hosts")
	}
}

func TestSSHAuthMethod_knownHosts(t *testing.T) {
	t.Parallel()

	key := testEd25519PrivateKeyPEM(t)
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	sshPub, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}

	known := []byte(knownhosts.Line([]string{"git.example"}, sshPub) + "\n")
	_, err = sshAuthMethod("git", key, SSHConfig{KnownHosts: known})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSSHAuthMethod_rejectsUnparseableKey(t *testing.T) {
	t.Parallel()

	_, err := sshAuthMethod("git", []byte("not a valid pem key"), SSHConfig{InsecureSkipVerify: true})
	if err == nil {
		t.Fatal("expected error for unparseable private key")
	}
	if !strings.Contains(err.Error(), "parse ssh private key") {
		t.Fatalf("error = %v, want parse ssh private key wrapper", err)
	}
}

func TestPublicKeysAuth_ClientConfig(t *testing.T) {
	t.Parallel()

	key := testEd25519PrivateKeyPEM(t)
	auth, err := sshAuthMethod("deploy", key, SSHConfig{InsecureSkipVerify: true})
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := auth.ClientConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.User != "deploy" {
		t.Fatalf("User = %q, want deploy", cfg.User)
	}
	if len(cfg.Auth) != 1 {
		t.Fatalf("Auth = %v, want one public-key auth method", cfg.Auth)
	}
	if cfg.HostKeyCallback == nil {
		t.Fatal("expected HostKeyCallback to be wired from the auth method")
	}
	if len(cfg.KeyExchanges) != len(defaultSSHKeyExchangeAlgorithms) {
		t.Fatalf("KeyExchanges = %v, want %v", cfg.KeyExchanges, defaultSSHKeyExchangeAlgorithms)
	}
}
