// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package git

import (
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGitFetchShallowAndPullRebase(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not in PATH")
	}

	remoteDir := createBareRemoteWithMainCommit(t)
	cloneDir := filepath.Join(t.TempDir(), "clone")
	runGit(t, "clone", "--branch", "main", "--single-branch", remoteDir, cloneDir)

	if err := gitFetchShallow(t.Context(), cloneDir, "main", 1, nil); err != nil {
		t.Fatalf("gitFetchShallow() error = %v", err)
	}
	if err := gitPullRebase(t.Context(), cloneDir, "main", nil); err != nil {
		t.Fatalf("gitPullRebase() error = %v", err)
	}
}
