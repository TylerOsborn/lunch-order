package repository

import (
	"database/sql"
	"lunchorder/constants"
	"lunchorder/models"
	"time"
)


type MealRepository struct {
	db *sql.DB
}

var mealRepository *MealRepository

func NewMealRepository(db *sql.DB) *MealRepository {
	if mealRepository == nil {
		mealRepository = &MealRepository{
			db: db,
		}
	}

	return mealRepository
}

func (r *MealRepository) CreateMeal(meal *models.Meal) error {
	_, err := r.db.Exec("INSERT INTO meal (date, description) VALUES (?, ?)", meal.Date, meal.Description)
	if err != nil {
		return err
	}

	return nil
}

func (r *MealRepository) GetMealsByDate(date string) ([]models.Meal, error) {
	rows, err := r.db.Query("SELECT id, description, date FROM meal WHERE date = ?", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal

	for rows.Next() {
		var meal models.Meal
		err := rows.Scan(&meal.Id, &meal.Description, &meal.Date)
		if err != nil {
			return nil, err
		}
		meal.Date = time.Now().Format(constants.DATE_FORMAT)
		meals = append(meals, meal)
	}

	return meals, nil
}

func (r *MealRepository) GetMealsByDates(startDate string, endDate string) ([]models.Meal, error) {

	rows, err := r.db.Query("SELECT id, description, date FROM meal WHERE date >= ? AND date <= ?", startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.Meal

	for rows.Next() {
		var meal models.Meal
		err := rows.Scan(&meal.Id, &meal.Description, &meal.Date)
		if err != nil {
			return nil, err
		}
		meal.Date = time.Now().Format(constants.DATE_FORMAT)
		meals = append(meals, meal)
	}

	return meals, nil
}