package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
)

type UserService interface {
	Register(input domain.UserParam) (domain.User, error)
	Login(input domain.LoginParam) (domain.User, error)
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
		return
	}

	responseFormat := domain.FormatUserResponse(newUser, "asdfghjkkkkkkkkkkkkkkkk")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", responseFormat)
	ctx.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(ctx *gin.Context) {
	var input domain.LoginParam
	err := ctx.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.userService.Login(input)

	if err != nil {

		errorMessage := gin.H{
			"errors": err.Error(),
		}
		response := helper.APIResponse("Login failed", http.StatusBadRequest, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	responseFormat := domain.FormatUserResponse(user, "asdfghjkkkkkkkkkkkkkkkk")

	response := helper.APIResponse("SUccess Login", http.StatusOK, "success", responseFormat)
	ctx.JSON(http.StatusOK, response)

}
