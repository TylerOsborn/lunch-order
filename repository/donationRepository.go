package repository

import (
	"log"
	"time"

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

func (r *DonationRepository) CreateDonation(donation *Donation) error {
	result := r.db.Create(&donation)

	return result.Error
}

func (r *DonationRepository) ClaimDonation(donationId uint, user *User) (bool, error) {

	result := r.db.Exec("UPDATE donations SET recipient_id = ? WHERE id = ? AND (recipient_id = 0 OR recipient_id IS NULL)", user.ID, donationId)
	if result.Error != nil {
		log.Println("Failed to claim donation:", result.Error)
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *DonationRepository) GetUnclaimedDonationsByDate(date string) ([]Donation, error) {
	var unclaimedDonations []Donation

	tx := r.db.Joins("Donor").Joins("Recipient").Joins("Meal").Where("(recipient_id <= 0 OR recipient_id IS NULL) AND DATE(Meal.date) = DATE(?)", date).Find(&unclaimedDonations)

	if tx.Error != nil {
		return unclaimedDonations, tx.Error
	}

	return unclaimedDonations, nil
}

func (r *DonationRepository) GetDonationsSummaryByDate(date string) (*[]Donation, error) {
	var donations []Donation
	tx := r.db.Joins("Donor").Joins("Recipient").Joins("Meal").Where("DATE(donations.created_at) = DATE(?)", date).Find(&donations)

	return &donations, tx.Error
}

func (r *DonationRepository) GetDonationClaimByClaimantName(name string) (Donation, error) {
	donation := Donation{}
	tx := r.db.Joins("Donor").Joins("Recipient").Joins("Meal").Where("recipient_id = (SELECT id FROM users WHERE name = ?) AND DATE(donations.created_at) = DATE(?)", name, time.Now()).First(&donation)

	if tx.Error != nil {
		return donation, tx.Error
	}

	return donation, nil
}

func (r *DonationRepository) GetDonationByDonorName(name string) (Donation, error) {
	donation := Donation{}
	tx := r.db.Joins("Donor").Joins("Recipient").Joins("Meal").Where("donor_id = (SELECT id FROM users WHERE name = ?) AND DATE(donations.created_at) = DATE(?)", name, time.Now()).First(&donation)

	if tx.Error != nil {
		return donation, tx.Error
	}

	return donation, nil
}
