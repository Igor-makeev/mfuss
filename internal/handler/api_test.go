package handler

import (
	"io"
	"mfuss/configs"
	"mfuss/internal/entity"
	"mfuss/internal/mock"
	"mfuss/internal/repositories"
	"mfuss/internal/service"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_PostJSONHandler(t *testing.T) {
	cfg := configs.Config{SrvAddr: "localhost:8080", BaseURL: "http://localhost:8080"}

	storage := mock.NewStorageMock(&cfg)
	rep := &repositories.Repository{
		URLStorager: storage,
		Config:      &cfg,
	}
	service := service.NewService(rep)
	exampleReq := entity.URLInput{URL: "https://kanobu.ru/"}
	body, _ := json.Marshal(exampleReq)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/shorten", strings.NewReader(string(body)))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	c.Set("userID", "test")
	h := NewHandler(service)
	h.PostJSONHandler(c)

	result := rr.Result()
	assert.Equal(t, http.StatusCreated, result.StatusCode, "wrong status code")
	exampleResp := entity.URLResponse{Result: "http://localhost:8080/0"}
	expectedBody, _ := json.Marshal(exampleResp)
	defer result.Body.Close()
	bodyRes, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.Equalf(t, string(expectedBody), string(bodyRes), "handler returned unexpected body: got %v want %v", string(bodyRes), expectedBody)

	assert.Equalf(t, http.StatusCreated, result.StatusCode, "handler returned wrong status code: got %v want %v", result.StatusCode, http.StatusCreated)

}
