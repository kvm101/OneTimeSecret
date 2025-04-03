package routes

import (
	"one_time_secret/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/home", controller.GetHome)

	route.GET("/message/:id", controller.GetMessage)
	route.POST("/message", controller.PostMessage)
	route.PATCH("/message/:id", controller.PatchMessage)
	route.DELETE("/message/:id", controller.DeleteMessage)

	route.POST("/registration", controller.PostRegistration)
	route.POST("/account", controller.PostAccount)
	route.PATCH("/account", controller.PatchAccount)
	route.DELETE("/account", controller.DeleteAccount)

	return route
}
