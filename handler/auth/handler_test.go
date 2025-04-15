package auth

import (
	"bytes"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
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

func Test_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	authSvc := services.NewMockAuth(ctrl)
	userSvc := services.NewMockUser(ctrl)

	result := map[string]interface{}{
		"access_token":  "abc",
		"expires_in":    3599,
		"refresh_token": "abc",
		"scope":         "abc",
		"token_type":    "Bearer",
		"id_token":      "abc",
	}
	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", []byte(`{"code":"4/0Ab_5qlly5YnrQiAWnlB81H80T8um6wEOMVR4xWXqakKOqUpzz5Ow9yitDXUO4i_OGF0g7Q"}`), result, nil,
			func(ctx *gofr.Context) {
				authSvc.EXPECT().GenerateGoogleToken(ctx, "4/0Ab_5qlly5YnrQiAWnlB81H80T8um6wEOMVR4xWXqakKOqUpzz5Ow9yitDXUO4i_OGF0g7Q").Return(result, nil)
			}},
		{"Failure Case: Error from service layer", []byte(`{"code":"4/0Ab_5qlly5YnrQiAWnlB81H80T8um6wEOMVR4xWXqakKOqUpzz5Ow9yitDXUO4i_OGF0g7Q"}`), nil, errors.New("error"),
			func(ctx *gofr.Context) {
				authSvc.EXPECT().GenerateGoogleToken(ctx, "4/0Ab_5qlly5YnrQiAWnlB81H80T8um6wEOMVR4xWXqakKOqUpzz5Ow9yitDXUO4i_OGF0g7Q").Return(nil, errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/google-token", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(authSvc, userSvc)

			output, err := h.CreateToken(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	authSvc := services.NewMockAuth(ctrl)
	userSvc := services.NewMockUser(ctrl)

	body := []byte(`{
  "provider" : "GOOGLE",
  "platform" : "WEB",
  "providerData" : {
    "token" : "eyJhbGciOiJSUzI1NiIsImtpZCI6ImMzN2RhNzVjOWZiZTE4YzJjZTkxMjViOWFhMWYzMDBkY2IzMWU4ZDkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2MzQ1NDAzMzA5MjYtdHB1Y3VqaTY1cjNiNW5sa2dwcTU0YXJnZ3FxY3JldDcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2MzQ1NDAzMzA5MjYtdHB1Y3VqaTY1cjNiNW5sa2dwcTU0YXJnZ3FxY3JldDcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDIwOTI5NDg4MzczNTg5NDMwNDYiLCJlbWFpbCI6InZlcm1hZGl0aTIwMjBAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJOYjBxeFoxd3VUdlNpQVdBN2pTc0tnIiwibmFtZSI6ImFkaXRpIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0xSeTdMckFBNTg3c1lxblFQdEtrbU05S1R5eFdsWmhVN29rdTlHNDByM3dlQk1fQTQ9czk2LWMiLCJnaXZlbl9uYW1lIjoiYWRpdGkiLCJpYXQiOjE3NDQ3MjU3ODMsImV4cCI6MTc0NDcyOTM4M30.Sg-MrxYbuL1hLvNiButSP-0cXI5r5TkhYz8tYafdWHA13D98cCJNCr5ajh9zieyZcmoFjdXLqJh_vsqAZqjZtNAQcpxUHXIiKFqiGtnPjzXW6t67tGnvk6I7_XyPryT_e43PxYPV2S8wCsT1Dt3g6UYpPwbHt1Fuh_Tr2N9wHHNLx-qOtSN3B7AvQqADdDm5I4zGmUr8Hd0FMD_ngWxJmxcZMisvgzwd78eCntMe59rRy2MA_L8feAvps_poUt7cA8kg2raStDSzddWx84BrJvaqA_yrubAIEY3yu47nncYwvurYnPgP-6klN2-rXy_10_BByB7cKmE_Vij70wnL2Q"
  }
}`)

	idToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImMzN2RhNzVjOWZiZTE4YzJjZTkxMjViOWFhMWYzMDBkY2IzMWU4ZDkiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhenAiOiI2MzQ1NDAzMzA5MjYtdHB1Y3VqaTY1cjNiNW5sa2dwcTU0YXJnZ3FxY3JldDcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2MzQ1NDAzMzA5MjYtdHB1Y3VqaTY1cjNiNW5sa2dwcTU0YXJnZ3FxY3JldDcuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMDIwOTI5NDg4MzczNTg5NDMwNDYiLCJlbWFpbCI6InZlcm1hZGl0aTIwMjBAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImF0X2hhc2giOiJOYjBxeFoxd3VUdlNpQVdBN2pTc0tnIiwibmFtZSI6ImFkaXRpIiwicGljdHVyZSI6Imh0dHBzOi8vbGgzLmdvb2dsZXVzZXJjb250ZW50LmNvbS9hL0FDZzhvY0xSeTdMckFBNTg3c1lxblFQdEtrbU05S1R5eFdsWmhVN29rdTlHNDByM3dlQk1fQTQ9czk2LWMiLCJnaXZlbl9uYW1lIjoiYWRpdGkiLCJpYXQiOjE3NDQ3MjU3ODMsImV4cCI6MTc0NDcyOTM4M30.Sg-MrxYbuL1hLvNiButSP-0cXI5r5TkhYz8tYafdWHA13D98cCJNCr5ajh9zieyZcmoFjdXLqJh_vsqAZqjZtNAQcpxUHXIiKFqiGtnPjzXW6t67tGnvk6I7_XyPryT_e43PxYPV2S8wCsT1Dt3g6UYpPwbHt1Fuh_Tr2N9wHHNLx-qOtSN3B7AvQqADdDm5I4zGmUr8Hd0FMD_ngWxJmxcZMisvgzwd78eCntMe59rRy2MA_L8feAvps_poUt7cA8kg2raStDSzddWx84BrJvaqA_yrubAIEY3yu47nncYwvurYnPgP-6klN2-rXy_10_BByB7cKmE_Vij70wnL2Q"
	claims := &models.GoogleClaims{
		"102092948837358943046",
		"vermaditi2020@gmail.com",
		"aditi",
		"https://lh3.googleusercontent.com/a/ACg8ocLRy7LrAA587sYqnQPtKkmM9KTyxWlZhU7oku9G40r3weBM_A4=s96-c",
		"aditi",
		"",
		0,
	}

	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InZlcm1hZGl0aTIwMjBAZ21haWwuY29tIiwiZXhwIjoxNzQ0ODEzMzY4LCJmYW1pbHlfbmFtZSI6IiIsImdpdmVuX25hbWUiOiJhZGl0aSIsIm5hbWUiOiJhZGl0aSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NMUnk3THJBQTU4N3NZcW5RUHRLa21NOUtUeXhXbFpoVTdva3U5RzQwcjN3ZUJNX0E0PXM5Ni1jIiwidXNlcklEIjoiMTAyMDkyOTQ4ODM3MzU4OTQzMDQ2In0.VwrR0CcBAbDH8TWjxSXSusJBC9A501src5_YRJEzftE"
	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", body, &models.Tokens{RefreshToken: refreshToken}, nil,
			func(ctx *gofr.Context) {
				authSvc.EXPECT().VerifyGoogleIDToken(ctx.Context, idToken).Return(claims, nil)
				authSvc.EXPECT().GenerateRefreshToken(claims).Return(refreshToken, nil)
			}},
		{"Failure Case: Error from service layer", body, nil, errors.New("error"),
			func(ctx *gofr.Context) {
				authSvc.EXPECT().VerifyGoogleIDToken(ctx.Context, idToken).Return(nil, errors.New("error"))
			}},
		{"Failure Case: Error from service layer", body, nil, errors.New("error"),
			func(ctx *gofr.Context) {
				authSvc.EXPECT().VerifyGoogleIDToken(ctx.Context, idToken).Return(claims, nil)
				authSvc.EXPECT().GenerateRefreshToken(claims).Return("", errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
		{"Failure Case: bind error", []byte(`{}`), nil, errors.New("missing id_token"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(authSvc, userSvc)

			output, err := h.Login(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}

func Test_Refresh(t *testing.T) {
	ctrl := gomock.NewController(t)
	authSvc := services.NewMockAuth(ctrl)
	userSvc := services.NewMockUser(ctrl)

	body := []byte(`{"refreshToken":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InZlcm1hZGl0aTIwMjBAZ21haWwuY29tIiwiZXhwIjoxNzQ0ODEzODE3LCJmYW1pbHlfbmFtZSI6IiIsImdpdmVuX25hbWUiOiJhZGl0aSIsIm5hbWUiOiJhZGl0aSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NMUnk3THJBQTU4N3NZcW5RUHRLa21NOUtUeXhXbFpoVTdva3U5RzQwcjN3ZUJNX0E0PXM5Ni1jIiwidXNlcklEIjoiMTAyMDkyOTQ4ODM3MzU4OTQzMDQ2In0.1c_Gs-albE_QuN4nicP1ynjhYeATrUCeuPS2koMXDTs"}`)

	refreshToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InZlcm1hZGl0aTIwMjBAZ21haWwuY29tIiwiZXhwIjoxNzQ0ODEzODE3LCJmYW1pbHlfbmFtZSI6IiIsImdpdmVuX25hbWUiOiJhZGl0aSIsIm5hbWUiOiJhZGl0aSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NMUnk3THJBQTU4N3NZcW5RUHRLa21NOUtUeXhXbFpoVTdva3U5RzQwcjN3ZUJNX0E0PXM5Ni1jIiwidXNlcklEIjoiMTAyMDkyOTQ4ODM3MzU4OTQzMDQ2In0.1c_Gs-albE_QuN4nicP1ynjhYeATrUCeuPS2koMXDTs"

	jwtClaims := jwt.MapClaims{
		"email":       "vermaditi2020@gmail.com",
		"exp":         1744813817,
		"family_name": "",
		"given_name":  "aditi",
		"name":        "aditi",
		"picture":     "https://lh3.googleusercontent.com/a/ACg8ocLRy7LrAA587sYqnQPtKkmM9KTyxWlZhU7oku9G40r3weBM_A4=s96-c",
		"userID":      "102092948837358943046",
	}

	claims := &models.GoogleClaims{
		"",
		"vermaditi2020@gmail.com",
		"aditi",
		"https://lh3.googleusercontent.com/a/ACg8ocLRy7LrAA587sYqnQPtKkmM9KTyxWlZhU7oku9G40r3weBM_A4=s96-c",
		"aditi",
		"",
		0,
	}

	accessToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InZlcm1hZGl0aTIwMjBAZ21haWwuY29tIiwiZXhwIjoxNzQ0NzI4MTA5LCJmYW1pbHlfbmFtZSI6IiIsImdpdmVuX25hbWUiOiJhZGl0aSIsIm5hbWUiOiJhZGl0aSIsInBpY3R1cmUiOiJodHRwczovL2xoMy5nb29nbGV1c2VyY29udGVudC5jb20vYS9BQ2c4b2NMUnk3THJBQTU4N3NZcW5RUHRLa21NOUtUeXhXbFpoVTdva3U5RzQwcjN3ZUJNX0E0PXM5Ni1jIiwidXNlcklEIjoxfQ.VxYeFcd5GSZyDhZ4InyKA8fm6bMWuiNXncjHwKuJhck"

	tests := []struct {
		description    string
		body           []byte
		expectedOutput interface{}
		expectedErr    error
		execMocks      func(ctx *gofr.Context)
	}{
		{"Success Case", body, &models.Tokens{AccessToken: accessToken}, nil,
			func(ctx *gofr.Context) {
				authSvc.EXPECT().ValidateRefreshToken(refreshToken).Return(jwtClaims, nil)
				userSvc.EXPECT().AuthAdaptor(ctx, claims).Return(nil)
				authSvc.EXPECT().GenerateAccessToken(claims).Return(accessToken, nil)
			}},
		{"Failure Case: Error from service layer", body, nil, errors.New("error"),
			func(ctx *gofr.Context) {
				authSvc.EXPECT().ValidateRefreshToken(refreshToken).Return(nil, errors.New("error"))
			}},
		{"Failure Case: Error from service layer", body, nil, errors.New("error"),
			func(ctx *gofr.Context) {
				authSvc.EXPECT().ValidateRefreshToken(refreshToken).Return(jwtClaims, nil)
				userSvc.EXPECT().AuthAdaptor(ctx, claims).Return(errors.New("error"))
			}},
		{"Failure Case: Error from service layer", body, nil, errors.New("error"),
			func(ctx *gofr.Context) {
				authSvc.EXPECT().ValidateRefreshToken(refreshToken).Return(jwtClaims, nil)
				userSvc.EXPECT().AuthAdaptor(ctx, claims).Return(nil)
				authSvc.EXPECT().GenerateAccessToken(claims).Return("", errors.New("error"))
			}},
		{"Failure Case: bind error", []byte(`{`), nil, errors.New("bind error"), func(ctx *gofr.Context) {
		}},
		{"Failure Case: bind error", []byte(`{}`), nil, errors.New("missing refresh token"), func(ctx *gofr.Context) {
		}},
	}

	for i, tc := range tests {
		i := i
		tc := tc

		t.Run(tc.description, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")
			gofrReq := gofrHTTP.NewRequest(req)
			mockContainer, _ := container.NewMockContainer(t)
			ctx := gofr.Context{Context: context.Background(), Request: gofrReq, Container: mockContainer, Out: nil}

			tc.execMocks(&ctx)

			h := New(authSvc, userSvc)

			output, err := h.Refresh(&ctx)

			assert.Equalf(t, tc.expectedOutput, output, "TEST[%d], failed.\n%s", i, tc.description)
			assert.Equalf(t, tc.expectedErr, err, "TEST[%d], failed.\n%s", i, tc.description)
		})
	}
}
