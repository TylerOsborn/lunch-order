package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"lunchorder/constants"
	"lunchorder/middleware"
	"lunchorder/models"
	"lunchorder/repository"
	"lunchorder/service"
	"lunchorder/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var donationRepository *repository.DonationRepository
var mealRepository *repository.MealRepository
var userRepository *repository.UserRepository
var donationRequestRepository *repository.DonationRequestRepository

var donationService *service.DonationService
var mealService *service.MealService
var donationRequestService *service.DonationRequestService
var authService *service.AuthService

var authMiddleware *middleware.AuthMiddleware

func main() {
	var err error

	err = loadEnvironmentVariables(err)

	err = getDBConfig()
	if err != nil {
		return
	}
	initDB(db)

	mealRepository = repository.NewMealRepository(db)
	userRepository = repository.NewUserRepository(db)
	donationRepository = repository.NewDonationRepository(db, userRepository)
	donationRequestRepository = repository.NewDonationRequestRepository(db, userRepository, donationRepository)

	donationService = service.NewDonationService(donationRepository, mealRepository, userRepository)
	mealService = service.NewMealService(mealRepository)
	donationRequestService = service.NewDonationRequestService(donationRequestRepository, donationRepository, userRepository)
	
	// Initialize auth service
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	if googleRedirectURL == "" {
		googleRedirectURL = "http://localhost:8080/Api/Auth/Callback"
	}
	authService = service.NewAuthService(userRepository, googleClientID, googleClientSecret, googleRedirectURL)
	
	// Initialize auth middleware
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		sessionSecret = "your-session-secret-change-in-production"
	}
	authMiddleware = middleware.NewAuthMiddleware(sessionSecret, userRepository)

	// Route setup
	r := gin.Default()
	setupCors(r)
	setupFrontEnd(r)
	setupRoutes(r)

	// Start server
	err = r.Run(":8080")
	if err != nil {
		return
	}
}

func loadEnvironmentVariables(err error) error {
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, reading from environment")
	}
	return err
}

