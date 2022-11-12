package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/user/handler"
)

type userHandler struct {
	userService handler.UserService
}

func NewUserHandler(usr handler.UserService) *userHandler {
	return &userHandler{userService: usr}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUser()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users})

}
