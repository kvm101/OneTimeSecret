package routes

import (
	"one_time_secret/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	route := gin.Default()

	route.GET("/message/:id", controller.GetMessage)
	route.POST("/message", controller.PostMessage)
	route.PATCH("/message", controller.PatchMessage)
	route.DELETE("/message/:id", controller.DeleteMessage)

	route.POST("/account/registration", controller.PostRegistration)
	route.GET("/account", controller.GetAccount)
	route.PATCH("/account", controller.PatchAccount)
	route.DELETE("/account", controller.DeleteAccount)

	return route
}
