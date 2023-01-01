package handler

import (
	"context"
	"errors"
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

	shortURL, err := h.Repo.URLStorage.SaveURL(string(body), userID)

	switch {
	case err != nil:
		if errors.Is(err, utilits.URLConflict{}) {
			if err = utilits.CheckURL(shortURL); err != nil {
				http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
			}
			c.Writer.WriteHeader(http.StatusConflict)
			c.Writer.Write([]byte(shortURL))
		}
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	default:
		c.Writer.WriteHeader(http.StatusCreated)
		c.Writer.Write([]byte(shortURL))
	}

}

func (h *Handler) GetURLHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		return
	}

	id := c.Param("id")

	sURL, err := h.Repo.URLStorage.GetShortURL(id, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}

	c.Writer.Header().Set("Location", sURL.Origin)

	c.Writer.WriteHeader(http.StatusTemporaryRedirect)

}

func (h *Handler) GetPingHandler(c *gin.Context) {
	if h.Repo.DB != nil {
		err := h.Repo.DB.Ping(context.Background())
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
		}
		c.Writer.WriteHeader(http.StatusOK)
	} else {
		c.Writer.Write([]byte("no "))
		c.Writer.WriteHeader(http.StatusInternalServerError)
	}

}
