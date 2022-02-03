package service

import (
	"test-api/pkg/entity"
	"test-api/pkg/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetJWTByCredentials(password, username string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type TodoList interface {
	Create(list entity.TodoList, userId int) (int, error)
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
