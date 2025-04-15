package users

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
	userSvc := services.NewMockUser(ctrl)

	user := &models.User{ID: 1, FirstName: "Aditi", LastName: "Verma", Email: "vermaditi2020@gmail.com"}

	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", []byte(`{"firstName": "Aditi", "lastName": "Verma", "email": "vermaditi2020@gmail.com"}`), user, nil,
			func(ctx *gofr.Context) {
				userSvc.EXPECT().Create(ctx, &models.User{FirstName: "Aditi", LastName: "Verma", Email: "vermaditi2020@gmail.com"}).Return(user, nil)
			}},
		{"Failure Case: Error from service layer", []byte(`{"firstName": "Aditi", "lastName": "Verma", "email": "vermaditi2020@gmail.com"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				userSvc.EXPECT().Create(ctx, &models.User{FirstName: "Aditi", LastName: "Verma", Email: "vermaditi2020@gmail.com"}).Return(nil, errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(userSvc)

			output, err := h.Create(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := services.NewMockUser(ctrl)

	user := &models.User{ID: 1, FirstName: "Aditi", LastName: "Verma", Email: "vermaditi2020@gmail.com"}

	tests := []struct {
		description    string
		id             string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []byte(`{"firstName": "Aditi", "lastName": "Verma", "email": "vermaditi2020@gmail.com"}`), user, nil,
			func(ctx *gofr.Context) {
				userSvc.EXPECT().Update(ctx, user).Return(user, nil)
			}},
		{"Failure Case: Error from service layer", "1", []byte(`{"firstName": "Aditi", "lastName": "Verma", "email": "vermaditi2020@gmail.com"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				userSvc.EXPECT().Update(ctx, user).Return(nil, errors.New("error"))
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
			req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(userSvc)

			output, err := h.Update(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := services.NewMockUser(ctrl)

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", "user deleted successfully", nil,
			func(ctx *gofr.Context) {
				userSvc.EXPECT().Delete(ctx, 1).Return(nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				userSvc.EXPECT().Delete(ctx, 1).Return(errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/user", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(userSvc)

			output, err := h.Delete(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := services.NewMockUser(ctrl)

	user := &models.User{ID: 1, FirstName: "Aditi", LastName: "Verma", Email: "vermaditi2020@gmail.com"}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", user, nil,
			func(ctx *gofr.Context) {
				userSvc.EXPECT().GetByID(ctx, 1).Return(user, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				userSvc.EXPECT().GetByID(ctx, 1).Return(nil, errors.New("error"))
			}},
		{"Failure Case: invalid id", "!", nil, errors.New("invalid id"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(userSvc)

			output, err := h.GetByID(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := services.NewMockUser(ctrl)

	user := &models.User{ID: 1, FirstName: "Aditi", LastName: "Verma", Email: "vermaditi2020@gmail.com"}

	tests := []struct {
		description    string
		id             string
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", "1", []*models.User{user}, nil,
			func(ctx *gofr.Context) {
				userSvc.EXPECT().GetAll(ctx, &filters.User{}).Return([]*models.User{user}, nil)
			}},
		{"Failure Case: Error from service layer", "1", nil, errors.New("error"),
			func(ctx *gofr.Context) {
				userSvc.EXPECT().GetAll(ctx, &filters.User{}).Return(nil, errors.New("error"))
			}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.Header.Set("Content-Type", "application/json")

			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(userSvc)

			output, err := h.GetAll(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}
