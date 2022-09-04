package repository

import (
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"gorm.io/gorm"
)

type Campaign interface {
}
type campaignDBRepository struct {
	db *gorm.DB
}

func NewCampaignDBRepository(db *gorm.DB) *campaignDBRepository {
	return &campaignDBRepository{
		db: db,
	}
}

func (r campaignDBRepository) FindAll() ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

func (r *campaignDBRepository) FindByUserID(userID int) ([]domain.Campaign, error) {
	var campaigns []domain.Campaign
	err := r.db.Where("user_id = ?", userID).Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
