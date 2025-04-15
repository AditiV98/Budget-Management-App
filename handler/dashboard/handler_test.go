package dashboard

import (
	"context"
	"errors"
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

func Test_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	dashboardSvc := services.NewMockDashboard(ctrl)

	dashboard := models.Dashboard{TotalIncome: 200, TotalExpense: 100, TotalSavings: 100}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", dashboard, nil,
			func(ctx *gofr.Context) {
				dashboardSvc.EXPECT().Get(ctx, &filters.Transactions{AccountID: []int{1}, StartDate: "01-01-2025 00:00:00", EndDate: "30-01-2025 23:59:59"}).Return(dashboard, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				dashboardSvc.EXPECT().Get(ctx, &filters.Transactions{AccountID: []int{1}, StartDate: "01-01-2025 00:00:00", EndDate: "30-01-2025 23:59:59"}).Return(models.Dashboard{}, errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
			req.Header.Set("Content-Type", "application/json")
			qParam := url.Values{"accountId": {tc.id}, "startDate": {"01-01-2025"}, "endDate": {"30-01-2025"}}

			req.URL.RawQuery = qParam.Encode()

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(dashboardSvc)

			output, err := h.Get(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}
