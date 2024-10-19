package service

import (
	"errors"
	"lunchorder/models"
	"lunchorder/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	if donationService == nil {
		donationService = &DonationService{
			donationRepository: donationRepository,
			mealRepository:     mealRepository,
			userRepository:     userRepository,
		}
	}

	return donationService
}

func (service *DonationService) CreateDonation(donationRequest *models.APIDonation) (*models.Donation, error) {
	var donation *models.Donation = &models.Donation{
	}

	donor, err := service.userRepository.GetUserByUUID(donationRequest.User.UUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if donationRequest.User.UUID == "" {
			donationRequest.User.UUID = uuid.NewString()
		}
		err = service.userRepository.CreateUser(&donationRequest.User)
		donor = &donationRequest.User
	}

	if err != nil {
		return nil, err
	}

	donation.DonorID = donor.ID
	donation.MealID = donationRequest.MealID

	err = service.donationRepository.CreateDonation(donation)

	if err != nil {
		return nil, err
	}

	return service.donationRepository.GetDonationByID(donation.ID)
}

func (service *DonationService) ClaimDonation(donationClaim *models.APIRecipient) (*models.Donation, error) {
	user, err := service.userRepository.GetUserByUUID(donationClaim.User.UUID)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		if donationClaim.User.UUID == "" {
			donationClaim.User.UUID = uuid.NewString()
		}
		err = service.userRepository.CreateUser(&donationClaim.User)
		user = &donationClaim.User
	}

	if err != nil {
		return nil, err
	}

	success, err := service.donationRepository.ClaimDonation(donationClaim.DonationID, user)
	if err != nil {
		return nil, err
	}

	if !success {
		return nil, ErrDonationNotFound
	}

	donation, err := service.donationRepository.GetDonationByID(donationClaim.DonationID)

	if err != nil {
		return nil, err
	}

	return donation, nil
}

func (service *DonationService) GetUnclaimedDonationsByDate(today string) ([]models.UnclaimedDonation, error) {
	return service.donationRepository.GetUnclaimedDonationsByDate(today)
}

func (service *DonationService) GetDonationsSummaryByDate(date string) ([]models.DonationClaimSummary, error) {
	var donationClaimSummaries []models.DonationClaimSummary

	err := service.donationRepository.GetDonationsSummaryByDate(date, &donationClaimSummaries)

	if err != nil {
		return nil, err
	}

	return donationClaimSummaries, nil
}

func (service *DonationService) GetDonationByRecipientUUIDAndDate(userUUID string, date string, donation *models.Donation) error {
	user, err := service.userRepository.GetUserByUUID(userUUID)

	if err != nil {
		return err
	}

	return service.donationRepository.GetDonationByRecipientIDAndDate(user.ID, date, donation)
}
