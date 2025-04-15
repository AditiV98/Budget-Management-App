package transactions

import (
	"bytes"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gofr.dev/pkg/gofr"
	"gofr.dev/pkg/gofr/container"
	gofrHTTP "gofr.dev/pkg/gofr/http"
	"moneyManagement/filters"
	"moneyManagement/models"
	"moneyManagement/services"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionSvc := services.NewMockTransactions(ctrl)

	transaction := &models.Transaction{ID: 1, UserID: 1, Account: models.AccountDetails{ID: 1}, Amount: 100, Type: "EXPENSE"}

	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", []byte(`{"userId":1,"account":{"id":1},"amount":100,"type":"EXPENSE"}`), transaction, nil,
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().Create(ctx, &models.Transaction{UserID: 1, Account: models.AccountDetails{ID: 1}, Amount: 100, Type: "EXPENSE"}).Return(transaction, nil)
			}},
		{"Failure Case: Error from service layer", []byte(`{"userId":1,"account":{"id":1},"amount":100,"type":"EXPENSE"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().Create(ctx, &models.Transaction{UserID: 1, Account: models.AccountDetails{ID: 1}, Amount: 100, Type: "EXPENSE"}).Return(nil, errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/transaction", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(transactionSvc)

			output, err := h.Create(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionSvc := services.NewMockTransactions(ctrl)

	transaction := &models.Transaction{ID: 1, UserID: 1, Account: models.AccountDetails{ID: 1}, Amount: 100, Type: "EXPENSE"}

	tests := []struct {
		description    string
		id             string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []byte(`{"userId":1,"account":{"id":1},"amount":100,"type":"EXPENSE"}`), transaction, nil,
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().Update(ctx, transaction).Return(transaction, nil)
			}},
		{"Failure Case: Error from service layer", "1", []byte(`{"userId":1,"account":{"id":1},"amount":100,"type":"EXPENSE"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().Update(ctx, transaction).Return(nil, errors.New("error"))
			}},
		{"Failure Case: bind error", "1", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
		{"Failure Case: invalid id", "!", []byte(`{`), nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/transaction", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(transactionSvc)

			output, err := h.Update(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionSvc := services.NewMockTransactions(ctrl)

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", "transaction deleted successfully", nil,
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().Delete(ctx, 1).Return(nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().Delete(ctx, 1).Return(errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/transaction", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(transactionSvc)

			output, err := h.Delete(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionSvc := services.NewMockTransactions(ctrl)

	transaction := &models.Transaction{ID: 1, UserID: 1, Account: models.AccountDetails{ID: 1}, Amount: 100, Type: "EXPENSE"}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", transaction, nil,
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().GetByID(ctx, 1).Return(transaction, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().GetByID(ctx, 1).Return(nil, errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/trasaction", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(transactionSvc)

			output, err := h.GetByID(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	transactionSvc := services.NewMockTransactions(ctrl)

	transaction := &models.Transaction{ID: 1, UserID: 1, Account: models.AccountDetails{ID: 1}, Amount: 100, Type: "EXPENSE"}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []*models.Transaction{transaction}, nil,
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().GetAll(ctx, &filters.Transactions{StartDate: "01-01-2025 00:00:00", EndDate: "30-01-2025 23:59:59"}).Return([]*models.Transaction{transaction}, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				transactionSvc.EXPECT().GetAll(ctx, &filters.Transactions{StartDate: "01-01-2025 00:00:00", EndDate: "30-01-2025 23:59:59"}).Return(nil, errors.New("error"))
			}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/transaction", nil)
			req.Header.Set("Content-Type", "application/json")
			qParam := url.Values{"startDate": {"01-01-2025"}, "endDate": {"30-01-2025"}}

			req.URL.RawQuery = qParam.Encode()
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(transactionSvc)

			output, err := h.GetAll(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}
