package service

import "github.com/syafrin-ibrahim/donasi.git/internal/app/domain"

type Campaign interface {
	FIndAll() ([]domain.Campaign, error)
	FindByUserID(userID int) ([]domain.Campaign, error)
}
type campaignService struct {
	campaignRepo Campaign
}

func NewCampaignService(camp Campaign) *campaignService {
	return &campaignService{
		campaignRepo: camp,
	}
}

//func (s *campaignService)()
