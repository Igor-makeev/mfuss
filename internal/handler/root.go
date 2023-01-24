package handler

import (
	"fmt"
	"io"
	"mfuss/internal/utilits"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	} else {
		if err := utilits.CheckURL(shortURL); err != nil {
			http.Error(c.Writer, fmt.Sprintf("output data: %v is invalid URL", shortURL), http.StatusInternalServerError)
		}
		c.Status(http.StatusCreated)
		c.Writer.Write([]byte(shortURL))
	}

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
	logrus.Printf("in handler: bool flag %v", sURL.IsDelited)
	if sURL.IsDelited {
		c.Status(http.StatusGone)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, sURL.Origin)
	}

}

func (h *Handler) GetPingHandler(c *gin.Context) {
	if err := h.Repo.URLStorager.Ping(); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)

	}
	c.Status(http.StatusOK)

}
