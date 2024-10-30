package repository

import (
	"log"
	"lunchorder/models"

	"gorm.io/gorm"
)

type DonationRepository struct {
	db             *gorm.DB
	userRepository *UserRepository
}

var donationRepository *DonationRepository

func NewDonationRepository(db *gorm.DB, userRepository *UserRepository) *DonationRepository {
	if donationRepository == nil {
		donationRepository = &DonationRepository{
			db:             db,
			userRepository: userRepository,
		}
	}

	return donationRepository
}

func (r *DonationRepository) CreateDonation(donation *models.Donation) error {
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

func (r *DonationRepository) GetDonationsSummaryByDate(date string, donationClaimSummaries *[]models.DonationClaimSummary) error {
	tx := r.db.Raw("SELECT donations.recipient_id > 0 AS claimed, meals.description AS description, donors.name AS donor_name, COALESCE(recipients.name, 'UNCLAIMED') AS recipient_name FROM donations INNER JOIN users donors ON donors.id = donations.donor_id LEFT JOIN users recipients ON recipients.id = donations.recipient_id INNER JOIN meals ON donations.meal_id = meals.id WHERE meals.date = ?", date).Scan(&donationClaimSummaries)

	return tx.Error
}
