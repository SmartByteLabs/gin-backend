package router

import (
	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/internal/config"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/logger"
)

type gitRouterHandler struct {
	log  logger.Logger
	conf *config.Config
}

func NewRouterHandler(log logger.Logger, conf *config.Config) *gitRouterHandler {
	return &gitRouterHandler{
		log:  log,
		conf: conf,
	}
}

func (h *gitRouterHandler) RegisterRoute(*gin.Engine) {
}
