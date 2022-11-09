package repository

import (
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	_ "github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"gorm.io/gorm"
)

type transactionDBRepository struct {
	db *gorm.DB
}

func NewTransactionDBRepository(db *gorm.DB) *transactionDBRepository {
	return &transactionDBRepository{
		db: db,
	}
}

func (r *transactionDBRepository) GetByCampaignID(campaignID int) ([]domain.Transaction, error) {
	var transaction []domain.Transaction
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}
func (r *transactionDBRepository) GetUserByID(userID int) ([]domain.Transaction, error) {
	var transactions []domain.Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Order("id desc").Where("user_id = ?", userID).Find(&transactions).Error

	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
