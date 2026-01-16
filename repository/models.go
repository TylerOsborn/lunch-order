package repository

import "time"

type Meal struct {
	ID          uint       `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	Description string     `json:"description" db:"description"`
	Date        string     `json:"date" db:"date"`
}

type User struct {
	ID        uint       `db:"id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
	Name      string     `json:"Name" db:"name"`
}

type Donation struct {
	ID          uint       `db:"id"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
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
	DeletedAt   *time.Time `db:"deleted_at"`
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