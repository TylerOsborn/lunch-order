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

	// Load meal details
	if order.MondayMealID != nil {
		meals, err := r.mealRepository.GetMealsByDate("") // Will need to get by ID
		if err == nil && len(meals) > 0 {
			for _, meal := range meals {
				if meal.ID == *order.MondayMealID {
					order.MondayMeal = &meal
					break
				}
			}
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
