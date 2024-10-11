package repository

import (
	"log"
	"lunchorder/models"

	"gorm.io/gorm"
)

type DonationRepository struct {
	db *gorm.DB
	userRepository *UserRepository
}

var donationRepository *DonationRepository

func NewDonationRepository(db *gorm.DB, userRepository *UserRepository) *DonationRepository {
	if donationRepository == nil {
		donationRepository = &DonationRepository{
			db: db,
			userRepository: userRepository,
		}	
	}

	return donationRepository
}

func (r *DonationRepository) CreateDonation(donationRequest *models.APIDonation) error {
	var donation models.Donation

	donor, err := r.userRepository.GetUserByName(donationRequest.DonorName)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	
	if err == gorm.ErrRecordNotFound {
		donor = &models.User{Name: donationRequest.DonorName}
		err = r.userRepository.CreateUser(donor)
	}

	if err != nil {
		return err
	}

	donation.DonorID = donor.ID
	donation.MealID = donationRequest.MealID
	
	result := r.db.Create(&donation)

	return result.Error
}

func (r *DonationRepository) ClaimDonation(donationId uint, user *models.User) (bool, error) {

	result := r.db.Exec("UPDATE donations SET recipient_id = ? WHERE id = ? AND (recipient_id = 0 OR recipient_id IS NULL)", user.ID, donationId)
	if result.Error != nil {
		log.Println("Failed to claim donation:", result.Error)
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *DonationRepository) GetUnclaimedDonationsByDate(date string) ([]models.UnclaimedDonation, error) {
	var unclaimedDonations []models.UnclaimedDonation

	result := r.db.Raw("SELECT DISTINCT meals.description, users.name AS donor_name, donations.id FROM donations INNER JOIN meals ON donations.meal_id = meals.id INNER JOIN users ON donations.donor_id = users.id WHERE (recipient_id <= 0  OR recipient_id IS NULL) AND meals.date = ? GROUP BY description ORDER BY description", date).Scan(&unclaimedDonations)

	if result.Error != nil {	
		return nil, result.Error
	}

	return unclaimedDonations, nil
}

func (r *DonationRepository) GetUnclaimedDonations() ([]models.UnclaimedDonation, error) {

	var unclaimedDonations []models.UnclaimedDonation

	result := r.db.Raw("SELECT DISTINCT description, name, id FROM donations INNER JOIN meals ON donations.meal_id = meal.id WHERE recipient_id <= 0  OR recipient IS NULL GROUP BY description ORDER BY description").Scan(&unclaimedDonations)
	
	if result.Error != nil {
		return nil, result.Error
	}

	return unclaimedDonations, nil
}