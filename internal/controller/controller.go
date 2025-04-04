package controller

import (
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
	str := c.Param("id")
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

func bindModelMessage(c *gin.Context) model.Message {
	var message model.Message

	err := c.BindJSON(&message)
	if err != nil {
		log.Println(err)
	}

	return message
}

func PostMessage(c *gin.Context) {
	message := bindModelMessage(c)
	if message.Text != nil {
		config.DB.Create(&message)
		return
	}

	c.Status(http.StatusNoContent)
	return
}

func PatchMessage(c *gin.Context) {
	message := bindModelMessage(c)
	config.DB.Save(&message)
}

func DeleteMessage(c *gin.Context) {
	str := c.Param("id")
	id, err := uuid.Parse(str)
	if err != nil {
		log.Println(err)
	}

	config.DB.Delete(&model.Message{}, id)
}

// func PostRegistration(c *gin.Context) {

// }

// func PostAccount(c *gin.Context) {

// }

// func PatchAccount(c *gin.Context) {

// }

// func DeleteAccount(c *gin.Context) {

// }
