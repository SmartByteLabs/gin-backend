package ginhelper

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/princeparmar/gin-backend.git/pkg/config"
	"github.com/princeparmar/gin-backend.git/pkg/logger"
)

func StartServer(lg logger.Logger, conf *config.Config) {
	router := getRouter(lg, conf)

	lg = lg.WithField("func", "ginhelper.StartServer")
	// Create a channel to receive OS signals
	sigs := make(chan os.Signal, 1)

	// Register the channel to receive SIGINT and SIGTERM signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	lg.Info("setting signal handler")
	// Start a goroutine that will do something when a signal is received
	go func() {
		sig := <-sigs
		if sig == syscall.SIGINT || sig == syscall.SIGTERM {
			lg.Info("Got signal, shutting down...")
			// Gracefully shutdown or cleanup here
			os.Exit(0)
		}
	}()

	// Start the server
	lg.Infof("Starting server on port %d", conf.App.Port)
	router.Run(fmt.Sprintf(":%d", conf.App.Port))
}
