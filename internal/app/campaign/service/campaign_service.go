package service

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
)

type Campaign interface {
	FindAll() ([]domain.Campaign, error)
	FindByUserID(userID int) ([]domain.Campaign, error)
	FindByID(ID int) (domain.Campaign, error)
	Save(campaign domain.Campaign) (domain.Campaign, error)
	Update(campaign domain.Campaign) (domain.Campaign, error)
	CreateImage(campaignImage domain.CampaignImage) (domain.CampaignImage, error)
	MarkAllImageIsPrimary(campaignID int) (bool, error)
}
type campaignService struct {
	campaignRepo Campaign
}

func NewCampaignService(camp Campaign) *campaignService {
	return &campaignService{
		campaignRepo: camp,
	}
}

func (s *campaignService) GetCampaigns(userID int) ([]domain.Campaign, error) {
	if userID != 0 {
		campaigns, err := s.campaignRepo.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.campaignRepo.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (s *campaignService) GetCampaignByID(param domain.InputIDParam) (domain.Campaign, error) {
	campaign, err := s.campaignRepo.FindByID(param.ID)
	if err != nil {
		return campaign, err
	}

	return campaign, nil

}

func (s *campaignService) CreateCampaign(param domain.CreateCampaignParam) (domain.Campaign, error) {
	campaign := domain.Campaign{}
	campaign.Name = param.Name
	campaign.ShortDescription = param.ShortDescription
	campaign.Description = param.Description
	campaign.GoalAmount = param.GoalAmount
	campaign.Perks = param.Perks
	campaign.UserID = param.User.ID
	newSlug := fmt.Sprintf("%s %d", campaign.Name, campaign.User.ID)
	campaign.Slug = slug.Make(newSlug)
	newCampaign, err := s.campaignRepo.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}

func (s *campaignService) UpdateCampaign(inputID domain.InputIDParam, param domain.CreateCampaignParam) (domain.Campaign, error) {
	campaign, err := s.campaignRepo.FindByID(inputID.ID)

	if err != nil {
		return campaign, err
	}

	if campaign.UserID != param.User.ID {
		return campaign, errors.New("Not an owner of the campaign")
	}
	campaign.Name = param.Name
	campaign.ShortDescription = param.ShortDescription
	campaign.Description = param.Description
	campaign.Perks = param.Perks
	campaign.GoalAmount = param.GoalAmount

	newCampaign, err := s.campaignRepo.Update(campaign)
	if err != nil {
		return newCampaign, err
	}

	return newCampaign, nil
}

func (s *campaignService) SaveCampaignImage(image domain.CreateCampaignImageParam, fileLocation string) (domain.CampaignImage, error) {
	campaign, err := s.campaignRepo.FindByID(image.CampaignID)
	if err != nil {
		return domain.CampaignImage{}, err
	}
	if campaign.UserID != image.User.ID {
		return domain.CampaignImage{}, errors.New("Not an Owner of the campaign")
	}
	isPrimary := 0
	if *image.IsPrimary {
		isPrimary = 1
		_, err := s.campaignRepo.MarkAllImageIsPrimary(image.CampaignID)
		if err != nil {
			return domain.CampaignImage{}, err
		}

	}

	campaignImage := domain.CampaignImage{}
	campaignImage.CampaignID = image.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.campaignRepo.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}
	return newCampaignImage, nil

}
