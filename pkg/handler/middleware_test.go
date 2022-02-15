package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"test-api/pkg/service"
	mock_service "test-api/pkg/service/mocks"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Empty header",
			headerName:           "",
			token:                "token",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"StatusCode":401,"error":"Empty Authorization header"}`,
		},
		{
			name:                 "Invalid header value",
			headerName:           "Authorization",
			headerValue:          "token",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"StatusCode":401,"error":"Invalid auth header"}`,
		},
		{
			name:                 "Invalid token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"StatusCode":401,"error":"Invalid auth header"}`,
		},
		{
			name:                 "Invalid bearer",
			headerName:           "Authorization",
			headerValue:          "Bear token",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"StatusCode":401,"error":"Invalid auth header"}`,
		},
		{
			name:        "Service failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, errors.New("parse token error"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"StatusCode":401,"error":"Parse token error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			srv := gin.New()
			srv.GET("/auth/sign-in", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			resWriter := httptest.NewRecorder()
			req := httptest.NewRequest(
				"GET",
				"/auth/sign-in",
				nil,
			)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			srv.ServeHTTP(resWriter, req)

			assert.Equal(t, resWriter.Code, testCase.expectedStatusCode)
			assert.Equal(t, resWriter.Body.String(), testCase.expectedResponseBody)
		})
	}
}
