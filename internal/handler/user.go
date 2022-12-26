package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUSERURLS(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	urls := h.Repo.GetAllURLS(userId)
	for i, v := range urls {
		urls[i].ResultURL = h.Repo.Config.BaseURL + "/" + v.ID
	}
	if len(urls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, urls)

}
