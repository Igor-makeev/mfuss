package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mfuss/internal/entity"
	errorsEntity "mfuss/internal/entity/errors"
	"mfuss/internal/utilits"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PostJSONHandler — хэндлер принимающий в теле запроса JSON-объект {"url":"<some_url>"} и возвращающий в ответ объект {"result":"<shorten_url>"}.
func (h *Handler) PostJSONHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}
	//инициализируем структуру в которую будем декодить входящие данные
	var input entity.URLInput
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(c.BindJSON(&input)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(buf).Decode(&input); err != nil && err != io.EOF {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.Service.SaveURL(c.Request.Context(), input.URL, userID)

	if err != nil {
		_, ok := err.(errorsEntity.URLConflict)

		if !ok {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := utilits.CheckURL(shortURL); err != nil {
			http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
		}

		c.JSON(http.StatusConflict, entity.URLResponse{Result: shortURL})
		return
	}

	if err := utilits.CheckURL(shortURL); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
	}
	c.JSON(http.StatusCreated, entity.URLResponse{Result: shortURL})

}

// PostJSONHandler — хэндлер  POST /api/shorten/batch, принимающий в теле запроса множество URL для сокращения.
func (h *Handler) MultipleShortHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}
	//инициализируем структуру в которую будем декодить входящие данные
	var input []entity.URLBatchInput

	err = json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseBatch, err := h.Service.MultipleShort(c.Request.Context(), input, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responseBatch)

}

// GetUserURLs - хэндлер GET / возвращает все URL скращенные пользователем.
func (h *Handler) GetUserURLs(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	urls := h.Service.GetAllURLs(c.Request.Context(), userID)
	for i, v := range urls {
		urls[i].ResultURL = h.Service.Cfg.BaseURL + "/" + v.ID
	}
	if len(urls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, urls)

}

// DeleteUrls — асинхронный хендлер DELETE /api/user/urls, который принимает список идентификаторов сокращённых URL для удаления в формате:[ "a", "b", "c", "d", ...]
func (h *Handler) DeleteUrls(c *gin.Context) {

	inputArray, err := getUrlsArray(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	h.Service.Queue.Write(inputArray)
	c.Status(http.StatusAccepted)
}
