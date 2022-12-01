package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostHandler(c *gin.Context) {

	b, err := io.ReadAll(c.Request.Body)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := url.ParseRequestURI(string(b)); err != nil {

		http.Error(c.Writer, fmt.Sprintf("invalid URL: %v", string(b)), http.StatusInternalServerError)

		return
	}

	short := "http://" + c.Request.Host + c.Request.URL.Path + h.repository.SaveURL(string(b))

	if _, err := url.ParseRequestURI(short); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", short), http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write([]byte(short))

}

func (h *Handler) GetURLHandler(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
	}

	sURL, err := h.repository.GetShortURL(id)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.Writer.Header().Set("Location", sURL.Origin)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)

}
