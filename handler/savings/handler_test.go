package savings

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
	"moneyManagement/models"
	"moneyManagement/services"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	savingsSvc := services.NewMockSavings(ctrl)

	saving := &models.Savings{ID: 1, UserID: 1, TransactionID: 1, Amount: 100, Type: "FD"}

	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", []byte(`{"userID":1,"transactionID":1,"amount":100,"type":"FD"}`), saving, nil,
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().Create(ctx, &models.Savings{UserID: 1, TransactionID: 1, Amount: 100, Type: "FD"}).Return(saving, nil)
			}},
		{"Failure Case: Error from service layer", []byte(`{"userID":1,"transactionID":1,"amount":100,"type":"FD"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().Create(ctx, &models.Savings{UserID: 1, TransactionID: 1, Amount: 100, Type: "FD"}).Return(nil, errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/savings", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(savingsSvc)

			output, err := h.Create(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	savingsSvc := services.NewMockSavings(ctrl)

	saving := &models.Savings{ID: 1, UserID: 1, TransactionID: 1, Amount: 100, Type: "FD"}

	tests := []struct {
		description    string
		id             string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []byte(`{"userID":1,"transactionID":1,"amount":100,"type":"FD"}`), saving, nil,
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().Update(ctx, saving).Return(saving, nil)
			}},
		{"Failure Case: Error from service layer", "1", []byte(`{"userID":1,"transactionID":1,"amount":100,"type":"FD"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().Update(ctx, saving).Return(nil, errors.New("error"))
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
			req := httptest.NewRequest(http.MethodPut, "/savings", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(savingsSvc)

			output, err := h.Update(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	savingsSvc := services.NewMockSavings(ctrl)

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", "saving deleted successfully", nil,
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().Delete(ctx, 1).Return(nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().Delete(ctx, 1).Return(errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/savings", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(savingsSvc)

			output, err := h.Delete(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	savingsSvc := services.NewMockSavings(ctrl)

	saving := &models.Savings{ID: 1, UserID: 1, TransactionID: 1, Amount: 100, Type: "FD"}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", saving, nil,
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().GetByID(ctx, 1).Return(saving, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().GetByID(ctx, 1).Return(nil, errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/savings", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(savingsSvc)

			output, err := h.GetByID(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	savingsSvc := services.NewMockSavings(ctrl)

	saving := &models.Savings{ID: 1, UserID: 1, TransactionID: 1, Amount: 100, Type: "FD"}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []*models.Savings{saving}, nil,
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().GetAll(ctx).Return([]*models.Savings{saving}, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				savingsSvc.EXPECT().GetAll(ctx).Return(nil, errors.New("error"))
			}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/saving", nil)
			req.Header.Set("Content-Type", "application/json")

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(savingsSvc)

			output, err := h.GetAll(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}
