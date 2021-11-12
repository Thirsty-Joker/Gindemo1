package router

import (
	"demo1/src/controller"
	"demo1/src/router/middleware"
	"github.com/gin-gonic/gin"
)

func CustomRoute(r *gin.Engine) {
	r.POST("/check", controller.Check)
	r.POST("/config/:action", middleware.Validate(), func(c *gin.Context) {
		controller.Config(c)
	})
}
