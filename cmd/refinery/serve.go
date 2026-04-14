package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/nanoteck137/refinery/apis"
	"github.com/nanoteck137/refinery/config"
	"github.com/nanoteck137/refinery/core"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the api server",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile, _ := cmd.Flags().GetString("config")

		config, err := config.Load(cfgFile)
		if err != nil {
			slog.Error("failed to load config", "err", err)
			os.Exit(1)
		}

		app := core.NewBaseApp(config)

		err = app.Bootstrap()
		if err != nil {
			slog.Error("failed to bootstrap app", "err", err)
			os.Exit(1)
		}

		e, err := apis.Server(app)
		if err != nil {
			slog.Error("failed to create server", "err", err)
			os.Exit(1)
		}

		slog.Info("starting server")

		done := make(chan bool, 1)

		// listen for interrupt signal to gracefully shutdown the application
		go func() {
			sigch := make(chan os.Signal, 1)
			signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
			<-sigch

			done <- true
		}()

		go func() {
			err = e.Start(app.Config().ListenAddr)
			if err != nil {
				slog.Error("failed to start server", "err", err)
				os.Exit(1)
			}

			done <- true
		}()

		<-done

		slog.Info("stopping server")

		err = app.Shutdown()
		if err != nil {
			slog.Error("failed to shutdown app", "err", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
