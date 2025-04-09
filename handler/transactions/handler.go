package transactions

import (
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
	"strconv"
	"strings"
)

type transactionsHandler struct {
	transactionSvc services.Transactions
}

func New(transactionSvc services.Transactions) handler.Transactions {
	return &transactionsHandler{transactionSvc: transactionSvc}
}

func (h *transactionsHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var transaction *models.Transaction

	err := ctx.Bind(&transaction)
	if err != nil {
		return nil, err
	}

	newTransaction, err := h.transactionSvc.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (h *transactionsHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	var f filters.Transactions

	f.Type = ctx.Params("type")
	f.Category = ctx.Params("category")
	startDate := ctx.Params("startDate")
	endDate := ctx.Params("endDate")

	f.StartDate = startDate[0] + " 00:00:00"
	f.EndDate = endDate[0] + " 23:59:59"

	transactions, err := h.transactionSvc.GetAll(ctx, &f)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (h *transactionsHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	transaction, err := h.transactionSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (h *transactionsHandler) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	var transaction *models.Transaction

	err = ctx.Bind(&transaction)
	if err != nil {
		return nil, err
	}

	transaction.ID = id

	updatedTransaction, err := h.transactionSvc.Update(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (h *transactionsHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, err
	}

	err = h.transactionSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "transaction deleted successfully", nil
}
