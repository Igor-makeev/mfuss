package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUSERURLS(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		return
	}

	urls := h.Repo.GetAllURLS(userID)
	for i, v := range urls {
		urls[i].ResultURL = h.Repo.Config.BaseURL + "/" + v.ID
	}
	if len(urls) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, urls)

}
