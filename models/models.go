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
