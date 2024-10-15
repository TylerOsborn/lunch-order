package service

import (
	"errors"
	"gorm.io/gorm"
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

	if donationService == nil {
		donationService = &DonationService{
			donationRepository: donationRepository,
			mealRepository:     mealRepository,
			userRepository:     userRepository,
		}
	}

	return donationService
}

func (service *DonationService) CreateDonation(donationRequest *models.APIDonation) error {
	var donation models.Donation

	donor, err := service.userRepository.GetUserByName(donationRequest.DonorName)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		donor = &models.User{Name: donationRequest.DonorName}
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

func (service *DonationService) ClaimDonation(donationClaim *models.APIRecipient) error {
	user, err := service.userRepository.GetUserByName(donationClaim.Name)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		user = &models.User{Name: donationClaim.Name}
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
