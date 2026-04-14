package main

import (
	"os"

	"github.com/nanoteck137/refinery"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     refinery.CliAppName,
	Version: refinery.Version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(refinery.VersionTemplate(refinery.CliAppName))
}
