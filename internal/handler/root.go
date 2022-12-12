package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostHandler(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(string(body)); err != nil {

		http.Error(c.Writer, fmt.Sprintf("invalid URL: %v", string(body)), http.StatusInternalServerError)

		return
	}

	shortURLId, err := h.storage.SaveURL(string(body))

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	short := fmt.Sprintf("%v/%v", h.Cfg.BaseURL, shortURLId)

	if _, err := url.ParseRequestURI(short); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", short), http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write([]byte(short))

}

func (h *Handler) GetURLHandler(c *gin.Context) {

	id := c.Param("id")

	sURL, err := h.storage.GetShortURL(id)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.Writer.Header().Set("Location", sURL.Origin)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)

}
