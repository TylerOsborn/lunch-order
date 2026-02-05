package repository

import "time"

type Meal struct {
	ID          uint       `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	Description string     `json:"description" db:"description"`
	Date        string     `json:"date" db:"date"`
}

type User struct {
	ID                uint       `db:"id"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	Name              string     `json:"name" db:"name"`
	Email             *string    `json:"email" db:"-"`
	GoogleID          *string    `json:"googleId" db:"-"`
	EmailHash         *string    `json:"-" db:"email_hash"`
	EmailEncrypted    *string    `json:"-" db:"email_encrypted"`
	GoogleIDHash      *string    `json:"-" db:"google_id_hash"`
	GoogleIDEncrypted *string    `json:"-" db:"google_id_encrypted"`
	FirstName         *string    `json:"firstName" db:"first_name"`
	LastName          *string    `json:"lastName" db:"last_name"`
	AvatarURL         *string    `json:"avatarUrl" db:"avatar_url"`
	IsAdmin           bool       `json:"isAdmin" db:"is_admin"`
}

type Donation struct {
	ID          uint       `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	MealID      uint       `json:"mealId" db:"meal_id"`
	Meal        Meal       `json:"meal" db:"meal"`
	DonorID     uint       `json:"donorId" db:"donor_id"`
	Donor       User       `json:"donor" db:"donor"`
	RecipientID *uint      `json:"recipientId" db:"recipient_id"`
	Recipient   User       `json:"recipient" db:"recipient"`
}

type DonationRequest struct {
	ID          uint       `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	RequesterID uint       `json:"requesterId" db:"requester_id"`
	Requester   User       `json:"requester" db:"requester"`
	Status      string     `json:"status" db:"status"` // "pending", "fulfilled", "cancelled"
	DonationID  *uint      `json:"donationId" db:"donation_id"`
	Donation    Donation   `json:"donation" db:"donation"`
}

type DonationRequestMeal struct {
	DonationRequestID uint `db:"donation_request_id"`
	MealID            uint `db:"meal_id"`
	Meal              Meal `db:"meal"`
}

type MealOrder struct {
	ID            uint             `db:"id"`
	CreatedAt     time.Time        `db:"created_at"`
	UpdatedAt     time.Time        `db:"updated_at"`
	UserID        uint             `json:"userId" db:"user_id"`
	User          User             `json:"user" db:"user"`
	WeekStartDate string           `json:"weekStartDate" db:"week_start_date"`
	Items         []MealOrderItem  `json:"items" db:"-"`
}

type MealOrderItem struct {
	ID          uint      `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	MealOrderID uint      `db:"meal_order_id"`
	DayOfWeek   string    `db:"day_of_week"`
	MealID      uint      `db:"meal_id"`
	Meal        *Meal     `json:"meal" db:"-"`
}