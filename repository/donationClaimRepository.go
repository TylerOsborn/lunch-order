package repository

import (
	"database/sql"
	"lunchorder/models"
)


type DonationClaimRepository struct {
	db *sql.DB
}

var donationClaimRepository *DonationClaimRepository

func NewDonationClaimRepository(db *sql.DB) *DonationClaimRepository {
	if donationClaimRepository == nil {
		donationClaimRepository = &DonationClaimRepository{
			db: db,
		}
	}

	return donationClaimRepository
}

func (r *DonationClaimRepository) CreateDonationClaim(donationClaim *models.DonationClaim) error {
	_, err := r.db.Exec("INSERT INTO donation_claim (donation_id, name) VALUES (?, ?)", donationClaim.DonationId, donationClaim.Name)
	if err != nil {
		return err
	}

	return nil
}
