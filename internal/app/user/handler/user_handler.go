package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
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
		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.Register(input)
	if err != nil {

		errors := helper.FormatError(err)
		errorMessage := gin.H{
			"errors": errors,
		}
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
	}

	responseFormat := domain.FormatUserResponse(newUser, "asdfghjkkkkkkkkkkkkkkkk")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", responseFormat)
	ctx.JSON(http.StatusOK, response)

}
