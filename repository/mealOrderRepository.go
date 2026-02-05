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
	// Start a transaction
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert the meal order
	result, err := tx.Exec(queries.CreateMealOrder,
		order.UserID,
		order.WeekStartDate,
	)
	if err != nil {
		return err
	}

	// Get the inserted order ID
	orderID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Insert meal order items
	for _, item := range order.Items {
		_, err := tx.Exec(queries.CreateMealOrderItem,
			orderID,
			item.DayOfWeek,
			item.MealID,
		)
		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
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

	// Load meal order items
	var items []MealOrderItem
	err = r.db.Select(&items, queries.GetMealOrderItems, order.ID)
	if err != nil {
		return nil, err
	}

	// Load meal details for each item
	for i := range items {
		meal, err := r.mealRepository.GetMealByID(items[i].MealID)
		if err == nil {
			items[i].Meal = meal
		}
	}

	order.Items = items
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
