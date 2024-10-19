package repository

import (
	"lunchorder/models"

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

func (r *UserRepository) CreateUser(user *models.User) error {
	result := r.db.Create(user)

	return result.Error
}

func (r *UserRepository) Save(user *models.User) error {
	result := r.db.Save(user)

	return result.Error
}


func (r *UserRepository) GetUserByName(name string) (*models.User, error) {
	var user models.User

	result := r.db.First(&user, "name = ?", name)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *UserRepository) GetUserByUUID(uuid string) (*models.User, error) {
	var user models.User

	if (uuid == "") {
		return nil, gorm.ErrRecordNotFound
	}

	result := r.db.First(&user, "uuid = ?", uuid)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}