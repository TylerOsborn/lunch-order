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

func (service *MealService) GetMealsByDates(start string, end string) ([]models.MealResponse, error) {
	var response []models.MealResponse
	meals, err := service.mealRepository.GetMealsByDates(start, end)

	if err != nil {
		return response, err
	}

	for _, meal := range meals {
		response = append(response, models.MealResponse{Date: meal.Date, Description: meal.Description})
	}
	return response, err
}

func (service *MealService) GetMealsByDate(today string) ([]models.MealResponse, error) {
	var results []models.MealResponse
	meals, err := service.mealRepository.GetMealsByDate(today)
	if err != nil {
		return results, err
	}

	for _, meal := range meals {
		results = append(results, models.MealResponse{ID: meal.ID, Date: meal.Date, Description: meal.Description})
	}

	return results, err
}

func (service *MealService) CreateMeal(meal *repository.Meal) error {
	return service.mealRepository.CreateMeal(meal)
}

func (service *MealService) CreateMeals(mealUpload models.MealUploadRequest) error {
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
		err := service.CreateMeal(&repository.Meal{Date: date, Description: description})
		if err != nil {
			return err
		}
	}

	return nil
}
