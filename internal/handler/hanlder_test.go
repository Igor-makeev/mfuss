package handler

import (
	"io/ioutil"
	"mfuss/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMyHandler_PostHandler(t *testing.T) {
	reader := strings.NewReader("https://kanobu.ru/")
	req, err := http.NewRequest(http.MethodPost, "localhost:8080/", reader)

	rec := httptest.NewRecorder()
	h := &MyHandler{
		store: storage.NewStorage(),
	}

	h.PostHandler(rec, req)

	res := rec.Result()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode, "wrong status code")

	defer res.Body.Close()

	resB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("cpould not read response: %v", err)
	}
	assert.Contains(t, string(resB), "http://localhost:8080/0")

}

func TestMyHandler_GetURLHandler(t *testing.T) {
	reader := strings.NewReader("https://kanobu.ru/")

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/", reader)

	rec := httptest.NewRecorder()

	h := &MyHandler{
		store: storage.NewStorage(),
	}
	h.PostHandler(rec, req)
	h.GetURLHandler(rec, req, 0)
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	res := rec.Result()

	assert.NoError(t, err)
	assert.Equal(t, client.CheckRedirect, res.StatusCode, "wrong status code")

}
