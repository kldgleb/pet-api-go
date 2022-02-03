package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"test-api/pkg/entity"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	var id int
	query := fmt.Sprintf(
		"INSERT INTO %s (name, username, password_hash) values ($1,$2,$3) RETURNING id",
		usersTable,
	)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByCredentials(username, hashPassword string) (entity.User, error) {
	var user entity.User
	query := fmt.Sprintf("SELECT id, name, username FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, hashPassword)
	if err != nil {
		return user, err
	}
	return user, nil
}
