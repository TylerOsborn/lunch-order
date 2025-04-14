package repository

import (
	"gorm.io/gorm"
)

type DonationRequestRepository struct {
	db                 *gorm.DB
	userRepository     *UserRepository
	donationRepository *DonationRepository
}

var donationRequestRepository *DonationRequestRepository

func NewDonationRequestRepository(db *gorm.DB, userRepository *UserRepository, donationRepository *DonationRepository) *DonationRequestRepository {
	if donationRequestRepository == nil {
		donationRequestRepository = &DonationRequestRepository{
			db:                 db,
			userRepository:     userRepository,
			donationRepository: donationRepository,
		}
	}

	return donationRequestRepository
}

func (r *DonationRequestRepository) CreateDonationRequest(requesterID uint, mealIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		request := DonationRequest{
			RequesterID: requesterID,
			Status:      "pending",
		}

		if err := tx.Create(&request).Error; err != nil {
			return err
		}

		for _, mealID := range mealIDs {
			preference := DonationRequestMeal{
				DonationRequestID: request.ID,
				MealID:            mealID,
			}

			if err := tx.Create(&preference).Error; err != nil {
				return err
			}
		}

		return r.CheckAndFulfillDonationRequests()
	})
}

func (r *DonationRequestRepository) GetDonationRequestsByStatus(status string) ([]DonationRequest, error) {
	var requests []DonationRequest

	result := r.db.Preload("Requester").
		Where("status = ?", status).
		Order("created_at ASC").
		Find(&requests)

	return requests, result.Error
}

func (r *DonationRequestRepository) UpdateDonationRequestStatus(requestID uint, status string, donationID *uint) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if donationID != nil {
		updates["donation_id"] = donationID
	}

	result := r.db.Model(&DonationRequest{}).
		Where("id = ?", requestID).
		Updates(updates)

	return result.Error
}

func (r *DonationRequestRepository) GetDonationRequestsByRequesterName(requesterName string, date string) ([]DonationRequest, error) {
	var user User
	if err := r.db.Where("name = ?", requesterName).First(&user).Error; err != nil {
		return nil, err
	}

	var requests []DonationRequest
	result := r.db.Joins("Requester").Joins("Donation").Joins("Donation.Donor").Joins("Donation.Meal").
		Where("requester_id = ? AND status = 'pending' AND DATE(donation_requests.created_at) = DATE(?) ", user.ID, date).
		Order("created_at DESC").
		Find(&requests)

	return requests, result.Error
}

func (r *DonationRequestRepository) GetDonationRequestById(id uint) (DonationRequest, error) {
	var request DonationRequest
	result := r.db.Preload("Requester").Preload("Donation").Preload("Donation.Donor").Preload("Donation.Meal").
		First(&request, id)

	return request, result.Error
}

func (r *DonationRequestRepository) GetDonationRequestMealPreferences(requestID uint) ([]Meal, error) {
	var meals []Meal

	result := r.db.Table("meals").
		Joins("JOIN donation_request_meals ON meals.id = donation_request_meals.meal_id").
		Where("donation_request_meals.donation_request_id = ?", requestID).
		Find(&meals)

	return meals, result.Error
}

func (r *DonationRequestRepository) CheckAndFulfillDonationRequests() error {

	pendingRequests, err := r.GetDonationRequestsByStatus("pending")
	if err != nil {
		return err
	}

	for _, request := range pendingRequests {
		preferredMeals, err := r.GetDonationRequestMealPreferences(request.ID)
		if err != nil || len(preferredMeals) == 0 {
			continue
		}

		var preferredMealIDs []uint
		for _, meal := range preferredMeals {
			preferredMealIDs = append(preferredMealIDs, meal.ID)
		}

		var mealDate string
		if err := r.db.Model(&Meal{}).Where("id = ?", preferredMealIDs[0]).Select("date").Scan(&mealDate).Error; err != nil {
			continue
		}

		unclaimedDonations, err := r.donationRepository.GetUnclaimedDonationsByDate(mealDate)
		if err != nil || len(unclaimedDonations) == 0 {
			continue
		}

		var matchingDonation *Donation
		for _, donation := range unclaimedDonations {
			for _, preferredMealID := range preferredMealIDs {
				if donation.MealID == preferredMealID {
					matchingDonation = &donation
					break
				}
			}
			if matchingDonation != nil {
				break
			}
		}

		if matchingDonation == nil {
			continue
		}

		err = r.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&Donation{}).
				Where("id = ?", matchingDonation.ID).
				Update("recipient_id", request.RequesterID).Error; err != nil {
				return err
			}

			if err := tx.Model(&DonationRequest{}).
				Where("id = ?", request.ID).
				Updates(map[string]interface{}{
					"status":      "fulfilled",
					"donation_id": matchingDonation.ID,
				}).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			continue
		}
	}

	return nil
}
