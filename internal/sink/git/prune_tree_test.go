// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package git

import (
	"os"
	"path/filepath"
	"testing"

	billy "github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

func TestRemoveBillyOrphans_RemovesOnlyUnwrittenFiles(t *testing.T) {
	t.Parallel()

	fs := memfs.New()
	mustWriteBillyFile(t, fs, "inventory/team-a/keep.json", "{}")
	mustWriteBillyFile(t, fs, "inventory/team-a/stale.json", "{}")
	mustWriteBillyFile(t, fs, "inventory/team-b/current.json", "{}")
	mustWriteBillyFile(t, fs, "inventory/team-b/old.json", "{}")

	written := []string{
		"inventory/team-a/keep.json",
		"inventory/team-b/current.json",
	}

	if err := removeBillyOrphans(fs, written); err != nil {
		t.Fatalf("removeBillyOrphans() error = %v", err)
	}
	assertBillyExists(t, fs, "inventory/team-a/keep.json", true)
	assertBillyExists(t, fs, "inventory/team-b/current.json", true)
	assertBillyExists(t, fs, "inventory/team-a/stale.json", false)
	assertBillyExists(t, fs, "inventory/team-b/old.json", false)
}

func TestRemoveBillyOrphans_MissingManagedDirIgnored(t *testing.T) {
	t.Parallel()

	fs := memfs.New()
	written := []string{
		"inventory/team-a/keep.json",
		"inventory/missing/keep.json",
	}

	mustWriteBillyFile(t, fs, "inventory/team-a/keep.json", "{}")
	if err := removeBillyOrphans(fs, written); err != nil {
		t.Fatalf("removeBillyOrphans() error = %v", err)
	}
}

func TestRemoveDiskOrphans_RemovesOnlyUnwrittenFiles(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	mustWriteDiskFile(t, workdir, "inventory/team-a/keep.json", "{}")
	mustWriteDiskFile(t, workdir, "inventory/team-a/stale.json", "{}")
	mustWriteDiskFile(t, workdir, "inventory/team-b/current.json", "{}")

	written := []string{
		"inventory/team-a/keep.json",
		"inventory/team-b/current.json",
	}
	if err := removeDiskOrphans(workdir, written); err != nil {
		t.Fatalf("removeDiskOrphans() error = %v", err)
	}

	assertDiskExists(t, workdir, "inventory/team-a/keep.json", true)
	assertDiskExists(t, workdir, "inventory/team-b/current.json", true)
	assertDiskExists(t, workdir, "inventory/team-a/stale.json", false)
}

func TestRemoveDiskOrphans_MissingManagedDirIgnored(t *testing.T) {
	t.Parallel()

	workdir := t.TempDir()
	mustWriteDiskFile(t, workdir, "inventory/team-a/keep.json", "{}")

	written := []string{
		"inventory/team-a/keep.json",
		"inventory/missing/keep.json",
	}
	if err := removeDiskOrphans(workdir, written); err != nil {
		t.Fatalf("removeDiskOrphans() error = %v", err)
	}
}

func mustWriteBillyFile(t *testing.T, fs billy.Filesystem, rel, content string) {
	t.Helper()

	if err := fs.MkdirAll(filepath.Dir(rel), 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", rel, err)
	}
	f, err := fs.Create(rel)
	if err != nil {
		t.Fatalf("Create(%q): %v", rel, err)
	}
	if _, err := f.Write([]byte(content)); err != nil {
		_ = f.Close()
		t.Fatalf("Write(%q): %v", rel, err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("Close(%q): %v", rel, err)
	}
}

func assertBillyExists(t *testing.T, fs billy.Filesystem, rel string, want bool) {
	t.Helper()

	_, err := fs.Stat(rel)
	if want && err != nil {
		t.Fatalf("Stat(%q) error = %v, want file present", rel, err)
	}
	if !want && !os.IsNotExist(err) {
		t.Fatalf("Stat(%q) error = %v, want not-exist", rel, err)
	}
}

func mustWriteDiskFile(t *testing.T, workdir, rel, content string) {
	t.Helper()

	full := filepath.Join(workdir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o750); err != nil {
		t.Fatalf("MkdirAll(%q): %v", rel, err)
	}
	if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q): %v", rel, err)
	}
}

func assertDiskExists(t *testing.T, workdir, rel string, want bool) {
	t.Helper()

	_, err := os.Stat(filepath.Join(workdir, filepath.FromSlash(rel)))
	if want && err != nil {
		t.Fatalf("Stat(%q) error = %v, want file present", rel, err)
	}
	if !want && !os.IsNotExist(err) {
		t.Fatalf("Stat(%q) error = %v, want not-exist", rel, err)
	}
}
