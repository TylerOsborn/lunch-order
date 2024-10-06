package repository

import (
	"database/sql"
	"log"
	. "lunchorder/models"
)

type DonationRepository struct {
	db *sql.DB
}

var donationRepository *DonationRepository

func NewDonationRepository(db *sql.DB) *DonationRepository {
	if donationRepository == nil {
		donationRepository = &DonationRepository{
			db: db,
		}	
	}

	return donationRepository
}

func (r *DonationRepository) CreateDonation(donation *Donation) error {
	_, err := r.db.Exec("INSERT INTO donation (name, description) VALUES (?, ?)", donation.Name, donation.Description)
	if err != nil {
		log.Println("Failed to create donation:", err)
		return err
	}

	return nil
}

func (r *DonationRepository) ClaimDonation(donationId int) (bool, error) {
	result, err := r.db.Exec("UPDATE donation SET claimed = 1 WHERE id = ?", donationId)
	if err != nil {
		log.Println("Failed to claim donation:", err)
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *DonationRepository) GetUnclaimedDonationsByDate(date string) ([]Donation, error) {
	rows, err := r.db.Query("SELECT DISTINCT description, name, id FROM donation WHERE claimed = 0 AND date = ? GROUP BY description ORDER BY description", date)
	if err != nil {
		log.Println("Failed to get unclaimed donations:", err)
		return nil, err
	}
	defer rows.Close()

	var donations []Donation

	for rows.Next() {
		var donation Donation
		err := rows.Scan(&donation.Description, &donation.Name, &donation.Id)
		if err != nil {
			return nil, err
		}
		donations = append(donations, donation)
	}


	return donations, nil
}

func (r *DonationRepository) GetUnclaimedDonations() ([]Donation, error) {
	rows, err := r.db.Query("SELECT DISTINCT description, name, id FROM donation WHERE claimed = 0 GROUP BY description ORDER BY description")
	if err != nil {
		log.Println("Failed to get unclaimed donations:", err)
		return nil, err
	}
	defer rows.Close()

	var donations []Donation

	for rows.Next() {
		var donation Donation
		err := rows.Scan(&donation.Description, &donation.Name, &donation.Id)
		if err != nil {
			return nil, err
		}
		donations = append(donations, donation)
	}

	return donations, nil
}