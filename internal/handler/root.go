package handler

import (
	"fmt"
	"io"
	errorsEntity "mfuss/internal/entity/errors"
	"mfuss/internal/utilits"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// PostHandler = хэндлер POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
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

	shortURL, err := h.Service.SaveURL(c.Request.Context(), string(body), userID)

	if err != nil {
		_, ok := err.(errorsEntity.URLConflict)

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

// GetURLHandler - хэндлер GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
func (h *Handler) GetURLHandler(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	id := c.Param("id")

	sURL, err := h.Service.GetShortURL(c.Request.Context(), id, userID)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadGateway)
		return
	}

	if sURL.IsDeleted {
		c.Status(http.StatusGone)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, sURL.Origin)
	}

}

// GetPingHandler хендлер GET /ping, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error
func (h *Handler) GetPingHandler(c *gin.Context) {
	if err := h.Service.Ping(c.Request.Context()); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)

	}
	c.Status(http.StatusOK)

}
