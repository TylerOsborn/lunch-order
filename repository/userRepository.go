package repository

import (
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var userRepo *UserRepository

func NewUserRepository(db *gorm.DB) *UserRepository {
	if userRepo == nil {
		userRepo = &UserRepository{
			db: db,
		}
	}

	return userRepo
}

func (r *UserRepository) CreateUser(user *User) error {
	result := r.db.Create(user)

	return result.Error
}

func (r *UserRepository) GetUserByName(name string) (*User, error) {
	var user User

	result := r.db.First(&user, "name = ?", name)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
