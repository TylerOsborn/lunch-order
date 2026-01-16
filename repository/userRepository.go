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

func (r *UserRepository) GetUserByGoogleID(googleID string) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByGoogleID, googleID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id uint) (*User, error) {
	var user User
	err := r.db.Get(&user, queries.GetUserByID, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpsertUser(user *User) error {
	_, err := r.db.NamedExec(queries.UpsertUserGoogle, user)
	if err != nil {
		return err
	}
	// Fetch the user back to get the ID
	updatedUser, err := r.GetUserByGoogleID(*user.GoogleID)
	if err != nil {
		return err
	}
	user.ID = updatedUser.ID
	return nil
}
