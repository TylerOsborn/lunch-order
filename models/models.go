package models

type DonationRequest struct {
	MealID    uint   `json:"mealId"`
	DonorName string `json:"donorName"`
}

type RecipientRequest struct {
	DonationID uint   `json:"donationId"`
	Name       string `json:"name"`
}

type UnclaimedDonationResponse struct {
	ID          uint   `json:"id"`
	DonorName   string `json:"donorName"`
	Description string `json:"description"`
}

type ClaimedDonationResponse struct {
	UnclaimedDonationResponse
}

type MealResponse struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type ApiResult struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}

type DonationClaimSummaryResponse struct {
	Claimed       bool   `json:"claimed"`
	Description   string `json:"description"`
	DonorName     string `json:"donorName"`
	RecipientName string `json:"recipientName"`
}

type MealUploadRequest struct {
	Csv string `json:"csv"`
}

type DonationRequestCreate struct {
	RequesterName string `json:"requesterName"`
	MealIds       []uint `json:"mealIds"`
}

type DonationRequestResponse struct {
	ID            uint   `json:"id"`
	RequesterName string `json:"requesterName"`
	Description   string `json:"description"`
	Status        string `json:"status"`
}

type MealOrderRequest struct {
	WeekStartDate   string `json:"weekStartDate"`
	MondayMealID    *uint  `json:"mondayMealId"`
	TuesdayMealID   *uint  `json:"tuesdayMealId"`
	WednesdayMealID *uint  `json:"wednesdayMealId"`
	ThursdayMealID  *uint  `json:"thursdayMealId"`
}

type MealOrderResponse struct {
	ID              uint          `json:"id"`
	WeekStartDate   string        `json:"weekStartDate"`
	MondayMeal      *MealResponse `json:"mondayMeal"`
	TuesdayMeal     *MealResponse `json:"tuesdayMeal"`
	WednesdayMeal   *MealResponse `json:"wednesdayMeal"`
	ThursdayMeal    *MealResponse `json:"thursdayMeal"`
}
