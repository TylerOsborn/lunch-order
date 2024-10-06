package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"lunchorder/constants"
	"lunchorder/models"
	"lunchorder/repository"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var donationRepository *repository.DonationRepository
var mealRepository *repository.MealRepository
var donationClaimRepository *repository.DonationClaimRepository

func main() {
	// DB setup
	fmt.Println("Setting up db...")
	var err error
	db, err = sql.Open("sqlite3", "./database/database.db")
	if err != nil {
		log.Fatal(err)
	}
	initDB()
	defer db.Close()

	donationRepository = repository.NewDonationRepository(db)
	mealRepository = repository.NewMealRepository(db)
	donationClaimRepository = repository.NewDonationClaimRepository(db)

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

func initDB() {

	sqlStmt := `
CREATE TABLE IF NOT EXISTS  meal 
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS donation
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        VARCHAR(255),
    description TEXT,
    claimed     BOOLEAN DEFAULT FALSE,
    date        DATE    DEFAULT CURRENT_DATE,
    doe         DATETIME   DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS donation_claim
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donation_id INT,
    name VARCHAR(255),
    doe DATETIME DEFAULT CURRENT_TIMESTAMP
)
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
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

	rows, err := db.Query("SELECT donation.claimed AS claimed, donation.description AS decription, donation.name AS donator_name, COALESCE(donation_claim.name, 'UNCLAIMED') AS claimer_name FROM donation LEFT JOIN donation_claim ON donation.id = donation_claim.donation_id WHERE donation.date = ?", date)

	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	defer rows.Close()

	for rows.Next() {
		var donationClaimSummary models.DonationClaimSummary
		err := rows.Scan(&donationClaimSummary.Claimed, &donationClaimSummary.Description, &donationClaimSummary.DonatorName, &donationClaimSummary.ClaimerName)
		if err != nil {
			context.JSON(http.StatusInternalServerError, models.ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}

		if donationClaimSummary.ClaimerName != "UNCLAIMED" {
			donationClaimSummary.Claimed = true
		}

		donationClaimSummaries = append(donationClaimSummaries, donationClaimSummary)
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
		Data:       donationClaimSummaries,
	})

}

func claimDonation(context *gin.Context) {
	var donationClaim models.DonationClaim
	err := context.BindJSON(&donationClaim)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	success, err := donationRepository.ClaimDonation(donationClaim.DonationId)
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

	err = donationClaimRepository.CreateDonationClaim(&donationClaim)

	if err != nil {
		context.JSON(http.StatusOK, models.ApiResult{
			StatusCode: http.StatusOK,
			Error:      "Meal was allocated but the following error was produced: " + err.Error(),
		})
	}

	context.JSON(http.StatusOK, models.ApiResult{
		StatusCode: http.StatusOK,
	})
}

func getDonations(context *gin.Context) {
	var donations []models.Donation

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
	var donation models.Donation
	err := context.BindJSON(&donation)
	if err != nil {
		context.JSON(http.StatusBadRequest, models.ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	err = donationRepository.CreateDonation(&donation)
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

func getTodayMeal(context *gin.Context) {
	var meals []models.Meal
	today := time.Now().Format(constants.DATE_FORMAT)

	meals, err := mealRepository.GetMealsByDate(today)
	if err != nil {
		context.JSON(http.StatusInternalServerError, models.ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
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
