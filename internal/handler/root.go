package handler

import (
	"context"
	"fmt"
	"io"
	"mfuss/internal/utilits"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
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

	shortURL, err := h.Repo.URLStorager.SaveURL(string(body), userID)

	if err != nil {
		_, ok := err.(utilits.URLConflict)

		if !ok {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := utilits.CheckURL(shortURL); err != nil {
			http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
		}

		c.Status(http.StatusConflict)
		c.Writer.Write([]byte(shortURL))
	}

	if err := utilits.CheckURL(shortURL); err != nil {
		http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
	}
	c.Status(http.StatusCreated)
	c.Writer.Write([]byte(shortURL))
}

func (h *Handler) GetURLHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	id := c.Param("id")

	sURL, err := h.Repo.URLStorager.GetShortURL(id, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, sURL.Origin)

}

func (h *Handler) GetPingHandler(c *gin.Context) {
	if h.Repo.DB != nil {
		err := h.Repo.DB.Ping(context.Background())
		if err != nil {
			c.Status(http.StatusInternalServerError)
		}
		c.Status(http.StatusOK)
	} else {
		c.Writer.Write([]byte("no "))
		c.Status(http.StatusInternalServerError)
	}

}
