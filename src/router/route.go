package router

import (
	"demo1/src/controller"
	"demo1/src/router/middleware"
	"github.com/gin-gonic/gin"
)

func CustomRoute(r *gin.Engine) {
	r.POST("/check", controller.Check)
	r.POST("/config", middleware.Validate(), func(c *gin.Context) {
		go controller.Config(c)
	})
}
