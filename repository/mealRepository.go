package repository

import (
	"github.com/jmoiron/sqlx"
	"lunchorder/queries"
)

type MealRepository struct {
	db *sqlx.DB
}

var mealRepository *MealRepository

func NewMealRepository(db *sqlx.DB) *MealRepository {
	return &MealRepository{
		db: db,
	}
}

func (r *MealRepository) CreateMeal(meal *Meal) error {
	var existingMeal Meal
	err := r.db.Get(&existingMeal, queries.GetMealByDescDate, meal.Description, meal.Date)

	if err == nil && existingMeal.ID != 0 {
		return nil // Already exists
	}

	_, err = r.db.Exec(queries.CreateMeal, meal.Description, meal.Date)
	return err
}

func (r *MealRepository) GetMealsByDate(date string) ([]Meal, error) {
	var meals []Meal
	err := r.db.Select(&meals, queries.GetMealsByDate, date)
	if err != nil {
		return nil, err
	}
	return meals, nil
}

func (r *MealRepository) GetMealsByDates(startDate string, endDate string) ([]Meal, error) {
	var meals []Meal
	err := r.db.Select(&meals, queries.GetMealsByRange, startDate, endDate)
	if err != nil {
		return nil, err
	}
	return meals, nil
}
