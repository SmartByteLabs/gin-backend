package ginhelper

import (
	"github.com/gin-gonic/gin"
	"github.com/princeparmar/9and9-templeCMS-backend.git/pkg/httphelper"
)

func Register(router *gin.RouterGroup, helper httphelper.CrudHandler, middlewares ...gin.HandlerFunc) {
	router = router.Group(helper.GetModelName())
	router.Use(middlewares...)

	router.GET("/", func(c *gin.Context) {
		helper.GetAll(c.Writer, c.Request)
	})

	router.GET("/:id", func(c *gin.Context) {
		helper.Get(c.Param("id"), c.Writer, c.Request)
	})

	router.POST("/", func(c *gin.Context) {
		helper.Create(c.Writer, c.Request)
	})

	router.PUT("/:id", func(c *gin.Context) {
		helper.Update(c.Param("id"), c.Writer, c.Request)
	})

	router.DELETE("/:id", func(c *gin.Context) {
		helper.Delete(c.Param("id"), c.Writer, c.Request)
	})
}
