package repository

import "gorm.io/gorm"

type Meal struct {
	gorm.Model
	Description string `json:"description"`
	Date        string `json:"date"`
}

type User struct {
	gorm.Model
	Name string `json:"Name" gorm:"unique not null"`
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
