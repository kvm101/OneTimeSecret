package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"one_time_secret/internal/model"
	"one_time_secret/internal/view"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ExtractBasic(c *gin.Context) []string {
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

func RenderHTML(c *gin.Context, path string, data any) {
	tmpl := template.Must(template.ParseFS(view.TmplFS, path))
	if err := tmpl.Execute(c.Writer, data); err != nil {
		log.Println(err)
		return
	}
}

func bindModelMessage(c *gin.Context) model.Message {
	arr_data := ExtractBasic(c)

	var message model.Message
	var user model.User

	err := c.BindJSON(&message)
	log.Println(message)
	if err != nil {
		log.Println(err)
	}

	model.DB.Find(&user, "username = ? AND password = ?", arr_data[0], arr_data[1])

	message.UserID = user.ID

	return message
}

func GetMessage(c *gin.Context) {
	var message model.Message
	str := c.Param("id")

	id, err := uuid.Parse(str)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusNotFound)
		RenderHTML(c, "templates/page_not_found.html", nil)
		return
	}

	model.DB.First(&message, id)

	if message.Times != nil {
		if *message.Times > 1 {
			*message.Times = *message.Times - 1
			model.DB.Save(&message)
		} else {
			*message.Times = *message.Times - 1
			model.DB.Save(&message)
			model.DB.Delete(&message)
		}
	}

	if message.ID != nil {
		var user model.User
		model.DB.Find(&user, "id = ?", message.UserID)

		data := model.MessageInfo{
			ID:        message.ID,
			Text:      message.Text,
			Times:     message.Times,
			Timestamp: message.Timestamp,
			Username:  user.Username,
		}

		RenderHTML(c, "templates/messages.html", data)

	} else {
		RenderHTML(c, "templates/messages.html", nil)
		c.Status(http.StatusNotFound)
		return
	}
}

func PostMessage(c *gin.Context) {
	message := bindModelMessage(c)
	if message.Text != nil {
		if message.UserID != nil {
			model.DB.Create(&message)
		}
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

	model.DB.Find(&message, "id = ?", id)
	message_change := bindModelMessage(c)
	message.Times = message_change.Times
	message.Text = message_change.Text
	message.ExpirationDate = message_change.ExpirationDate
	message.MessagePassword = message_change.MessagePassword

	model.DB.Save(&message)
}

func DeleteMessage(c *gin.Context) {
	str := c.Param("id")
	id, err := uuid.Parse(str)
	if err != nil {
		log.Println(err)
		return
	}

	model.DB.Delete(&model.Message{}, id)
}

func PostRegistration(c *gin.Context) {
	arr_data := ExtractBasic(c)

	model.DB.Create(&model.User{
		Username: &arr_data[0],
		Password: &arr_data[1],
	})
}

func GetAccount(c *gin.Context) {
	arr_data := ExtractBasic(c)

	var user model.User
	var messages []model.Message

	model.DB.Find(&user, "username = ? AND password = ?", arr_data[0], arr_data[1])
	model.DB.Find(&messages, "user_id = ?", user.ID)

	log.Println(user)

	isAuth := false
	if user.Username != nil {
		isAuth = true
	}

	data := model.AccountData{
		Username: user.Username,
		Messages: &messages,
		IsAuth:   &isAuth,
	}

	RenderHTML(c, "templates/account.html", data)
}

func PatchAccount(c *gin.Context) {
	arr_data := ExtractBasic(c)
	var user_changes model.User
	var user model.User

	if err := c.BindJSON(&user_changes); err != nil {
		log.Println(err)
	}

	model.DB.Find(&user, "username = ? and password = ?", arr_data[0], arr_data[1])

	if user_changes.Username != nil {
		user.Username = user_changes.Username
	}

	if user_changes.Password != nil {
		str_password := user_changes.Password
		sum := fmt.Sprintf("%x", sha256.Sum256([]byte(*str_password)))
		user.Password = &sum
	}

	model.DB.Save(&user)
}

func DeleteAccount(c *gin.Context) {
	var user model.User
	var messages model.Message
	arr_data := ExtractBasic(c)

	model.DB.Find(&user, "username = ? AND password = ?", arr_data[0], arr_data[1])
	model.DB.Delete(&messages, "user_id = ?", user.ID)
	model.DB.Delete(&user)
}
