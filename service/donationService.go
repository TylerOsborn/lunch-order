package service

import (
	"database/sql"
	"errors"
	"lunchorder/models"
	"lunchorder/repository"
)

var ErrDonationNotFound = errors.New("donation not found")

type DonationService struct {
	donationRepository *repository.DonationRepository
	mealRepository     *repository.MealRepository
	userRepository     *repository.UserRepository
}

var donationService *DonationService

func NewDonationService(
	donationRepository *repository.DonationRepository,
	mealRepository *repository.MealRepository,
	userRepository *repository.UserRepository) *DonationService {

	return &DonationService{
		donationRepository: donationRepository,
		mealRepository:     mealRepository,
		userRepository:     userRepository,
	}
}

func (service *DonationService) CreateDonation(donationRequest *models.DonationRequest) error {
	var donation repository.Donation

	donor, err := service.userRepository.GetUserByName(donationRequest.DonorName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		donor = &repository.User{Name: donationRequest.DonorName}
		err = service.userRepository.CreateUser(donor)
	}

	if err != nil {
		return err
	}

	donation.DonorID = donor.ID
	donation.MealID = donationRequest.MealID

	err = service.donationRepository.CreateDonation(&donation)

	return err
}

func (service *DonationService) ClaimDonation(donationClaim *models.RecipientRequest) error {
	user, err := service.userRepository.GetUserByName(donationClaim.Name)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		user = &repository.User{Name: donationClaim.Name}
		err = service.userRepository.CreateUser(user)
	}

	if err != nil {
		return err
	}

	success, err := service.donationRepository.ClaimDonation(donationClaim.DonationID, user)
	if err != nil {
		return err
	}

	if !success {
		return ErrDonationNotFound
	}

	return nil
}

func (service *DonationService) GetUnclaimedDonationsByDate(today string) ([]models.UnclaimedDonationResponse, error) {
	var results []models.UnclaimedDonationResponse

	unclaimedDonations, err := service.donationRepository.GetUnclaimedDonationsByDate(today)

	for _, donation := range unclaimedDonations {
		results = append(results, models.UnclaimedDonationResponse{
			ID:          donation.ID,
			Description: donation.Meal.Description,
			DonorName:   donation.Donor.Name,
		})
	}

	return results, err
}

func (service *DonationService) GetDonationsSummaryByDate(date string) ([]models.DonationClaimSummaryResponse, error) {
	var donationClaimSummaries []models.DonationClaimSummaryResponse

	donations, err := service.donationRepository.GetDonationsSummaryByDate(date)

	if err != nil {
		return nil, err
	}

	for _, donation := range *donations {
		donationClaimSummaries = append(donationClaimSummaries, models.DonationClaimSummaryResponse{
			Claimed:       donation.Recipient.ID != 0,
			Description:   donation.Meal.Description,
			DonorName:     donation.Donor.Name,
			RecipientName: donation.Recipient.Name,
		})
	}

	return donationClaimSummaries, nil
}

func (service *DonationService) GetDonationClaimByClaimantName(name string) (models.ClaimedDonationResponse, error) {
	donation, err := service.donationRepository.GetDonationClaimByClaimantName(name)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.ClaimedDonationResponse{}, err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return models.ClaimedDonationResponse{}, nil
	}

	return models.ClaimedDonationResponse{
		UnclaimedDonationResponse: models.UnclaimedDonationResponse{
			ID:          donation.ID,
			DonorName:   donation.Donor.Name,
			Description: donation.Meal.Description,
		},
	}, nil
}

