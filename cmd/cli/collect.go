// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Konrad Heimel

package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCollectCmd() *cobra.Command {
	flags := &collectFlags{}

	cmd := &cobra.Command{
		Use:   "collect",
		Short: "Collect Kubernetes inventory from a kubeconfig without installing the operator",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runCollect(cmd, flags)
		},
	}

	bindCollectFlags(cmd, flags)

	return cmd
}

func runCollect(cmd *cobra.Command, flags *collectFlags) error {
	if flags.config == "" {
		return fmt.Errorf("--config is required")
	}

	if _, ok := validLogLevels[flags.logLevel]; !ok {
		return fmt.Errorf("invalid --log-level %q: must be one of debug|info|warn|error", flags.logLevel)
	}

	return runCollectPipeline(cmd, flags)
}

// runCollectPipeline wires config loading, context resolution, collection, and export.
// P-001 only validates flags; the full pipeline is wired in P-005.
func runCollectPipeline(_ *cobra.Command, _ *collectFlags) error {
	return nil
}
