package cmd

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/netbirdio/management-refactor/internals/server"
)

var newServer = func() server.Server {
	return server.NewServer()
}

func SetNewServer(fn func() server.Server) {
	newServer = fn
}

// mgmtCmd starts the management server
var mgmtCmd = &cobra.Command{
	Use:   "management",
	Short: "start NetBird Management Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		srv := newServer()

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
