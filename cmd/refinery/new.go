package main

import (
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use: "new <NAME>",
	Short: "Create new config",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		conf := Config{
			Owner:            "<REPO OWNER>",
			RepoName:         "<REPO NAME>",
			RegistryImage:    "<NAME OF THE FINAL IMAGE>",
			RegistryUsername: "<USERNAME FOR REGISTRY AUTH>",
			RegistryPassword: "<PASSWORD FOR REGISTRY AUTH>",
		}

		d, err := toml.Marshal(conf)
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(name + ".toml", d, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
