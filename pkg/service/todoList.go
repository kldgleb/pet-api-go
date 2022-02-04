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

func (s *TodoListService) GetAllLists(userId int) ([]entity.TodoList, error) {
	return s.repo.GetAllLists(userId)
}

func (s *TodoListService) GetListById(userId, listId int) (entity.TodoList, error) {
	return s.repo.GetListById(userId, listId)
}

func (s *TodoListService) UpdateList(userId, listId int, list entity.TodoList) error {
	return s.repo.UpdateList(userId, listId, list)
}

func (s *TodoListService) DeleteList(userId, listId int) error {
	return s.repo.DeleteList(userId, listId)
}
