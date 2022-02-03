package repository

import (
	"github.com/jmoiron/sqlx"
	"test-api/pkg/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUserByCredentials(username, hashPassword string) (entity.User, error)
}

type TodoList interface {
	Create(todoList entity.TodoList, user int) (int, error)
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoItem
	TodoList
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoPostgres(db),
	}
}
