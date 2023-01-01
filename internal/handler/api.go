package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
		return
	}

	var input entity.URLInput
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(c.BindJSON(&input)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(buf).Decode(&input); err == io.EOF {

	} else if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.Repo.URLStorage.SaveURL(input.URL, userID)

	switch {
	case err != nil:
		if errors.Is(err, utilits.URLConflict{}) {
			if err = utilits.CheckURL(shortURL); err != nil {
				http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
			}

			c.JSON(http.StatusConflict, entity.URLResponse{Result: shortURL})
		}
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	default:
		c.JSON(http.StatusCreated, entity.URLResponse{Result: shortURL})
	}

}

func (h *Handler) MultipleShortHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}
	var input []entity.URLBatchInput
	var resOutput entity.URLBatchResponse
	var responseBatch []entity.URLBatchResponse

	err = json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

	for _, v := range input {
		res, err := h.Repo.URLStorage.SaveURL(v.URL, userID)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		}
		resOutput.CorrelID = v.CorrelID
		resOutput.URL = res
		responseBatch = append(responseBatch, resOutput)

	}

	c.JSON(http.StatusCreated, responseBatch)

}
