package service

import (
	"lunchorder/models"
	"lunchorder/repository"
	"lunchorder/utils"
)

type MealService struct {
	mealRepository *repository.MealRepository
}

var mealService *MealService

func NewMealService(mealRepository *repository.MealRepository) *MealService {
	if mealService == nil {
		mealService = &MealService{
			mealRepository: mealRepository,
		}
	}

	return mealService
}

func (service *MealService) GetMealsByDates(start string, end string) ([]models.Meal, error) {
	return service.mealRepository.GetMealsByDates(start, end)
}

func (service *MealService) GetMealsByDate(today string) ([]models.Meal, error) {
	return service.mealRepository.GetMealsByDate(today)
}

func (service *MealService) CreateMeal(meal *models.Meal) error {
	return service.mealRepository.CreateMeal(meal)
}

func (service *MealService) CreateMeals(mealUpload models.MealUpload) error {
	csvString := mealUpload.Csv
	records, err := utils.ParseCSV(csvString)
	if err != nil {
		return err
	}

	for _, record := range records {
		if len(record) != 2 {
			return utils.ErrIncorrectCSVFormat
		}
	}

	for _, record := range records {
		date, description := record[0], record[1]
		err := service.CreateMeal(&models.Meal{Date: date, Description: description})
		if err != nil {
			return err
		}
	}

	return nil
}
