package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
	"net/http/httptest"
	"test-api/pkg/entity"
	"test-api/pkg/service"
	mock_service "test-api/pkg/service/mocks"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user entity.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           entity.User
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"user","username":"test","password":"password"}`,
			inputUser: entity.User{
				Name:     "user",
				Username: "test",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:                "Empty fields",
			inputBody:           `{}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, user entity.User) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"StatusCode":400,"error":"invalid input body"}`,
		},
		{
			name:      "Service failure",
			inputBody: `{"name":"user","username":"test","password":"password"}`,
			inputUser: entity.User{
				Name:     "user",
				Username: "test",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user entity.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"StatusCode":500,"error":"service failure"}`,
		},
	}
	for i := range testTable {
		t.Run(testTable[i].name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testTable[i].mockBehavior(auth, testTable[i].inputUser)

			services := &service.Service{
				Authorization: auth,
			}
			handler := NewHandler(services)

			srv := gin.New()
			srv.POST("/auth/sign-up", handler.signUp)

			resWriter := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/auth/sign-up",
				bytes.NewBufferString(testTable[i].inputBody),
			)
			srv.ServeHTTP(resWriter, req)

			assert.Equal(t, resWriter.Code, testTable[i].expectedStatusCode)
			assert.Equal(t, resWriter.Body.String(), testTable[i].expectedRequestBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, input entity.SignInInput)
	testTable := []struct {
		name                string
		inputBody           string
		signInInput         entity.SignInInput
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"user","password":"password"}`,
			signInInput: entity.SignInInput{
				Username: "user",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input entity.SignInInput) {
				s.EXPECT().GetJWTByCredentials(input).Return("token", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"token"}`,
		},
		{
			name:                "Bad request",
			inputBody:           `{"password":"password"}`,
			mockBehavior:        func(s *mock_service.MockAuthorization, input entity.SignInInput) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"StatusCode":400,"error":"Key: 'SignInInput.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"user","password":"password"}`,
			signInInput: entity.SignInInput{
				Username: "user",
				Password: "password",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, input entity.SignInInput) {
				s.EXPECT().GetJWTByCredentials(input).Return("token", errors.New("internal server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"StatusCode":500,"error":"internal server error"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockAuth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(mockAuth, testCase.signInInput)
			services := &service.Service{
				Authorization: mockAuth,
			}
			handler := NewHandler(services)

			srv := gin.New()
			srv.POST("/auth/sign-in", handler.signIn)
			resWriter := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/auth/sign-in",
				bytes.NewBufferString(testCase.inputBody),
			)
			srv.ServeHTTP(resWriter, req)

			assert.Equal(t, resWriter.Code, testCase.expectedStatusCode)
			assert.Equal(t, resWriter.Body.String(), testCase.expectedRequestBody)
		})
	}
}
