package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mfuss/internal/entity"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostJSONHandler(c *gin.Context) {

	var input entity.URLInput
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(c.BindJSON(&input)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	for {
		if err := json.NewDecoder(buf).Decode(&input); err == io.EOF {
			break
		} else if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}
	}

	shortURLId, err := h.Repo.URLStorage.SaveURL(input.URL)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	short := fmt.Sprintf("%v/%v", h.Repo.Config.BaseURL, shortURLId)

	if _, err := url.ParseRequestURI(short); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", short), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, entity.URLResponse{Result: short})

}
