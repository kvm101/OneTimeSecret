package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"one_time_secret/config"
	"one_time_secret/internal/model"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestExtractBasic(t *testing.T) {
	username := "testuser"
	password := "secret"

	sum := sha256.Sum256([]byte(password))
	expectedHash := fmt.Sprintf("%x", sum)

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	req := httptest.NewRequest(http.MethodGet, "/dummy", nil)
	req.Header.Set("Authorization", "Basic "+auth)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	result := extractBasic(c)

	assert.Equal(t, username, result[0])
	assert.Equal(t, expectedHash, result[1])
}

func TestGetMessage(t *testing.T) {
	// Перевірка підключення до бази даних
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	// Створення роутера
	router := gin.Default()

	// Маршрут для створення повідомлення
	router.POST("/message", func(c *gin.Context) {
		// Тіло запиту
		var input struct {
			Text string `json:"text"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Збереження повідомлення в базі даних
		message := model.Message{
			Text: &input.Text,
		}
		if err := config.DB.Create(&message).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Повернення ID нового повідомлення
		c.JSON(http.StatusOK, gin.H{"id": message.ID})
	})

	// Маршрут для отримання повідомлення
	router.GET("/message/:id", GetMessage)

	// Створення тестового повідомлення
	messageText := "Test message"
	w := httptest.NewRecorder()
	body := fmt.Sprintf(`{"text": "%s"}`, messageText)
	req := httptest.NewRequest(http.MethodPost, "/message", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Перевіряємо, чи створено повідомлення
	var response struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal("Failed to parse response body:", err)
	}

	// Використовуємо ID для запиту GET
	id := response.ID
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/message/"+id, nil)
	router.ServeHTTP(w, req)

	// Перевірка статусу відповіді
	assert.Equal(t, http.StatusOK, w.Code)

	// Перевірка, чи відповідає тексту повідомлення
	assert.Contains(t, w.Body.String(), messageText)
}

func TestPOSTMessage(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.POST("/message", PostMessage)

	testID := uuid.New()
	jsonMessage := fmt.Sprintf(`{"ID": "%s", "Text": "This is a test message"}`, testID.String())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/message", bytes.NewBufferString(jsonMessage))
	req.Header.Add("Authorization", "Basic dGVzdDp0ZXN0") // Base64("test:test")
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	config.DB.Delete(&model.Message{}, "id = ?", testID)
}

func TestDeteleMessage(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.DELETE("/message/:id", DeleteMessage)

	textMessage := "This is a test message"
	testID := uuid.New()

	message := model.Message{
		ID:   &testID,
		Text: &textMessage,
	}

	config.DB.Create(&message)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/message/"+testID.String(), nil)
	req.Header.Add("Authorization", "dGVzdDp0ZXN0")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPatchMessage(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.PATCH("/message/:id", PatchMessage)

	jsonMessage := `{
		"Text": "This is a test message"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/message/40b629b3-e0c3-4527-baa7-ef8716456758", bytes.NewBufferString(jsonMessage))
	req.Header.Add("Authorization", "dGVzdDp0ZXN0")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPostRegistration(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.POST("/account/registration", PostRegistration)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/account/registration", nil)

	username := "test_username"
	password := "test_password"

	data := fmt.Sprintf("%s:%s", username, password)

	b64_auth := base64.StdEncoding.EncodeToString([]byte(data))

	req.Header.Add("Authorization", b64_auth)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAccount(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.GET("/account", GetAccount)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/account", nil)

	username := "test_username"
	password := "test_password"

	data := fmt.Sprintf("%s:%s", username, password)

	b64_auth := base64.StdEncoding.EncodeToString([]byte(data))

	req.Header.Add("Authorization", b64_auth)
	router.ServeHTTP(w, req)

	chech_html := fmt.Sprintf("Account: %s", username)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), chech_html)
}

func TestPatchAccount(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.PATCH("/account", PatchAccount)

	username := "test_username"
	password := "test_password"
	data := fmt.Sprintf("%s:%s", username, password)

	jsonAccountChanges := `{
		"username": "test_username",
		"password": "test_password"
	}`

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPatch, "/account", bytes.NewBufferString(jsonAccountChanges))
	b64_auth := base64.StdEncoding.EncodeToString([]byte(data))
	req.Header.Add("Authorization", b64_auth)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteAccount(t *testing.T) {
	if err := config.ConnectDatabase(); err != nil {
		t.Error(err)
	}

	if config.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	router := gin.Default()
	router.DELETE("/account", DeleteAccount)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/account", nil)

	username := "test_username"
	password := "test_password"

	data := fmt.Sprintf("%s:%s", username, password)

	b64_auth := base64.StdEncoding.EncodeToString([]byte(data))

	req.Header.Add("Authorization", b64_auth)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
