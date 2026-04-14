package repository

import (
	"github.com/jmoiron/sqlx"
	"lunchorder/queries"
)

type DonationRequestRepository struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

var donationRequestRepository *DonationRequestRepository

func NewDonationRequestRepository(db *sqlx.DB, userRepository *UserRepository) *DonationRequestRepository {
	return &DonationRequestRepository{
		db:             db,
		userRepository: userRepository,
	}
}

func (r *DonationRequestRepository) CreateDonationRequest(requesterID uint, mealIDs []uint) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create Request
	result, err := tx.Exec(queries.CreateDonationRequest, requesterID, "pending")
	if err != nil {
		return err
	}
	requestID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Create Meals
	for _, mealID := range mealIDs {
		_, err := tx.Exec(queries.CreateDonationRequestMeal, requestID, mealID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *DonationRequestRepository) GetDonationRequestsByStatus(status string) ([]DonationRequest, error) {
	var requests []DonationRequest

	rows, err := r.db.Queryx(queries.GetRequestsByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dr DonationRequest
		var u User
		
		err := rows.Scan(
			&dr.ID, &dr.CreatedAt, &dr.UpdatedAt, &dr.RequesterID, &dr.Status, &dr.DonationID,
			&u.ID, &u.Name,
		)
		if err != nil {
			return nil, err
		}
		dr.Requester = u
		requests = append(requests, dr)
	}

	return requests, nil
}

func (r *DonationRequestRepository) UpdateDonationRequestStatus(requestID uint, status string, donationID *uint) error {
	_, err := r.db.Exec(queries.UpdateRequestStatus, status, donationID, requestID)
	return err
}

func (r *DonationRequestRepository) GetDonationRequestsByRequesterName(requesterName string, date string) ([]DonationRequest, error) {
	user, err := r.userRepository.GetUserByName(requesterName)
	if err != nil {
		return nil, err
	}

	var requests []DonationRequest
	rows, err := r.db.Queryx(queries.GetRequestsByRequester, user.ID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var dr DonationRequest
		var u User
		var d Donation
		var m Meal
		var donor User
		
		var donationID *uint
		var donationMealID *uint
		var donationDonorID *uint
		var donorID *uint
		var donorName *string
		var mealID *uint
		var mealDesc *string
		var mealDate *string

		err := rows.Scan(
			&dr.ID, &dr.CreatedAt, &dr.UpdatedAt, &dr.RequesterID, &dr.Status, &dr.DonationID,
			&u.ID, &u.Name,
			&donationID, &donationMealID, &donationDonorID,
			&donorID, &donorName,
			&mealID, &mealDesc, &mealDate,
		)
		if err != nil {
			return nil, err
		}
		
		dr.Requester = u
		
		if donationID != nil {
			d.ID = *donationID
			d.MealID = *donationMealID
			d.DonorID = *donationDonorID
			
			if donorID != nil {
				donor.ID = *donorID
				donor.Name = *donorName
				d.Donor = donor
			}
			if mealID != nil {
				m.ID = *mealID
				m.Description = *mealDesc
				m.Date = *mealDate
				d.Meal = m
			}
			dr.Donation = d
		}
		
		requests = append(requests, dr)
	}

	return requests, nil
}

func (r *DonationRequestRepository) GetDonationRequestMealPreferences(requestID uint) ([]Meal, error) {
	var meals []Meal
	err := r.db.Select(&meals, queries.GetRequestMeals, requestID)
	return meals, err
}

func (r *DonationRequestRepository) FulfillRequest(requestID uint, donationID uint, requesterID uint) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(queries.ClaimDonation, requesterID, donationID); err != nil {
		return err
	}

	if _, err = tx.Exec(queries.UpdateRequestStatus, "fulfilled", donationID, requestID); err != nil {
		return err
	}

	return tx.Commit()
}
