package handlers

import (
	"lunchorder/constants"
	"lunchorder/models"
	"lunchorder/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DonationHandler struct {
	donationService        *service.DonationService
	donationRequestService *service.DonationRequestService
}

func NewDonationHandler(donationService *service.DonationService, donationRequestService *service.DonationRequestService) *DonationHandler {
	return &DonationHandler{
		donationService:        donationService,
		donationRequestService: donationRequestService,
	}
}

func (h *DonationHandler) HandleDonateMeal(context *gin.Context) {
	var donationRequest models.DonationRequest
	err := context.BindJSON(&donationRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = h.donationService.CreateDonation(&donationRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if h.donationRequestService != nil {
		_ = h.donationRequestService.CheckAndFulfillDonationRequests()
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func (h *DonationHandler) HandleGetUnclaimedDonations(context *gin.Context) {
	today := time.Now().Format(constants.DateFormat)

	donations, err := h.donationService.GetUnclaimedDonationsByDate(today)

	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       donations,
	})
}

func (h *DonationHandler) HandleDonationClaim(context *gin.Context) {
	var donationClaim models.RecipientRequest
	err := context.BindJSON(&donationClaim)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = h.donationService.ClaimDonation(&donationClaim)

	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func (h *DonationHandler) HandleGetDonationClaim(context *gin.Context) {
	claimantName := context.Query("name")

	if claimantName == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "name is a required query parameter",
		})
		return
	}

	claimed, err := h.donationService.GetDonationClaimByClaimantName(claimantName)

	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	if claimed.ID <= 0 {
		context.JSON(http.StatusNotFound, models.ApiResult{
			StatusCode: http.StatusNotFound,
			Error:      "No claimed donations found",
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       claimed,
	})
}

func (h *DonationHandler) HandleGetDonationSummary(context *gin.Context) {
	date := context.Query("date")

	if date == "" {
		date = time.Now().Format(constants.DateFormat)
	}

	donationClaimSummaries, err := h.donationService.GetDonationsSummaryByDate(date)

	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       donationClaimSummaries,
	})
}
