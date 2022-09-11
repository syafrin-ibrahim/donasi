package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
)

type CampaignService interface {
	GetCampaigns(userID int) ([]domain.Campaign, error)
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
