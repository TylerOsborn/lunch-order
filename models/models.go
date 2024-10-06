package models


type Meal struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type Donation struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Date        string `json:"date"`
	Doe         string `json:"doe"`
}

type DonationClaim struct {
	Id         int    `json:"id"`
	DonationId int    `json:"donationId"`
	Name       string `json:"name"`
	Doe        string `json:"doe"`
}

type ApiResult struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}

type DonationClaimSummary struct {
	Claimed     bool   `json:"claimed"`
	Description string `json:"description"`
	DonatorName string `json:"donatorName"`
	ClaimerName string `json:"claimerName"`
}