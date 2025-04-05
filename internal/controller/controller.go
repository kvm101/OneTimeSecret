package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"one_time_secret/config"
	"one_time_secret/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetHome(c *gin.Context) {
	tmpl := template.Must(template.ParseFiles("/home/omni/Desktop/Projects/OneTimeSecret/internal/view/index.html"))
	if err := tmpl.Execute(c.Writer, nil); err != nil {
		log.Println(err)
		return
	}

}

func extractBasic(c *gin.Context) []string {
	b64_auth, _ := strings.CutPrefix(c.Request.Header.Get("Authorization"), "Basic ")

	auth, err := base64.StdEncoding.DecodeString(b64_auth)
	if err != nil {
		log.Println(err)
		return nil
	}

	arr_data := strings.Split(string(auth), ":")

	sum := sha256.Sum256([]byte(arr_data[1]))
	arr_data[1] = fmt.Sprintf("%x", sum)

	return arr_data
}

func GetMessage(c *gin.Context) {
	var message model.Message
	str := c.Param("id")
	id, err := uuid.Parse(string(str))
	if err != nil {
		log.Println(err)
		return
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
			return
		}

	} else {
		tmpl := template.Must(template.ParseFiles("/home/omni/Desktop/Projects/OneTimeSecret/internal/view/not_found.html"))

		c.Status(http.StatusNotFound)
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			log.Println(err)
			return
		}
	}
}

func bindModelMessage(c *gin.Context) model.Message {
	arr_data := extractBasic(c)

	var message model.Message
	var user model.User

	err := c.BindJSON(&message)
	log.Println(message)
	if err != nil {
		log.Println(err)
	}

	config.DB.Find(&user, "username = ? AND password = ?", arr_data[0], arr_data[1])

	message.UserID = user.ID

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
		return
	}

	config.DB.Delete(&model.Message{}, id)
}

func PostRegistration(c *gin.Context) {
	arr_data := extractBasic(c)

	config.DB.Create(&model.User{
		Username: &arr_data[0],
		Password: &arr_data[1],
	})
}

func GetAccount(c *gin.Context) {
	arr_data := extractBasic(c)
	var user model.User
	var messages []model.Message

	config.DB.Find(&user, "username = ? AND password = ?", arr_data[0], arr_data[1])
	config.DB.Find(&messages, "user_id = ?", user.ID)

	c.JSON(http.StatusOK, messages)
}

func PatchAccount(c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		log.Println(err)
	}

	config.DB.Save(&user)
}

func DeleteAccount(c *gin.Context) {
	var user model.User
	var messages model.Message
	arr_data := extractBasic(c)

	config.DB.Find(&user, "username = ? AND password = ?", arr_data[0], arr_data[1])
	config.DB.Delete(&messages, "user_id = ?", user.ID)
	config.DB.Delete(&user)
}
