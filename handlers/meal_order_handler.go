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

	// Convert the denormalized request to normalized items
	var items []repository.MealOrderItem
	if request.MondayMealID != nil {
		items = append(items, repository.MealOrderItem{
			DayOfWeek: "Monday",
			MealID:    *request.MondayMealID,
		})
	}
	if request.TuesdayMealID != nil {
		items = append(items, repository.MealOrderItem{
			DayOfWeek: "Tuesday",
			MealID:    *request.TuesdayMealID,
		})
	}
	if request.WednesdayMealID != nil {
		items = append(items, repository.MealOrderItem{
			DayOfWeek: "Wednesday",
			MealID:    *request.WednesdayMealID,
		})
	}
	if request.ThursdayMealID != nil {
		items = append(items, repository.MealOrderItem{
			DayOfWeek: "Thursday",
			MealID:    *request.ThursdayMealID,
		})
	}

	// Create the meal order
	order := &repository.MealOrder{
		UserID:        userID,
		WeekStartDate: request.WeekStartDate,
		Items:         items,
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

	// Convert normalized structure to denormalized response format
	response := models.MealOrderResponse{
		ID:            order.ID,
		WeekStartDate: order.WeekStartDate,
	}

	// Map items to the denormalized response structure
	for _, item := range order.Items {
		if item.Meal != nil {
			mealResponse := &models.MealResponse{
				ID:          item.Meal.ID,
				Description: item.Meal.Description,
				Date:        item.Meal.Date,
			}

			switch item.DayOfWeek {
			case "Monday":
				response.MondayMeal = mealResponse
			case "Tuesday":
				response.TuesdayMeal = mealResponse
			case "Wednesday":
				response.WednesdayMeal = mealResponse
			case "Thursday":
				response.ThursdayMeal = mealResponse
			}
		}
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       response,
	})
}
