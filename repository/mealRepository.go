package repository

import (
	"lunchorder/models"

	"gorm.io/gorm"
)


type MealRepository struct {
	db *gorm.DB
}

var mealRepository *MealRepository

func NewMealRepository(db *gorm.DB) *MealRepository {
	if mealRepository == nil {
		mealRepository = &MealRepository{
			db: db,
		}
	}

	return mealRepository
}

func (r *MealRepository) CreateMeal(meal *models.Meal) error {
	var existingMeal models.Meal
	result := r.db.First(&existingMeal, "description = ? AND date = ?", meal.Description, meal.Date)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if existingMeal.ID != 0 {
		return nil
	}

	result = r.db.Create(meal)

	return result.Error
}

func (r *MealRepository) GetMealsByDate(date string) ([]models.Meal, error) {
	var meals []models.Meal

	result := r.db.Find(&meals, "date = ?", date)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return meals, nil
}

func (r *MealRepository) GetMealsByDates(startDate string, endDate string) ([]models.Meal, error) {

	var meals []models.Meal

	result := r.db.Where("date >= ? AND date <= ?", startDate, endDate).Find(&meals)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return meals, nil
}