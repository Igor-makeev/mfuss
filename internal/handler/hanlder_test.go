package handler

import (
	"io"
	"mfuss/configs"
	mock "mfuss/internal/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_PostHandler(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("https://kanobu.ru/"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	store := mock.NewStorageMock()
	cfg := configs.Config{SrvAddr: "localhost:8080", BaseURL: "http://localhost:8080"}
	h := NewHandler(store, cfg)
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
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	store := mock.NewStorageMock()
	cfg := configs.Config{SrvAddr: "localhost:8080", BaseURL: "http://localhost:8080"}
	h := NewHandler(store, cfg)
	h.storage.SaveURL("https://kanobu.ru/")
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = req
	c.AddParam("id", "0")

	h.GetURLHandler(c)
	go http.ListenAndServe("127.0.0.1:8080", h.Router)

	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	expectedHeader := "https://kanobu.ru/"
	result := rr.Result()
	defer result.Body.Close()
	resHeader := result.Header.Get("Location")
	assert.Equalf(t, expectedHeader, resHeader, "handler return wrong header: got %v want %v", resHeader, expectedHeader)

	expectedStatusCode := http.StatusTemporaryRedirect
	resStatus := res.StatusCode
	assert.Equalf(t, expectedStatusCode, resStatus, "handler return wrong status code: got %v want %v", resStatus, expectedStatusCode)

}
