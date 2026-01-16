package service

import (
	"database/sql"
	"errors"
	"lunchorder/models"
	"lunchorder/repository"
)

type DonationRequestService struct {
	donationRequestRepository *repository.DonationRequestRepository
	donationRepository        *repository.DonationRepository
	userRepository            *repository.UserRepository
}

var donationRequestService *DonationRequestService

func NewDonationRequestService(
	donationRequestRepository *repository.DonationRequestRepository,
	donationRepository *repository.DonationRepository,
	userRepository *repository.UserRepository) *DonationRequestService {

	return &DonationRequestService{
		donationRequestRepository: donationRequestRepository,
		donationRepository:        donationRepository,
		userRepository:            userRepository,
	}
}

func (s *DonationRequestService) CreateDonationRequest(request *models.DonationRequestCreate) error {
	// Get or create user
	user, err := s.userRepository.GetUserByName(request.RequesterName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		user = &repository.User{Name: request.RequesterName}
		if err := s.userRepository.CreateUser(user); err != nil {
			return err
		}
	}

	// Create the donation request with meal preferences
	return s.donationRequestRepository.CreateDonationRequest(user.ID, request.MealIds)
}

func (s *DonationRequestService) GetDonationRequestsByStatus(status string) ([]models.DonationRequestResponse, error) {
	requests, err := s.donationRequestRepository.GetDonationRequestsByStatus(status)
	if err != nil {
		return nil, err
	}

	var response []models.DonationRequestResponse
	for _, request := range requests {
		// Get meal preferences for this request
		meals, err := s.donationRequestRepository.GetDonationRequestMealPreferences(request.ID)
		if err != nil {
			continue
		}

		// Format the meal descriptions
		var description string
		for i, meal := range meals {
			if i > 0 {
				description += ", "
			}
			description += meal.Description
		}

		response = append(response, models.DonationRequestResponse{
			ID:            request.ID,
			RequesterName: request.Requester.Name,
			Description:   description,
			Status:        request.Status,
		})
	}

	return response, nil
}

func (s *DonationRequestService) GetDonationRequestsByRequesterName(name string, date string) ([]models.DonationRequestResponse, error) {
	user, err := s.userRepository.GetUserByName(name)
	if err != nil {
		return nil, err
	}

	requests, err := s.donationRequestRepository.GetDonationRequestsByRequesterName(name, date)
	if err != nil {
		return nil, err
	}

	var response []models.DonationRequestResponse
	for _, request := range requests {
		meals, err := s.donationRequestRepository.GetDonationRequestMealPreferences(request.ID)
		if err != nil {
			continue
		}

		var description string
		for i, meal := range meals {
			if i > 0 {
				description += ", "
			}
			description += meal.Description
		}

		donationResponse := models.DonationRequestResponse{
			ID:            request.ID,
			RequesterName: user.Name,
			Description:   description,
			Status:        request.Status,
		}

		response = append(response, donationResponse)
	}

	if response == nil {
		response = []models.DonationRequestResponse{}
	}

	return response, nil
}

func (s *DonationRequestService) UpdateDonationRequestStatus(id uint, status string) error {
	return s.donationRequestRepository.UpdateDonationRequestStatus(id, status, nil)
}

func (s *DonationRequestService) CheckAndFulfillDonationRequests() error {
	return s.donationRequestRepository.CheckAndFulfillDonationRequests()
}

