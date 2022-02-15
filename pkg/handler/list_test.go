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

func TestHandler_createList(t *testing.T) {
	type mockBehavior func(s *mock_service.MockTodoList, user entity.TodoList, userId int)

	testTable := []struct {
		name                string
		inputBody           string
		inputTodo           entity.TodoList
		userId              int
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"title":"title","description":"desc"}`,
			inputTodo: entity.TodoList{
				Title:       "title",
				Description: "desc",
			},
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoList, user entity.TodoList, userId int) {
				s.EXPECT().Create(user, userId).Return(1, nil)
			},
			expectedStatusCode:  201,
			expectedRequestBody: `{"listId":1}`,
		},
		{
			name:                "Bad request",
			inputBody:           `{"ttle":"title"}`,
			inputTodo:           entity.TodoList{},
			userId:              1,
			mockBehavior:        func(s *mock_service.MockTodoList, user entity.TodoList, userId int) {},
			expectedStatusCode:  400,
			expectedRequestBody: `{"StatusCode":400,"error":"Key: 'TodoList.Title' Error:Field validation for 'Title' failed on the 'required' tag"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"title":"title","description":"desc"}`,
			inputTodo: entity.TodoList{
				Title:       "title",
				Description: "desc",
			},
			userId: 1,
			mockBehavior: func(s *mock_service.MockTodoList, user entity.TodoList, userId int) {
				s.EXPECT().Create(user, userId).Return(1, errors.New("internal server error"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"StatusCode":500,"error":"error while creating todoList"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock_service.NewMockTodoList(c)
			testCase.mockBehavior(mockService, testCase.inputTodo, testCase.userId)

			services := &service.Service{
				TodoList: mockService,
			}
			handler := NewHandler(services)

			srv := gin.New()
			srv.POST("/api/lists", func(ctx *gin.Context) {
				ctx.Set(userCtx, testCase.userId)
			}, handler.createList)

			resWriter := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				"/api/lists",
				bytes.NewBufferString(testCase.inputBody),
			)
			srv.ServeHTTP(resWriter, req)

			assert.Equal(t, resWriter.Code, testCase.expectedStatusCode)
			assert.Equal(t, resWriter.Body.String(), testCase.expectedRequestBody)
		})
	}
}

func TestHandler_getAllLists(t *testing.T) {
}

func TestHandler_getListById(t *testing.T) {
}

func TestHandler_updateList(t *testing.T) {
}

func TestHandler_deleteList(t *testing.T) {
}
