package service

import (
	"test-api/pkg/entity"
	"test-api/pkg/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	//getJWTByCredentials(password, username string) (entity.User, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos repository.Authorization) *Service {
	return &Service{Authorization: NewAuthService(repos)}
}
