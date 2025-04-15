package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"management/internal/server"
	"management/pkg/logging"
)

var log = logging.LoggerForThisPackage()

// mgmtCmd starts the management server
var mgmtCmd = &cobra.Command{
	Use:   "management",
	Short: "start NetBird Management Server",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := logging.Init("logging.yaml")
		if err != nil {
			log.Fatalf("Failed to init logging: %v", err)
		}

		srv := server.NewServer()

		go func() {
			log.Info("Starting server on :8080")
			if err := srv.Start(); err != nil {
				log.Fatalf("Server error: %v", err)
			}
		}()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
		<-stopChan

		log.Info("Shutting down server...")
		if err := srv.Stop(); err != nil {
			log.Errorf("Error stopping server: %v", err)
		}
		log.Info("Server stopped gracefully.")

		return nil
	},
}

func init() {
	// Attach serveCmd to the rootCmd
	rootCmd.AddCommand(mgmtCmd)
}
