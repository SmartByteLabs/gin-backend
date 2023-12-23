package main

import (
	"github.com/princeparmar/gin-backend.git/internal/config"
	"github.com/princeparmar/gin-backend.git/internal/router"
	"github.com/princeparmar/gin-backend.git/pkg/ginhelper"
	"github.com/princeparmar/gin-backend.git/pkg/logger"
)

func main() {
	log := logger.New()
	conf := config.NewConfigFromEnv()

	ginRouterHelper := router.NewRouterHandler(log, conf)
	log.Info("Starting server...")
	ginhelper.StartServer(log, conf.App.Port, ginRouterHelper.RegisterRoute)
}
