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
