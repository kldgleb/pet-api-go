package service

import (
	"test-api/pkg/entity"
	"test-api/pkg/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetJWTByCredentials(input entity.SignInInput) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	Create(list entity.TodoList, userId int) (int, error)
	GetAllLists(userId int) ([]entity.TodoList, error)
	GetListById(userId, listId int) (entity.TodoList, error)
	UpdateList(userId, listId int, list entity.TodoList) error
	DeleteList(userId, listId int) error
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
	}
}
