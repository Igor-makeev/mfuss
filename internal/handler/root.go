package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}
	body = []byte(strings.Trim(string(body), "\n"))

	if _, err := url.ParseRequestURI(string(body)); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	shortURLId, err := h.Repo.URLStorage.SaveURL(string(body), userID)

	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	short := fmt.Sprintf("%v/%v", h.Repo.Config.BaseURL, shortURLId)

	if _, err := url.ParseRequestURI(short); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", short), http.StatusInternalServerError)
		return
	}

	c.Writer.WriteHeader(http.StatusCreated)
	c.Writer.Write([]byte(short))

}

func (h *Handler) GetURLHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	id := c.Param("id")

	sURL, err := h.Repo.URLStorage.GetShortURL(id, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusNotFound)
		return
	}

	c.Writer.Header().Set("Location", sURL.Origin)
	c.Writer.WriteHeader(http.StatusTemporaryRedirect)

}
