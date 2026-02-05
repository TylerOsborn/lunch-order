package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"lunchorder/queries"
)

type MealOrderRepository struct {
	db             *sqlx.DB
	userRepository *UserRepository
	mealRepository *MealRepository
}

func NewMealOrderRepository(db *sqlx.DB, userRepository *UserRepository, mealRepository *MealRepository) *MealOrderRepository {
	return &MealOrderRepository{
		db:             db,
		userRepository: userRepository,
		mealRepository: mealRepository,
	}
}

func (r *MealOrderRepository) CreateMealOrder(order *MealOrder) error {
	_, err := r.db.Exec(queries.CreateMealOrder,
		order.UserID,
		order.WeekStartDate,
		order.MondayMealID,
		order.TuesdayMealID,
		order.WednesdayMealID,
		order.ThursdayMealID,
	)
	return err
}

func (r *MealOrderRepository) GetMealOrderByUserAndWeek(userID uint, weekStartDate string) (*MealOrder, error) {
	var order MealOrder
	err := r.db.Get(&order, queries.GetMealOrderByUserAndWeek, userID, weekStartDate)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Load meal details using GetMealByID for better performance
	if order.MondayMealID != nil {
		meal, err := r.mealRepository.GetMealByID(*order.MondayMealID)
		if err == nil {
			order.MondayMeal = meal
		}
	}

	if order.TuesdayMealID != nil {
		meal, err := r.mealRepository.GetMealByID(*order.TuesdayMealID)
		if err == nil {
			order.TuesdayMeal = meal
		}
	}

	if order.WednesdayMealID != nil {
		meal, err := r.mealRepository.GetMealByID(*order.WednesdayMealID)
		if err == nil {
			order.WednesdayMeal = meal
		}
	}

	if order.ThursdayMealID != nil {
		meal, err := r.mealRepository.GetMealByID(*order.ThursdayMealID)
		if err == nil {
			order.ThursdayMeal = meal
		}
	}

	return &order, nil
}

func (r *MealOrderRepository) GetMealOrderByID(id uint) (*MealOrder, error) {
	var order MealOrder
	err := r.db.Get(&order, queries.GetMealOrderByID, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &order, err
}
