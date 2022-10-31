package domain

import "time"

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	User       User
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Transactionparam struct {
	ID   int `uri:"id" binding:"required"`
	User User
}

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func FormatTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{
		ID:        transaction.ID,
		Name:      transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
	return formatter
}

func FormatTransactionList(transactions []Transaction) []CampaignTransactionFormatter {
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}
	var listTransaction []CampaignTransactionFormatter

	for _, transaction := range transactions {
		formatter := FormatTransaction(transaction)
		listTransaction = append(listTransaction, formatter)

	}

	return listTransaction
}
