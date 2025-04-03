package controller

import (
	"encoding/base64"
	"html/template"
	"log"
	"net/http"

	"one_time_secret/config"
	"one_time_secret/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetHome(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("/home/omni/Desktop/Projects/OneTimeSecret/internal/view/index.html"))
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		log.Println(err)
	}
}

func GetMessage(c *gin.Context) {
	var message model.Message
	base64_str := c.Param("id")
	str, err := base64.StdEncoding.DecodeString(base64_str)
	if err != nil {
		log.Println(err)
	}

	id, err := uuid.Parse(string(str))
	if err != nil {
		log.Println(err)
	}

	config.DB.First(&message, id)

	if message.Times != nil {
		if *message.Times > 1 {
			*message.Times = *message.Times - 1
			config.DB.Save(&message)
		} else {
			*message.Times = *message.Times - 1
			config.DB.Save(&message)
			config.DB.Delete(&message)
		}
	}

	if message.ID != nil {
		tmpl := template.Must(template.ParseFiles("/home/omni/Desktop/Projects/OneTimeSecret/internal/view/messages.html"))

		data := map[string]any{
			"ID":        message.ID,
			"Text":      message.Text,
			"Times":     message.Times,
			"Timestamp": message.Timestamp,
			"UserID":    message.UserID,
		}

		if err := tmpl.Execute(c.Writer, data); err != nil {
			log.Println(err)
		}

	} else {
		tmpl := template.Must(template.ParseFiles("/home/omni/Desktop/Projects/OneTimeSecret/internal/view/not_found.html"))

		c.Status(http.StatusNotFound)
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			log.Println(err)
		}
	}
}

func PostMessage(c *gin.Context) {
	var message model.Message

	err := c.BindJSON(&message)
	if err != nil {
		log.Println(err)
	}

	config.DB.Create(&message)
	if err != nil {
		log.Println(err)
	}
}

func PatchMessage(c *gin.Context) {

}

func DeleteMessage(c *gin.Context) {

}

func PostRegistration(c *gin.Context) {

}

func PostAccount(c *gin.Context) {

}

func PatchAccount(c *gin.Context) {

}

func DeleteAccount(c *gin.Context) {

}
