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

	"github.com/gin-gonic/gin"
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
	c, _ := gin.CreateTestContext(rr)
	c.Request = req
	h := NewHandler(repositories.NewRepository(&StorageMock{store: make(map[int]entity.ShortURL), ID: 0}))
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
	h := NewHandler(repositories.NewRepository(&StorageMock{store: make(map[int]entity.ShortURL), ID: 0}))
	h.repository.URLStorage.SaveURL("https://kanobu.ru/")
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/0", nil)
	if err != nil {
		t.Fatal(err)
	}

	c.Request = req
	c.AddParam("id", "0")

	h.GetURLHandler(c)
	go http.ListenAndServe(":8080", h.InitRoutes())

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
