package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mfuss/internal/entity"
	"mfuss/internal/utilits"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostJSONHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

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

	shortURL, err := h.Repo.URLStorager.SaveURL(input.URL, userID)

	if err != nil {
		_, ok := err.(utilits.URLConflict)

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

func (h *Handler) MultipleShortHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}
	var input []entity.URLBatchInput

	err = json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	responseBatch, err := h.Repo.URLStorager.MultipleShort(input, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, responseBatch)

}

func (h *Handler) GetUSERURLS(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	urls := h.Repo.GetAllURLS(userID)
	for i, v := range urls {
		urls[i].ResultURL = h.Repo.Config.BaseURL + "/" + v.ID
	}
	if len(urls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, urls)

}

func (h *Handler) DeleteUrls(c *gin.Context) {

	id, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	inputArray, err := getUrlsArray(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	_ = h.Repo.URLStorager.MarkAsDeleted(inputArray, id)

	c.Status(http.StatusAccepted)
}
