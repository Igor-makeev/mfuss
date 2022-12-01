package handler

import (
	"fmt"
	"io"
	"mfuss/internal/entity"
	"mfuss/internal/repositories"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StorageMock struct {
	store map[int]entity.ShortURL
	ID    int
}

func (store *StorageMock) SaveURL(input string) string {
	url := entity.ShortURL{
		ID:     store.ID,
		Result: strconv.Itoa(store.ID),
		Origin: input}

	store.store[store.ID] = url

	return url.Result
}

func (store *StorageMock) GetShortURL(id int) (sURL entity.ShortURL, er error) {
	s, ok := store.store[id]
	if ok {
		return s, nil
	}
	return entity.ShortURL{}, fmt.Errorf("url with id=%d not found", id)

}

func TestHandler_PostHandler(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/", strings.NewReader("https://kanobu.ru/"))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h := NewHandler(repositories.NewRepository(&StorageMock{store: make(map[int]entity.ShortURL), ID: 0}))
	h.PostHandler(rr, req)
	defer rr.Result().Body.Close()
	assert.Equal(t, http.StatusCreated, rr.Result().StatusCode, "wrong status code")

	expected := "http://localhost:8080/0"
	result := rr.Result().Body
	bodyRes, err := io.ReadAll(result)
	require.NoError(t, err)
	assert.Equalf(t, expected, string(bodyRes), "handler returned unexpected body: got %v want %v", string(bodyRes), expected)
	logrus.Print("body checked")
	assert.Equalf(t, http.StatusCreated, rr.Result().StatusCode, "handler returned wrong status code: got %v want %v", rr.Result().StatusCode, http.StatusCreated)
	logrus.Print("status code checked")
}

func TestHandler_GetURLHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/0", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	h := NewHandler(repositories.NewRepository(&StorageMock{store: make(map[int]entity.ShortURL), ID: 0}))
	h.repository.URLStorage.SaveURL("https://kanobu.ru/")
	h.GetURLHandler(rr, req, 0)
	defer rr.Result().Body.Close()
	expectedHeader := "https://kanobu.ru/"
	result := rr.Result()
	resHeader := result.Header.Get("Location")
	assert.Equalf(t, expectedHeader, resHeader, "handler return wrong header: got %v want %v", resHeader, expectedHeader)
	logrus.Print("header checked")
	expectedStatusCode := http.StatusTemporaryRedirect
	resStatus := result.StatusCode
	assert.Equalf(t, expectedStatusCode, resStatus, "handler return wrong status code: got %v want %v", resStatus, expectedStatusCode)
	logrus.Print("status code checked")
}
