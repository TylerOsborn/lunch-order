package main

import (
	"encoding/csv"
	"fmt"
	"lunchorder/constants"
	"lunchorder/models"
	"lunchorder/repository"
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
	r.GET("/Api/Meal", getMeals)

	r.POST("/Api/Meal/Upload", uploadWeeklyMeal)
	r.GET("/Api/Meal/Today", getTodayMeal)

	r.POST("/Api/Donation", donateMeal)
	r.GET("/Api/Donation", getDonations)

	r.POST("/Api/Donation/Claim", claimDonation)

	r.GET("/Api/Stats/Claims/Summary", getClaimsSummaryToday)
}

func getClaimsSummaryToday(context *gin.Context) {
	date := context.Query("date")

	if date == "" {
		date = time.Now().Format(constants.DATE_FORMAT)
	}

	var donationClaimSummaries []models.DonationClaimSummary
	tx := db.Raw("SELECT donations.recipient_id > 0 AS claimed, meals.description AS description, donors.name AS donor_name, COALESCE(recipients.name, 'UNCLAIMED') AS recipient_name FROM donations INNER JOIN users donors ON donors.id = donations.donor_id LEFT JOIN users recipients ON recipients.id = donations.recipient_id INNER JOIN meals ON donations.meal_id = meals.id WHERE meals.date = ?", date).Scan(&donationClaimSummaries)

	if tx.Error != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      tx.Error.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       donationClaimSummaries,
	})

}

func claimDonation(context *gin.Context) {
	var donationClaim models.APIRecipient
	err := context.BindJSON(&donationClaim)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	user, err := userRepository.GetUserByName(donationClaim.Name)

	if err != nil && err != gorm.ErrRecordNotFound {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	if err == gorm.ErrRecordNotFound {
		user = &models.User{Name: donationClaim.Name}
		userRepository.CreateUser(user)
	}

	success, err := donationRepository.ClaimDonation(donationClaim.DonationID, user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	if !success {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "Donation not found",
		})
		return
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func getDonations(context *gin.Context) {
	var donations []models.UnclaimedDonation

	today := time.Now().Format(constants.DATE_FORMAT)
	
	donations, err := donationRepository.GetUnclaimedDonationsByDate(today)
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

func donateMeal(context *gin.Context) {
	var donationRequest models.APIDonation
	err := context.BindJSON(&donationRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = donationRepository.CreateDonation(&donationRequest)
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

func getTodayMeal(context *gin.Context) {
	var meals []models.Meal
	today := time.Now().Format(constants.DATE_FORMAT)

	meals, err := mealRepository.GetMealsByDate(today)
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

type MealUpload struct {
	Csv string `json:"csv"`
}

func uploadWeeklyMeal(context *gin.Context) {
	var mealUpload MealUpload
	err := context.BindJSON(&mealUpload)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	csvString := mealUpload.Csv
	records, err := parseCSV(csvString)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// Print parsed records
	for _, record := range records {
		if len(record) != 2 {
			context.JSON(http.StatusBadRequest, models.ApiResult{
				StatusCode: http.StatusBadRequest,
				Error:      "CSV must have 2 columns",
			})
			return
		}
	}

	for _, record := range records {
		date, description := record[0], record[1]
		err := mealRepository.CreateMeal(&models.Meal{Date: date, Description: description})
		if err != nil {
			context.JSON(http.StatusInternalServerError, models.ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
			return
		}
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func parseCSV(input string) ([][]string, error) {
	r := csv.NewReader(strings.NewReader(input))

	r.FieldsPerRecord = 2

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV: %v", err)
	}

	return records, nil
}

func getMeals(context *gin.Context) {
	startDate := context.Query("startDate")
	endDate := context.Query("endDate")

	if startDate == "" || endDate == "" {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "startDate and endDate are required query parameters",
		})
		return
	}

	meals, err := mealRepository.GetMealsByDates(startDate, endDate)
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