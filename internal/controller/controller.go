package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"one_time_secret/config"
	"one_time_secret/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SystemPath(after_path string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}

	project_name := "/OneTimeSecret"
	index := strings.Index(cwd, project_name)
	last_index := len(project_name) + index

	return cwd[:last_index] + after_path
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

func GetMessage(c *gin.Context) {
	var message model.Message
	str := c.Param("id")
	id, err := uuid.Parse(string(str))
	if err != nil {
		log.Println(err)
		c.Status(http.StatusBadRequest)
		return
	}

	// Знайти повідомлення за ID
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
		if os.Getenv("IS_TESTING") == "true" {
			var user model.User
			config.DB.Find(&user, "id = ?", message.UserID)

			data := model.MessageInfo{
				ID:        message.ID,
				Text:      message.Text,
				Times:     message.Times,
				Timestamp: message.Timestamp,
				Username:  user.Username,
			}

			c.JSON(http.StatusOK, data)
			return
		}

		tmpl := template.Must(template.ParseFiles(SystemPath("/internal/view/messages.html")))

		var user model.User
		config.DB.Find(&user, "id = ?", message.UserID)

		data := model.MessageInfo{
			ID:        message.ID,
			Text:      message.Text,
			Times:     message.Times,
			Timestamp: message.Timestamp,
			Username:  user.Username,
		}

		if err := tmpl.Execute(c.Writer, data); err != nil {
			log.Println(err)
			return
		}

	} else {
		tmpl := template.Must(template.ParseFiles(SystemPath("/internal/view/not_found.html")))

		c.Status(http.StatusNotFound)
		if err := tmpl.Execute(c.Writer, nil); err != nil {
			log.Println(err)
			return
		}
	}
}

func PostMessage(c *gin.Context) {
	message := bindModelMessage(c)
	if message.Text != nil {
		config.DB.Create(&message)
		return
	}

	c.Status(http.StatusNoContent)
}

func PatchMessage(c *gin.Context) {
	str := c.Param("id")
	id, err := uuid.Parse(str)
	if err != nil {
		log.Println(err)
		return
	}

	var message model.Message

	config.DB.Find(&message, "id = ?", id)
	message_change := bindModelMessage(c)
	message.Times = message_change.Times
	message.Text = message_change.Text
	message.ExpirationDate = message_change.ExpirationDate
	message.MessagePassword = message_change.MessagePassword

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

	log.Println(user)

	isAuth := false
	if user.Username != nil {
		isAuth = true
	}

	data := model.AccountData{
		Username: user.Username,
		Messages: &messages,
		IsAuth:   isAuth,
	}

	if os.Getenv("IS_TESTING") == "true" {
		c.JSON(200, data)
		return
	}

	tmpl := template.Must(template.ParseFiles(SystemPath("/internal/view/account.html")))
	if err := tmpl.Execute(c.Writer, data); err != nil {
		log.Println(err)
		return
	}
}

func PatchAccount(c *gin.Context) {
	arr_data := extractBasic(c)
	var user_changes model.User
	var user model.User

	if err := c.BindJSON(&user_changes); err != nil {
		log.Println(err)
	}

	config.DB.Find(&user, "username = ? and password = ?", arr_data[0], arr_data[1])

	if user_changes.Username != nil {
		user.Username = user_changes.Username
	}

	if user_changes.Password != nil {
		str_password := *user_changes.Password
		sum := fmt.Sprintf("%x", sha256.Sum256([]byte(str_password)))
		user.Password = &sum
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
