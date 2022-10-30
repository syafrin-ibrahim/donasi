package repository

import (
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"gorm.io/gorm"
)

// type Campaign interface {
// }
type campaignDBRepository struct {
	db *gorm.DB
}

func NewCampaignDBRepository(db *gorm.DB) *campaignDBRepository {
	return &campaignDBRepository{
		db: db,
	}
}

func (r *campaignDBRepository) FindAll() ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *campaignDBRepository) FindByUserID(userID int) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *campaignDBRepository) FindByID(ID int) (domain.Campaign, error) {
	var campaign domain.Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *campaignDBRepository) Save(campaign domain.Campaign) (domain.Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

func (r *campaignDBRepository) Update(campaign domain.Campaign) (domain.Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}
func (r *campaignDBRepository) CreateImage(campaignImage domain.CampaignImage) (domain.CampaignImage, error) {
	err := r.db.Create(&campaignImage).Error
	if err != nil {
		return campaignImage, err
	}
	return campaignImage, nil
}
func (r *campaignDBRepository) MarkAllImageIsPrimary(campaignID int) (bool, error) {
	err := r.db.Model(&domain.CampaignImage{}).Where("campaign_id = ? ", campaignID).Update("is_primary", false).Error
	if err != nil {
		return false, err
	}
	return true, nil

}
