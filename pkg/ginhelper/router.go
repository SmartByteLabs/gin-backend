package ginhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/princeparmar/gin-backend.git/pkg/config"
	"github.com/princeparmar/gin-backend.git/pkg/logger"
	"github.com/princeparmar/gin-backend.git/pkg/middleware"

	"github.com/gin-contrib/pprof"
)

func getRouter(log logger.Logger, conf *config.Config) *gin.Engine {
	log.Info("Setting up router")
	router := gin.Default()

	router.Use(httpHandlerToGinHandlerWithNext(middleware.LoggerMiddleware(log)))

	// Panic recovery middleware
	router.Use(gin.Recovery())

	// CORS middleware
	router.Use(httpHandlerToGinHandler(middleware.CORSMiddleware(conf)))

	// Add PProf routes
	pprof.Register(router)

	v1 := router.Group("api/v1")

	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	return router
}

func httpHandlerToGinHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

func httpHandlerToGinHandlerWithNext(h middleware.MiddlewareFuncWithNext) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request, c.Next)
	}
}
