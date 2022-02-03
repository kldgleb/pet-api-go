package service

import (
	"test-api/pkg/entity"
	"test-api/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (s *TodoListService) Create(todoList entity.TodoList, userId int) (int, error) {
	return s.repo.Create(todoList, userId)
}
