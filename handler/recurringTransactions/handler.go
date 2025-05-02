package recurringTransactions

import (
	"errors"
	"gofr.dev/pkg/gofr"
	"moneyManagement/filters"
	"moneyManagement/handler"
	"moneyManagement/models"
	"moneyManagement/services"
	"strconv"
	"strings"
)

type recurringTransactionsHandler struct {
	recurringTransactionSvc services.RecurringTransactions
}

func New(recurringTransactionSvc services.RecurringTransactions) handler.RecurringTransactions {
	return &recurringTransactionsHandler{recurringTransactionSvc: recurringTransactionSvc}
}

func (h *recurringTransactionsHandler) Create(ctx *gofr.Context) (interface{}, error) {
	var recurringTransaction *models.RecurringTransaction

	err := ctx.Bind(&recurringTransaction)
	if err != nil {
		return nil, errors.New("bind error")
	}

	newTransaction, err := h.recurringTransactionSvc.Create(ctx, recurringTransaction)
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (h *recurringTransactionsHandler) GetAll(ctx *gofr.Context) (interface{}, error) {
	var f filters.RecurringTransactions

	f.Type = ctx.Params("type")
	f.Category = ctx.Params("category")
	startDate := ctx.Params("startDate")
	endDate := ctx.Params("endDate")

	if len(startDate) != 0 {
		f.StartDate = startDate[0] + " 00:00:00"
	}

	if len(endDate) != 0 {
		f.EndDate = endDate[0] + " 23:59:59"
	}

	transactions, err := h.recurringTransactionSvc.GetAll(ctx, &f)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (h *recurringTransactionsHandler) GetByID(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	transaction, err := h.recurringTransactionSvc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (h *recurringTransactionsHandler) Update(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	var recurringTransaction *models.RecurringTransaction

	err = ctx.Bind(&recurringTransaction)
	if err != nil {
		return nil, errors.New("bind error")
	}

	recurringTransaction.ID = id

	updatedTransaction, err := h.recurringTransactionSvc.Update(ctx, recurringTransaction)
	if err != nil {
		return nil, err
	}

	return updatedTransaction, nil
}

func (h *recurringTransactionsHandler) Delete(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	err = h.recurringTransactionSvc.Delete(ctx, id)
	if err != nil {
		return nil, err
	}

	return "transaction deleted successfully", nil
}

func (h *recurringTransactionsHandler) SkipNextRun(ctx *gofr.Context) (interface{}, error) {
	idString := strings.TrimSpace(ctx.PathParam("id"))

	id, err := strconv.Atoi(idString)
	if err != nil {
		return nil, errors.New("invalid id")
	}

	err = h.recurringTransactionSvc.SkipNextRun(ctx, id)
	if err != nil {
		return nil, err
	}

	return "Transaction skipped successfully", nil
}
