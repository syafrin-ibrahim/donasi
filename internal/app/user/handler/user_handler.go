package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
)

type UserService interface {
	Register(input domain.UserParam) (domain.User, error)
	Login(input domain.LoginParam) (domain.User, error)
	IsEmailAvailable(input domain.CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (domain.User, error)
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

func (h *userHandler) CheckEmailAvailability(ctx *gin.Context) {
	var input domain.CheckEmailInput
	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		response := helper.APIResponse("Check Email Failedd", http.StatusUnprocessableEntity, "error", err.Error())
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{
			"errors": "server error",
		}
		response := helper.APIResponse("Check email Failed", http.StatusBadRequest, "error", errorMessage)
		ctx.JSON(http.StatusInternalServerError, response)
		return

	}

	data := gin.H{
		"is_available": isAvailable,
	}

	metaMessage := "Email has been registered"
	if isAvailable {
		metaMessage = "Email is available"
	}
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(ctx *gin.Context) {
	//input dari user
	//simpan gambarnya di folder images
	//di service kita panggil repository
	//dapatkan jwt (sementara hardcode, seakan user yang login = 1)
	//repo ambil data user yang id = 1
	// repo update data user simpan lokasi file

	file, err := ctx.FormFile("avatar")

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := 1
	path := fmt.Sprintf("internal/app/images/%d-%s", userID, file.Filename)

	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helper.APIResponse("Avatar Succesfuly uploaded", http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}
