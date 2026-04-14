package main

import (
	"os"

	"github.com/nanoteck137/refinery"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     refinery.AppName,
	Version: refinery.Version,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(refinery.VersionTemplate(refinery.AppName))

	rootCmd.PersistentFlags().StringP("config", "c", "", "Config File")
}
