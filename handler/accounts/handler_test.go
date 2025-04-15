package accounts

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
	"testing"
)

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountSvc := services.NewMockAccount(ctrl)

	account := &models.Account{ID: 1, Name: "Cash", Type: "CASH", Balance: 2000}

	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", []byte(`{"name":"Cash","type":"CASH","balance":2000}`), account, nil,
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().Create(ctx, &models.Account{Name: "Cash", Type: "CASH", Balance: 2000}).Return(account, nil)
			}},
		{"Failure Case: Error from service layer", []byte(`{"name":"Cash","type":"CASH","balance":2000}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().Create(ctx, &models.Account{Name: "Cash", Type: "CASH", Balance: 2000}).Return(nil, errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/account", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(accountSvc)

			output, err := h.Create(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountSvc := services.NewMockAccount(ctrl)

	account := &models.Account{ID: 1, Name: "Cash", Type: "CASH", Balance: 2000}

	tests := []struct {
		description    string
		id             string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []byte(`{"name":"Cash","type":"CASH","balance":2000}`), account, nil,
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().Update(ctx, account).Return(account, nil)
			}},
		{"Failure Case: Error from service layer", "1", []byte(`{"name":"Cash","type":"CASH","balance":2000}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().Update(ctx, account).Return(nil, errors.New("error"))
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
			req := httptest.NewRequest(http.MethodPut, "/account", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(accountSvc)

			output, err := h.Update(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountSvc := services.NewMockAccount(ctrl)

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", "account deleted successfully", nil,
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().Delete(ctx, 1).Return(nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().Delete(ctx, 1).Return(errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/account", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(accountSvc)

			output, err := h.Delete(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountSvc := services.NewMockAccount(ctrl)

	account := &models.Account{ID: 1, Name: "Cash", Type: "CASH", Balance: 2000}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", account, nil,
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().GetByID(ctx, 1).Return(account, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().GetByID(ctx, 1).Return(nil, errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/account", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(accountSvc)

			output, err := h.GetByID(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	accountSvc := services.NewMockAccount(ctrl)

	account := &models.Account{ID: 1, Name: "Cash", Type: "CASH", Balance: 2000}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []*models.Account{account}, nil,
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().GetAll(ctx, &filters.Account{}).Return([]*models.Account{account}, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				accountSvc.EXPECT().GetAll(ctx, &filters.Account{}).Return(nil, errors.New("error"))
			}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/account", nil)
			req.Header.Set("Content-Type", "application/json")

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(accountSvc)

			output, err := h.GetAll(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}
