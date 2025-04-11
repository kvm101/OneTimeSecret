package routes

import (
	"errors"
	"net/http"
	"one_time_secret/internal/controller"
	"one_time_secret/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		basic := c.GetHeader("Authorization")
		model.ConnectDatabase()

		if basic == "" {
			c.Status(http.StatusForbidden)
			controller.RenderHTML(c, "templates/forbidden.html", nil)
			c.Abort()
		}

		auth := controller.ExtractBasic(c)
		var user model.User
		err := model.DB.First(&user, "username = ? and password = ?", auth[0], auth[1]).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(http.StatusForbidden)
			controller.RenderHTML(c, "templates/forbidden.html", nil)
			c.Abort()
		}

		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()
	route.Use(AuthMiddleware())

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
