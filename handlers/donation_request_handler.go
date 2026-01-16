package handlers

import (
	"lunchorder/models"
	"lunchorder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DonationRequestHandler struct {
	donationRequestService *service.DonationRequestService
}

func NewDonationRequestHandler(donationRequestService *service.DonationRequestService) *DonationRequestHandler {
	return &DonationRequestHandler{donationRequestService: donationRequestService}
}

func (h *DonationRequestHandler) HandleCreateDonationRequest(context *gin.Context) {
	var donationRequestData models.DonationRequestCreate
	err := context.BindJSON(&donationRequestData)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// Validate request data
	if donationRequestData.RequesterName == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "requesterName is required",
		})
		return
	}

	if len(donationRequestData.MealIds) == 0 {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "at least one meal must be selected",
		})
		return
	}

	err = h.donationRequestService.CreateDonationRequest(&donationRequestData)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// After creating the request, check if any requests can be fulfilled with available donations
	h.donationRequestService.CheckAndFulfillDonationRequests()

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func (h *DonationRequestHandler) HandleGetPendingDonationRequests(context *gin.Context) {
	requests, err := h.donationRequestService.GetDonationRequestsByStatus("pending")
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       requests,
	})
}

func (h *DonationRequestHandler) HandleGetUserDonationRequests(context *gin.Context) {
	userName := context.Query("name")
	if userName == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "name is a required query parameter",
		})
		return
	}

	date := context.Query("date")

	requests, err := h.donationRequestService.GetDonationRequestsByRequesterName(userName, date)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       requests,
	})
}
