package repository

import (
	"github.com/jmoiron/sqlx"
	"lunchorder/queries"
)

type DonationRequestRepository struct {
	db                 *sqlx.DB
	userRepository     *UserRepository
	donationRepository *DonationRepository
}

var donationRequestRepository *DonationRequestRepository

func NewDonationRequestRepository(db *sqlx.DB, userRepository *UserRepository, donationRepository *DonationRepository) *DonationRequestRepository {
	return &DonationRequestRepository{
		db:                 db,
		userRepository:     userRepository,
		donationRepository: donationRepository,
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

	if err := tx.Commit(); err != nil {
		return err
	}

	return r.CheckAndFulfillDonationRequests()
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
			&dr.ID, &dr.CreatedAt, &dr.UpdatedAt, &dr.DeletedAt, &dr.RequesterID, &dr.Status, &dr.DonationID,
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
			&dr.ID, &dr.CreatedAt, &dr.UpdatedAt, &dr.DeletedAt, &dr.RequesterID, &dr.Status, &dr.DonationID,
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

		if len(preferredMeals) == 0 {
			continue
		}
		mealDate := preferredMeals[0].Date

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

		// Fulfill
		tx, err := r.db.Beginx()
		if err != nil {
			return err
		}

		// Update Donation
		_, err = tx.Exec(queries.ClaimDonation, request.RequesterID, matchingDonation.ID)
		if err != nil {
			tx.Rollback()
			continue
		}

		// Update Request
		_, err = tx.Exec(queries.UpdateRequestStatus, "fulfilled", matchingDonation.ID, request.ID)
		if err != nil {
			tx.Rollback()
			continue
		}

		tx.Commit()
	}

	return nil
}
