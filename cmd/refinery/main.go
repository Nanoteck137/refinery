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

func init() {
	rootCmd.SetVersionTemplate(refinery.VersionTemplate(refinery.AppName))
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
