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
	GetAllLists(userId int) ([]entity.TodoList, error)
	GetListById(userId, listId int) (entity.TodoList, error)
	UpdateList(userId, listId int, list entity.TodoList) error
	DeleteList(userId, listId int) error
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
