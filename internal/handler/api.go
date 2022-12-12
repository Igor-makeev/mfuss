package handler

import (
	"fmt"
	"mfuss/internal/entity"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostJSONHandler(c *gin.Context) {

	var input entity.URLInput
	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
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
