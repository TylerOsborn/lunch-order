package repository

import (
	"lunchorder/models"

	"gorm.io/gorm"
)

type DonationRequestRepository struct {
	db *gorm.DB
}

var donationRequestRepository *DonationRequestRepository

func NewDonationRequestRepository(db *gorm.DB) *DonationRequestRepository {
	if donationRequestRepository == nil {
		donationRequestRepository = &DonationRequestRepository{db: db}
	}

	return donationRequestRepository
}

func (r *DonationRequestRepository) CreateDonationRequest(donationRequest *models.DonationRequest) error {
	result := r.db.Create(donationRequest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *DonationRequestRepository) Save(donationRequest *models.DonationRequest) error {
	result := r.db.Save(donationRequest)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *DonationRepository) GetRequestByMealID(mealID uint) (*models.DonationRequest, error) {
	var donationRequest models.DonationRequest
	result := r.db.Where("meal_id = ?", mealID).First(&donationRequest)
	if result.Error != nil {
		return nil, result.Error
	}

	return &donationRequest, nil
}