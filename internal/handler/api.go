package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mfuss/internal/entity"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostJSONHandler(c *gin.Context) {

	var input entity.URLInput
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(c.Request.Body); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
	if err := json.NewDecoder(buf).Decode(&input); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	shortURLId, err := h.storage.SaveURL(input.URL)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	short := fmt.Sprintf("http://%v/%v", c.Request.Host, shortURLId)

	if _, err := url.ParseRequestURI(short); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", short), http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, entity.URLResponse{Result: short})

}