func getDBConfig() error {
	var err error

	user, foundUser := os.LookupEnv("MYSQL_USER")
	password, foundPassword := os.LookupEnv("MYSQL_PASSWORD")
	host, foundHost := os.LookupEnv("MYSQL_HOST")
	port, foundPort := os.LookupEnv("MYSQL_PORT")
	database, foundDatabase := os.LookupEnv("MYSQL_DATABASE")

	if !foundUser || !foundPassword || !foundHost || !foundPort || !foundDatabase {
		panic("Missing required environment variables")
	}

	finalString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	db, err = gorm.Open(mysql.Open(finalString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return err
}

func setupCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func initDB(db *gorm.DB) {
	err := db.Migrator().AutoMigrate(&repository.Meal{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().AutoMigrate(&repository.User{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().AutoMigrate(&repository.Donation{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().AutoMigrate(&repository.DonationRequest{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().AutoMigrate(&repository.DonationRequestMeal{})
	if err != nil {
		panic(err)
	}
}

func setupFrontEnd(r *gin.Engine) {
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/vite.svg", "./frontend/dist/vite.svg")
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api") && !strings.Contains(path, ".") {
			c.File("./frontend/dist/index.html")
		} else {
			c.Next()
		}
	})
}

func setupRoutes(r *gin.Engine) {
	// Auth routes (public)
	r.GET("/Api/Auth/Login", HandleGoogleLogin)
	r.GET("/Api/Auth/Callback", HandleGoogleCallback)
	r.POST("/Api/Auth/Logout", HandleLogout)
	r.GET("/Api/Auth/Profile", authMiddleware.RequireAuth(), HandleGetProfile)

	// Public routes
	r.GET("/Api/Meal", HandleGetMeals)
	r.GET("/Api/Meal/Today", HandleGetMealsToday)
	r.GET("/Api/Donation", HandleGetUnclaimedDonations)
	r.GET("/Api/DonationRequest", HandleGetPendingDonationRequests)

	// Protected routes (require authentication)
	authenticated := r.Group("/Api")
	authenticated.Use(authMiddleware.RequireAuth())
	{
		authenticated.POST("/Donation", HandleDonateMeal)
		authenticated.POST("/Donation/Claim", HandleDonationClaim)
		authenticated.GET("/Donation/Claim", HandleGetDonationClaim)
		authenticated.POST("/DonationRequest", HandleCreateDonationRequest)
		authenticated.GET("/DonationRequest/User", HandleGetUserDonationRequests)
	}

	// Admin-only routes
	admin := r.Group("/Api")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireAdmin())
	{
		admin.POST("/Meal/Upload", HandleMealUpload)
		admin.GET("/Stats/Claims/Summary", HandleGetDonationSummary)
	}
}

func HandleGetDonationSummary(context *gin.Context) {
	date := context.Query("date")

	if date == "" {
		date = time.Now().Format(constants.DATE_FORMAT)
	}

	donationClaimSummaries, err := donationService.GetDonationsSummaryByDate(date)

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

func HandleDonationClaim(context *gin.Context) {
	var donationClaim models.RecipientRequest
	err := context.BindJSON(&donationClaim)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = donationService.ClaimDonation(&donationClaim)

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

func HandleGetDonationClaim(context *gin.Context) {
	var claimantName string
	claimantName = context.Query("name")

	if claimantName == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "name is a required query parameter",
		})
		return
	}

	claimed, err := donationService.GetDonationClaimByClaimantName(claimantName)

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

func HandleGetUnclaimedDonations(context *gin.Context) {
	var donations []models.UnclaimedDonationResponse

	today := time.Now().Format(constants.DATE_FORMAT)

	donations, err := donationService.GetUnclaimedDonationsByDate(today)

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

func HandleDonateMeal(context *gin.Context) {
	var donationRequest models.DonationRequest
	err := context.BindJSON(&donationRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = donationService.CreateDonation(&donationRequest)

	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if donationRequestService != nil {
		_ = donationRequestService.CheckAndFulfillDonationRequests()
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func HandleGetMealsToday(context *gin.Context) {
	var meals []models.MealResponse
	today := time.Now().Format(constants.DATE_FORMAT)

	meals, err := mealService.GetMealsByDate(today)

	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       meals,
	})
}

func HandleMealUpload(context *gin.Context) {
	var mealUpload models.MealUploadRequest
	err := context.BindJSON(&mealUpload)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = mealService.CreateMeals(mealUpload)

	if errors.Is(err, utils.ErrIncorrectCSVFormat) {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func HandleGetMeals(context *gin.Context) {
	startDate := context.Query("startDate")
	endDate := context.Query("endDate")

	if startDate == "" || endDate == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "startDate and endDate are required query parameters",
		})
		return
	}

	meals, err := mealService.GetMealsByDates(startDate, endDate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       meals,
	})
}

// HandleCreateDonationRequest creates a new donation request
func HandleCreateDonationRequest(context *gin.Context) {
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

	err = donationRequestService.CreateDonationRequest(&donationRequestData)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// After creating the request, check if any requests can be fulfilled with available donations
	donationRequestService.CheckAndFulfillDonationRequests()

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func HandleGetPendingDonationRequests(context *gin.Context) {
	requests, err := donationRequestService.GetDonationRequestsByStatus("pending")
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

func HandleGetUserDonationRequests(context *gin.Context) {
	userName := context.Query("name")
	if userName == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "name is a required query parameter",
		})
		return
	}

	date := context.Query("date")

	requests, err := donationRequestService.GetDonationRequestsByRequesterName(userName, date)
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
// Auth handlers
func HandleGoogleLogin(context *gin.Context) {
	state := "random-state-string" // In production, use a proper random state
	authURL := authService.GetAuthURL(state)
	context.JSON(http.StatusOK, gin.H{"auth_url": authURL})
}

func HandleGoogleCallback(context *gin.Context) {
	code := context.Query("code")
	if code == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
		return
	}

	user, err := authService.HandleCallback(code)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = authMiddleware.SetUserSession(context, user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set session"})
		return
	}

	// Redirect to frontend
	context.Redirect(http.StatusFound, "http://localhost:5173/")
}

func HandleLogout(context *gin.Context) {
	err := authMiddleware.ClearUserSession(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func HandleGetProfile(context *gin.Context) {
	user, exists := context.Get("user")
	if !exists {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	userData := user.(*repository.User)
	context.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    userData.ID,
			"name":  userData.Name,
			"email": userData.Email,
			"role":  userData.Role,
		},
	})
}
