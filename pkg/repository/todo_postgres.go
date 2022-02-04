package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"test-api/pkg/entity"
)

type TodoPostgres struct {
	db *sqlx.DB
}

func NewTodoPostgres(db *sqlx.DB) *TodoPostgres {
	return &TodoPostgres{db: db}
}

func (r *TodoPostgres) Create(todoList entity.TodoList, userId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, todoList.Title, todoList.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoPostgres) GetAllLists(userId int) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable,
		usersListsTable,
	)
	err := r.db.Select(&lists, query, userId)
	return lists, err
}

func (r *TodoPostgres) GetListById(userId, listId int) (entity.TodoList, error) {
	var list entity.TodoList
	query := fmt.Sprintf(
		`SELECT tl.id, tl.title, tl.description FROM %s tl
				INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable,
		usersListsTable,
	)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}

func (r *TodoPostgres) UpdateList(userId, listId int, list entity.TodoList) error {
	query := fmt.Sprintf(
		`UPDATE %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2
				SET tl.title = $3 AND tl.description = $4`,
		todoListsTable,
		usersListsTable,
	)
	_, err := r.db.Exec(query, userId, listId, list.Title, list.Description)
	return err
}

func (r *TodoPostgres) DeleteList(userId, listId int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`,
		todoListsTable,
		usersListsTable,
	)
	_, err := r.db.Exec(query, userId, listId)
	return err
}
