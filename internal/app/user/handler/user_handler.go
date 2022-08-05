package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
)

type UserService interface {
	Register(input domain.UserParam) (domain.User, error)
}

type userHandler struct {
	userService UserService
}

func NewUserhandler(service UserService) *userHandler {
	return &userHandler{
		userService: service,
	}
}

func (h *userHandler) Register(ctx *gin.Context) {
	var input domain.UserParam

	err := ctx.ShouldBind(&input)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := h.userService.Register(input)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "ok",
		"data":    user,
	})

}
