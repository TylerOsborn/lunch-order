package models

import "time"

type Meal struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UUID	  string    `json:"uuid" gorm:"unique"`
	Name      string    `json:"name" gorm:"unique not null"`
	CreatedAt time.Time `json:"createdAt"`
}

type APIDonation struct {
	MealID    uint   `json:"mealId"`
	User      User   `json:"user"`
}

type APIRecipient struct {
	DonationID uint   `json:"donationId"`
	User       User   `json:"user"`
}

type Donation struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	MealID      uint      `json:"mealId"`
	Meal        Meal      `json:"meal" gorm:"foreignKey:MealID"`
	DonorID     uint      `json:"donorId"`
	Donor       User      `json:"donor" gorm:"foreignKey:DonorID"`
	RecipientID uint      `json:"recipientId"`
	Recipient   User      `json:"recipient" gorm:"foreignKey:RecipientID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (donation *Donation) IsEmpty() bool {
	return donation.ID == 0
}

type UnclaimedDonation struct {
	ID          uint   `json:"id"`
	DonorName   string `json:"donorName"`
	Description string `json:"description"`
}

type ApiResult struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}

type DonationClaimSummary struct {
	Claimed       bool   `json:"claimed"`
	Description   string `json:"description"`
	DonorName     string `json:"donorName"`
	RecipientName string `json:"recipientName"`
}

type MealUpload struct {
	Csv string `json:"csv"`
}
