package service

import (
	"errors"

	"github.com/syafrin-ibrahim/donasi.git/internal/app/campaign/service"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
)

type Transaction interface {
	GetByCampaignID(campaignID int) ([]domain.Transaction, error)
	GetUserByID(userID int) ([]domain.Transaction, error)
}

type transactionService struct {
	TransactionRepo Transaction
	CampaignRepo    service.Campaign
}

func NewTransactionService(trx Transaction, campaignRepo service.Campaign) *transactionService {
	return &transactionService{
		TransactionRepo: trx,
		CampaignRepo:    campaignRepo,
	}
}

func (s *transactionService) GetTransactionByCampaignID(campaignID domain.Transactionparam) ([]domain.Transaction, error) {
	campaign, err := s.CampaignRepo.FindByID(campaignID.ID)
	if err != nil {
		return []domain.Transaction{}, err
	}

	if campaign.UserID != campaignID.User.ID {
		return []domain.Transaction{}, errors.New("Not An Owner Campaign")
	}
	transactions, err := s.TransactionRepo.GetByCampaignID(campaignID.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *transactionService) GetTransactonByuserID(userID int) ([]domain.Transaction, error) {
	transactions, err := s.TransactionRepo.GetUserByID(userID)

	if err != nil {
		return transactions, err
	}

	return transactions, nil

}
