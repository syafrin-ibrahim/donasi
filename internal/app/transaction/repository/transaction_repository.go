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
