//go:build unit

package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"one_time_secret/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetMessage(t *testing.T) {
	// Ініціалізація з'єднання з БД
	if err := model.ConnectDatabase(); err != nil {
		t.Error(err)
	}
	if model.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	// Створення нового маршруту
	router := gin.Default()
	router.GET("/message/:id", GetMessage)

	// Створення повідомлення в БД
	message_id := uuid.New()
	message_text := "TEST"
	model.DB.Create(&model.Message{
		ID:   &message_id,
		Text: &message_text,
	})

	// Створення користувача для аутентифікації
	username := "unit_test"
	password := fmt.Sprintf("%x", sha256.Sum256([]byte("unit_test")))
	model.DB.Create(&model.User{
		Username: &username,
		Password: &password,
	})

	// Створення HTTP запиту
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/message/%s", message_id.String()), nil)
	req.Header.Add("Authorization", "Basic dWludF90ZXN0OnVuaXRfdGVzdA") // Base64("unit_test:unit_test")
	req.Header.Set("Content-Type", "application/json")

	// Виконання запиту через сервер
	router.ServeHTTP(w, req)

	// Перевірка результату
	assert.Equal(t, http.StatusOK, w.Code)

	// Очистка даних після тесту
	model.DB.Delete(&model.Message{}, "id = ?", message_id)
	model.DB.Delete(&model.User{}, "username = ? and password = ?", username, password)
}

func TestPOSTMessage(t *testing.T) {
	// Ініціалізація з'єднання з БД
	if err := model.ConnectDatabase(); err != nil {
		t.Error(err)
	}
	if model.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	// Створення нового маршруту
	router := gin.Default()
	router.POST("/message", PostMessage)

	// Підготовка тестових даних
	testID := uuid.New()
	jsonMessage := fmt.Sprintf(`{"ID": "%s", "Text": "This is a test message"}`, testID.String())

	// Створення HTTP запиту
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/message", bytes.NewBufferString(jsonMessage))
	req.Header.Add("Authorization", "Basic dGVzdDp0ZXN0") // Base64("test:test")
	req.Header.Set("Content-Type", "application/json")

	// Виконання запиту через сервер
	router.ServeHTTP(w, req)

	// Перевірка результату
	assert.Equal(t, http.StatusOK, w.Code)

	// Очистка даних після тесту
	model.DB.Delete(&model.Message{}, "id = ?", testID)
}

func TestDeleteMessage(t *testing.T) {
	// Ініціалізація з'єднання з БД
	if err := model.ConnectDatabase(); err != nil {
		t.Error(err)
	}
	if model.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	// Створення нового маршруту
	router := gin.Default()
	router.DELETE("/message/:id", DeleteMessage)

	// Створення тестового повідомлення
	textMessage := "This is a test message"
	testID := uuid.New()
	message := model.Message{
		ID:   &testID,
		Text: &textMessage,
	}
	model.DB.Create(&message)

	// Створення HTTP запиту для видалення
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/message/"+testID.String(), nil)
	req.Header.Add("Authorization", "dGVzdDp0ZXN0") // Base64("test:test")
	router.ServeHTTP(w, req)

	// Перевірка результату
	assert.Equal(t, http.StatusOK, w.Code)

	// Очистка даних після тесту
	model.DB.Delete(&model.Message{}, "id = ?", testID)
}

func TestPostRegistration(t *testing.T) {
	// Ініціалізація з'єднання з БД
	if err := model.ConnectDatabase(); err != nil {
		t.Error(err)
	}
	if model.DB == nil {
		t.Fatal("Database connection not initialized!")
	}

	// Створення нового маршруту
	router := gin.Default()
	router.POST("/account/registration", PostRegistration)

	// Підготовка тестових даних для реєстрації
	username := "test_username"
	password := "test_password"
	data := fmt.Sprintf("%s:%s", username, password)
	b64_auth := base64.StdEncoding.EncodeToString([]byte(data))

	// Створення HTTP запиту для реєстрації
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/account/registration", nil)
	req.Header.Add("Authorization", b64_auth)

	// Виконання запиту через сервер
	router.ServeHTTP(w, req)

	// Перевірка результату
	assert.Equal(t, http.StatusOK, w.Code)
}
