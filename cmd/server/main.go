package main

import (
	"github.com/princeparmar/gin-backend.git/pkg/config"
	"github.com/princeparmar/gin-backend.git/pkg/ginhelper"
	"github.com/princeparmar/gin-backend.git/pkg/logger"
)

func main() {
	log := logger.New()
	conf := config.NewConfigFromEnv()

	log.Info("Starting server...")
	ginhelper.StartServer(log, conf)
}
