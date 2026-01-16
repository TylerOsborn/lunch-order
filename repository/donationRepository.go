package repository

import (
	"github.com/jmoiron/sqlx"
	"log"
	"lunchorder/queries"
	"time"
)

type DonationRepository struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

var donationRepository *DonationRepository

func NewDonationRepository(db *sqlx.DB, userRepository *UserRepository) *DonationRepository {
	return &DonationRepository{
		db:             db,
		userRepository: userRepository,
	}
}

func (r *DonationRepository) CreateDonation(donation *Donation) error {
	_, err := r.db.Exec(queries.CreateDonation, donation.MealID, donation.DonorID)
	return err
}

func (r *DonationRepository) ClaimDonation(donationId uint, user *User) (bool, error) {
	result, err := r.db.Exec(queries.ClaimDonation, user.ID, donationId)
	if err != nil {
		log.Println("Failed to claim donation:", err)
		return false, err
	}

	rows, err := result.RowsAffected()
	return rows > 0, err
}

func (r *DonationRepository) GetUnclaimedDonationsByDate(date string) ([]Donation, error) {
	var unclaimedDonations []Donation

	rows, err := r.db.Queryx(queries.GetUnclaimedDonations, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d Donation
		var m Meal
		var u User
		
		err := rows.Scan(
			&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.MealID, &d.DonorID, &d.RecipientID,
			&m.ID, &m.Description, &m.Date,
			&u.ID, &u.Name,
		)
		if err != nil {
			return nil, err
		}
		
		d.Meal = m
		d.Donor = u
		unclaimedDonations = append(unclaimedDonations, d)
	}

	return unclaimedDonations, nil
}

func (r *DonationRepository) GetDonationsSummaryByDate(date string) (*[]Donation, error) {
	var donations []Donation
	
	rows, err := r.db.Queryx(queries.GetDonationsSummary, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d Donation
		var m Meal
		var donor User
		var recipient User
		var recipientID *uint
		var recipientName *string

		err := rows.Scan(
			&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.MealID, &d.DonorID, &d.RecipientID,
			&m.ID, &m.Description, &m.Date,
			&donor.ID, &donor.Name,
			&recipientID, &recipientName,
		)
		
		if err != nil {
			return nil, err
		}

		d.Meal = m
		d.Donor = donor
		
		if recipientID != nil {
			recipient.ID = *recipientID
			if recipientName != nil {
				recipient.Name = *recipientName
			}
			d.Recipient = recipient
		}
		
		donations = append(donations, d)
	}

	return &donations, nil
}

func (r *DonationRepository) GetDonationClaimByClaimantName(name string) (Donation, error) {
	var d Donation
	var m Meal
	var donor User
	
	row := r.db.QueryRowx(queries.GetDonationClaimByName, name, time.Now())
	
	err := row.Scan(
		&d.ID, &d.CreatedAt, &d.UpdatedAt, &d.MealID, &d.DonorID, &d.RecipientID,
		&m.ID, &m.Description, &m.Date,
		&donor.ID, &donor.Name,
	)
	
	if err != nil {
		return d, err
	}
	
	d.Meal = m
	d.Donor = donor
	
	return d, nil
}
