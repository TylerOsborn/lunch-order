package repository

import (
	"github.com/jmoiron/sqlx"
	"lunchorder/queries"
)

type UserRepository struct {
	db *sqlx.DB
}

var userRepo *UserRepository

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(user *User) error {
	result, err := r.db.Exec(queries.CreateUser, user.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint(id)
	return nil
}

func (r *UserRepository) GetUserByName(name string) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByName, name)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
