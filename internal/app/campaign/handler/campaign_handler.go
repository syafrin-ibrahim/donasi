package handler

import "github.com/syafrin-ibrahim/donasi.git/internal/app/domain"

type CampaignService interface {
	FindCampaigns(userID int) ([]domain.Campaign, error)
}

type campaignHandler struct {
	campaignService CampaignService
}

func NewCampaignHandler(service CampaignService) *campaignHandler {
	return &campaignHandler{
		campaignService: service,
	}
}
