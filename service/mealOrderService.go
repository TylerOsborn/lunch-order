package service

import (
	"lunchorder/repository"
)

type MealOrderService struct {
	mealOrderRepository *repository.MealOrderRepository
	mealRepository      *repository.MealRepository
}

func NewMealOrderService(mealOrderRepository *repository.MealOrderRepository, mealRepository *repository.MealRepository) *MealOrderService {
	return &MealOrderService{
		mealOrderRepository: mealOrderRepository,
		mealRepository:      mealRepository,
	}
}

func (service *MealOrderService) CreateMealOrder(order *repository.MealOrder) error {
	return service.mealOrderRepository.CreateMealOrder(order)
}

func (service *MealOrderService) GetMealOrderByUserAndWeek(userID uint, weekStartDate string) (*repository.MealOrder, error) {
	return service.mealOrderRepository.GetMealOrderByUserAndWeek(userID, weekStartDate)
}
