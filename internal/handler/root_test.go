package handler

import (
	"context"
	"io"
	"mfuss/configs"
	"mfuss/internal/mock"
	"mfuss/internal/repositories"
	"mfuss/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_PostHandler(t *testing.T) {
	cfg := configs.Config{SrvAddr: "localhost:8080", BaseURL: "http://localhost:8080"}
	storage := mock.NewStorageMock(&cfg)
	rep := &repositories.Repository{
		URLStorager: storage,
		Config:      &cfg,
	}

	service := service.NewService(rep)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("https://kanobu.ru/"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Set("userID", "test")
	h := NewHandler(service)

	h.PostHandler(c)
	result := rr.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode, "wrong status code")

	expected := "http://localhost:8080/0"
	defer result.Body.Close()
	bodyRes, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.Equalf(t, expected, string(bodyRes), "handler returned unexpected body: got %v want %v", string(bodyRes), expected)

	assert.Equalf(t, http.StatusCreated, result.StatusCode, "handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusCreated)

}

func TestHandler_GetURLHandler(t *testing.T) {

	cfg := configs.Config{SrvAddr: "localhost:8080", BaseURL: "http://localhost:8080"}
	storage := mock.NewStorageMock(&cfg)
	rep := &repositories.Repository{
		URLStorager: storage,
		Config:      &cfg,
	}
	service := service.NewService(rep)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Set("userID", "test")
	h := NewHandler(service)

	h.Service.SaveURL(context.Background(), "https://kanobu.ru/", "")
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = req
	c.AddParam("id", "0")

	h.GetURLHandler(c)

	expectedHeader := "https://kanobu.ru/"
	result := rr.Result()
	defer result.Body.Close()
	resHeader := result.Header.Get("Location")
	assert.Equalf(t, expectedHeader, resHeader, "handler return wrong header: got %v want %v", resHeader, expectedHeader)

	expectedStatusCode := http.StatusTemporaryRedirect
	resStatus := result.StatusCode
	assert.Equalf(t, expectedStatusCode, resStatus, "handler return wrong status code: got %v want %v", resStatus, expectedStatusCode)

}
