package service

import (
	"lunchorder/models"
	"lunchorder/repository"

	"github.com/google/uuid"
)


type UserService struct  {
	userRepository *repository.UserRepository
}

var userService *UserService

func NewUserService(userRepository *repository.UserRepository) *UserService {
	if userService == nil {
		userService = &UserService{
			userRepository: userRepository,
		}
	}

	return userService
}

func (service *UserService) CreateUser(user *models.User) error {
	return service.userRepository.CreateUser(user)
}

func (service *UserService) GetUserByName(name string) (*models.User, error) {
	return service.userRepository.GetUserByName(name)
}

func (service *UserService) GetUserByUUID(uuid string) (*models.User, error) {
	return service.userRepository.GetUserByUUID(uuid)
}

func (service *UserService) Save(user *models.User) error {
	return service.userRepository.Save(user)
}

func (service *UserService) SaveEmpty(user *models.User) error {
	if (user.UUID != "") {
		return nil
	}

	// pray this doesn't panic
	user.UUID = uuid.NewString()

	err := service.userRepository.CreateUser(user)

	return err
}
