//go:build integration

package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"one_time_secret/internal/controller"
	"one_time_secret/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFullFlow(t *testing.T) {
	if err := model.ConnectDatabase(); err != nil {
		t.Fatalf("Database connection failed: %v", err)
	}

	router := gin.Default()
	router.POST("/message", controller.PostMessage)
	router.GET("/message/:id", controller.GetMessage)

	messageID := uuid.New().String()
	jsonMessage := fmt.Sprintf(`{"ID": "%s", "Text": "This is a test message"}`, messageID)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/message", bytes.NewBufferString(jsonMessage))
	req.Header.Add("Authorization", "dGVzdDp0ZXN0") // Basic test:test
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var savedMessage model.Message
	err := model.DB.Where("id = ?", messageID).First(&savedMessage).Error
	if err != nil {
		t.Fatalf("Failed to fetch saved message: %v", err)
	}

	assert.Equal(t, "This is a test message", *savedMessage.Text)

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/message/"+messageID, nil)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	assert.Contains(t, w2.Body.String(), "This is a test message")
}
