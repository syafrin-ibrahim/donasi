package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
)

type CampaignService interface {
	GetCampaigns(userID int) ([]domain.Campaign, error)
	GetCampaignByID(param domain.InputIDParam) (domain.Campaign, error)
	CreateCampaign(param domain.CreateCampaignParam) (domain.Campaign, error)
	UpdateCampaign(inputID domain.InputIDParam, param domain.CreateCampaignParam) (domain.Campaign, error)
	SaveCampaignImage(image domain.CreateCampaignImageParam, fileLocation string) (domain.CampaignImage, error)
}

type campaignHandler struct {
	campaignService CampaignService
}

func NewCampaignHandler(service CampaignService) *campaignHandler {
	return &campaignHandler{
		campaignService: service,
	}
}

func (h *campaignHandler) GetCampaigns(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.Query("user_id"))

	campaigns, err := h.campaignService.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Error to get Campaigns", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of Campaigns", http.StatusOK, "success", domain.FormatCampaigns(campaigns))
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(ctx *gin.Context) {
	var param domain.InputIDParam
	err := ctx.ShouldBindUri(&param)
	if err != nil {
		response := helper.APIResponse("Error to get detail of Campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	campaign, err := h.campaignService.GetCampaignByID(param)
	if err != nil {
		response := helper.APIResponse("Error to get param input", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign Detail", http.StatusOK, "success", domain.FormatDetailCampaign(campaign))
	ctx.JSON(http.StatusOK, response)
}
func (h *campaignHandler) CreateCampaign(ctx *gin.Context) {
	var input domain.CreateCampaignParam

	err := ctx.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to Create Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(domain.User)
	input.User = currentUser

	newCampaign, err := h.campaignService.CreateCampaign(input)

	if err != nil {
		response := helper.APIResponse("Failed to Create Campaign", http.StatusUnprocessableEntity, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to Create Campaing", http.StatusOK, "success", domain.FormatCampaign(newCampaign))
	ctx.JSON(http.StatusOK, response)
}

func (h *campaignHandler) UpdateCampaign(ctx *gin.Context) {
	var inputID domain.InputIDParam
	err := ctx.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to get Id Campaign", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var inputCampaign domain.CreateCampaignParam

	err = ctx.ShouldBindJSON(&inputCampaign)

	if err != nil {
		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(domain.User)
	inputCampaign.User = currentUser
	updatedCampaign, err := h.campaignService.UpdateCampaign(inputID, inputCampaign)

	if err != nil {
		response := helper.APIResponse("Failed update Campaign", http.StatusUnprocessableEntity, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to Update Campaing", http.StatusOK, "success", domain.FormatCampaign(updatedCampaign))
	ctx.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UploadImage(ctx *gin.Context) {
	var input domain.CreateCampaignImageParam

	err := ctx.ShouldBind(&input)
	if err != nil {

		errors := helper.FormatError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Failed to upload Campaign Image", http.StatusUnprocessableEntity, "error", errorMessage)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(domain.User)
	input.User = currentUser
	userID := currentUser.ID
	file, err := ctx.FormFile("file")
	if err != nil {

		errorMessage := gin.H{"is_uploaded": false}

		response := helper.APIResponse("Failed to upload Campaign Image", http.StatusBadRequest, "error", errorMessage)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	path := fmt.Sprintf("internal/app/images/%d-%s", userID, file.Filename)

	err = ctx.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.campaignService.SaveCampaignImage(input, path)

	if err != nil {
		data := gin.H{
			"is_uploaded": false,
		}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_uploaded": true,
	}
	response := helper.APIResponse("Campaign Image Succesfuly uploaded", http.StatusOK, "success", data)
	ctx.JSON(http.StatusOK, response)
}
