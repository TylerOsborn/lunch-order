package main

import (
	"errors"
	"lunchorder/constants"
	"lunchorder/models"
	"lunchorder/repository"
	"lunchorder/service"
	"lunchorder/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

var donationRepository *repository.DonationRepository
var mealRepository *repository.MealRepository
var userRepository *repository.UserRepository

var donationService *service.DonationService
var mealService *service.MealService

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("./database/database.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	initDB(db)

	mealRepository = repository.NewMealRepository(db)
	userRepository = repository.NewUserRepository(db)
	donationRepository = repository.NewDonationRepository(db, userRepository)

	donationService = service.NewDonationService(donationRepository, mealRepository, userRepository)
	mealService = service.NewMealService(mealRepository)

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
	err := db.Migrator().AutoMigrate(&models.Meal{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}

	err = db.Migrator().AutoMigrate(&models.Donation{})
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
	r.GET("/Api/Meal", HandleGetMeals)

	r.POST("/Api/Meal/Upload", HandleMealUpload)
	r.GET("/Api/Meal/Today", HandleGetMealsToday)

	r.POST("/Api/Donation", HandleDonateMeal)
	r.GET("/Api/Donation", HandleGetUnclaimedDonations)

	r.POST("/Api/Donation/Claim", HandleDonationClaim)
	r.GET("/Api/Donation/Claim", HandleGetDonationClaim)

	r.GET("/Api/Stats/Claims/Summary", HandleGetDonationSummary)
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
	var donationClaim models.APIRecipient
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
	var donations []models.UnclaimedDonation

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
	var donationRequest models.APIDonation
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

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func HandleGetMealsToday(context *gin.Context) {
	var meals []models.Meal
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
	var mealUpload models.MealUpload
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
