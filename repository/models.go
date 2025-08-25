package repository

import "gorm.io/gorm"

type Meal struct {
	gorm.Model
	Description string `json:"description"`
	Date        string `json:"date"`
}

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	GoogleID string `json:"googleId" gorm:"unique;not null"`
	Role     string `json:"role" gorm:"default:'standard'"`
}

type Donation struct {
	gorm.Model
	MealID      uint  `json:"mealId"`
	Meal        Meal  `json:"meal" gorm:"foreignKey:MealID"`
	DonorID     uint  `json:"donorId"`
	Donor       User  `json:"donor" gorm:"foreignKey:DonorID"`
	RecipientID *uint `json:"recipientId"`
	Recipient   User  `json:"recipient" gorm:"foreignKey:RecipientID"`
}

type DonationRequest struct {
	gorm.Model
	RequesterID uint     `json:"requesterId"`
	Requester   User     `json:"requester" gorm:"foreignKey:RequesterID"`
	Status      string   `json:"status"` // "pending", "fulfilled", "cancelled"
	DonationID  *uint    `json:"donationId"`
	Donation    Donation `json:"donation" gorm:"foreignKey:DonationID"`
}

type DonationRequestMeal struct {
	DonationRequestID uint `gorm:"primaryKey"`
	MealID            uint `gorm:"primaryKey"`
	Meal              Meal
}
