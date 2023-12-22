package main

import (
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/router"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/ginhelper"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
)

func main() {
	log := logger.New()
	conf := config.NewConfigFromEnv()

	ginRouterHelper := router.NewRouterHandler(log, conf)
	log.Info("Starting server...")
	ginhelper.StartServer(log, conf.App.Port, ginRouterHelper.RegisterRoute)
}
