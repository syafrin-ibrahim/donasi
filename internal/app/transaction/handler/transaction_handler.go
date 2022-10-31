package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syafrin-ibrahim/donasi.git/internal/app/domain"
	"github.com/syafrin-ibrahim/donasi.git/internal/package/helper"
)

type TransactionService interface {
	GetTransactionByCampaignID(campaignID domain.Transactionparam) ([]domain.Transaction, error)
}

type TransactionHandler struct {
	transactionService TransactionService
}

func NewTransactionHandler(trx TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: trx,
	}
}

func (h *TransactionHandler) GetCampaignTransaction(ctx *gin.Context) {
	var param domain.Transactionparam
	err := ctx.ShouldBindUri(&param)
	if err != nil {
		response := helper.APIResponse("Error to get campaign transaction", http.StatusUnprocessableEntity, "error", nil)
		ctx.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := ctx.MustGet("currentUser").(domain.User)
	param.User = currentUser

	transactions, err := h.transactionService.GetTransactionByCampaignID(param)

	if err != nil {
		response := helper.APIResponse("Error to get campaign transaction", http.StatusBadRequest, "error", nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Transaction Campaign", http.StatusOK, "success", domain.FormatTransactionList(transactions))
	ctx.JSON(http.StatusOK, response)
}
