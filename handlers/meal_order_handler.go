package handlers

import (
	"lunchorder/models"
	"lunchorder/repository"
	"lunchorder/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MealOrderHandler struct {
	mealOrderService *service.MealOrderService
	mealService      *service.MealService
}

func NewMealOrderHandler(mealOrderService *service.MealOrderService, mealService *service.MealService) *MealOrderHandler {
	return &MealOrderHandler{
		mealOrderService: mealOrderService,
		mealService:      mealService,
	}
}

func (h *MealOrderHandler) HandleCreateMealOrder(context *gin.Context) {
	var request models.MealOrderRequest
	err := context.BindJSON(&request)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// Get user ID from context (set by AuthMiddleware)
	userIDInterface, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusUnauthorized, models.ApiResult{
			StatusCode: http.StatusUnauthorized,
			Error:      "User not authenticated",
		})
		return
	}
	userID := userIDInterface.(uint)

	// Check if order already exists for this user and week
	existingOrder, err := h.mealOrderService.GetMealOrderByUserAndWeek(userID, request.WeekStartDate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	if existingOrder != nil {
		context.JSON(http.StatusConflict, models.ApiResult{
			StatusCode: http.StatusConflict,
			Error:      "Meal order already exists for this week",
		})
		return
	}

	// Create the meal order
	order := &repository.MealOrder{
		UserID:          userID,
		WeekStartDate:   request.WeekStartDate,
		MondayMealID:    request.MondayMealID,
		TuesdayMealID:   request.TuesdayMealID,
		WednesdayMealID: request.WednesdayMealID,
		ThursdayMealID:  request.ThursdayMealID,
	}

	err = h.mealOrderService.CreateMealOrder(order)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       "Meal order created successfully",
	})
}

func (h *MealOrderHandler) HandleGetMealOrder(context *gin.Context) {
	weekStartDate := context.Query("weekStartDate")
	if weekStartDate == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "weekStartDate is required query parameter",
		})
		return
	}

	// Get user ID from context (set by AuthMiddleware)
	userIDInterface, exists := context.Get("userID")
	if !exists {
		context.JSON(http.StatusUnauthorized, models.ApiResult{
			StatusCode: http.StatusUnauthorized,
			Error:      "User not authenticated",
		})
		return
	}
	userID := userIDInterface.(uint)

	// Get the meal order
	order, err := h.mealOrderService.GetMealOrderByUserAndWeek(userID, weekStartDate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	if order == nil {
		context.JSON(http.StatusOK, models.ApiResult{
			StatusCode: http.StatusOK,
			Data:       nil,
		})
		return
	}

	// Convert to response format
	response := models.MealOrderResponse{
		ID:            order.ID,
		WeekStartDate: order.WeekStartDate,
	}

	// Load meal details using GetMealByID for better performance
	if order.MondayMealID != nil {
		meal, err := h.mealService.GetMealByID(*order.MondayMealID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, models.ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      "Failed to load Monday meal details: " + err.Error(),
			})
			return
		}
		response.MondayMeal = meal
	}

	if order.TuesdayMealID != nil {
		meal, err := h.mealService.GetMealByID(*order.TuesdayMealID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, models.ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      "Failed to load Tuesday meal details: " + err.Error(),
			})
			return
		}
		response.TuesdayMeal = meal
	}

	if order.WednesdayMealID != nil {
		meal, err := h.mealService.GetMealByID(*order.WednesdayMealID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, models.ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      "Failed to load Wednesday meal details: " + err.Error(),
			})
			return
		}
		response.WednesdayMeal = meal
	}

	if order.ThursdayMealID != nil {
		meal, err := h.mealService.GetMealByID(*order.ThursdayMealID)
		if err != nil {
			context.JSON(http.StatusInternalServerError, models.ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      "Failed to load Thursday meal details: " + err.Error(),
			})
			return
		}
		response.ThursdayMeal = meal
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       response,
	})
}
