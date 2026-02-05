package queries

import _ "embed"

// Meal
//go:embed meal/create_meal.sql
var CreateMeal string

//go:embed meal/get_meal_by_desc_date.sql
var GetMealByDescDate string

//go:embed meal/get_meals_by_date.sql
var GetMealsByDate string

//go:embed meal/get_meals_by_range.sql
var GetMealsByRange string

// User
//go:embed user/create_user.sql
var CreateUser string

//go:embed user/get_user_by_name.sql
var GetUserByName string

//go:embed user/upsert_user_google.sql
var UpsertUserGoogle string

//go:embed user/insert_user_google.sql
var InsertUserGoogle string

//go:embed user/update_user_google.sql
var UpdateUserGoogle string

//go:embed user/get_user_by_google_id.sql
var GetUserByGoogleID string

//go:embed user/get_user_by_id.sql
var GetUserByID string

//go:embed user/get_user_by_email.sql
var GetUserByEmail string

// Donation
//go:embed donation/create_donation.sql
var CreateDonation string

//go:embed donation/claim_donation.sql
var ClaimDonation string

//go:embed donation/get_unclaimed_donations.sql
var GetUnclaimedDonations string

//go:embed donation/get_donations_summary.sql
var GetDonationsSummary string

//go:embed donation/get_donation_claim_by_name.sql
var GetDonationClaimByName string

// Donation Request
//go:embed donation_request/create_donation_request.sql
var CreateDonationRequest string

//go:embed donation_request/create_donation_request_meal.sql
var CreateDonationRequestMeal string

//go:embed donation_request/get_requests_by_status.sql
var GetRequestsByStatus string

//go:embed donation_request/update_request_status.sql
var UpdateRequestStatus string

//go:embed donation_request/get_requests_by_requester.sql
var GetRequestsByRequester string

//go:embed donation_request/get_request_meals.sql
var GetRequestMeals string

// Meal Order
//go:embed meal_order/create_meal_order.sql
var CreateMealOrder string

//go:embed meal_order/get_meal_order_by_user_and_week.sql
var GetMealOrderByUserAndWeek string

//go:embed meal_order/get_meal_order_by_id.sql
var GetMealOrderByID string
