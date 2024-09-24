package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
	"time"
)

const DATE_FORMAT = "2006-01-02"

var db *sql.DB

func main() {
	// DB setup
	fmt.Println("Setting up db...")
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	initDB()
	defer db.Close()

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
}

func claimDonation(context *gin.Context) {
	var donationClaim DonationClaim
	err := context.BindJSON(&donationClaim)
	if err != nil {
		context.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	result, err := db.Exec("UPDATE donation SET claimed = 1 WHERE id = ?", donationClaim.DonationId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		context.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "Donation not found",
		})
		return
	}

	_, err = db.Exec("INSERT INTO donation_claim (donation_id, name) VALUES (?, ?)", donationClaim.DonationId, donationClaim.Name)

	if err != nil {
		context.JSON(http.StatusOK, ApiResult{
			StatusCode: http.StatusOK,
			Error:      "Meal was allocated but the follwing error was produced: " + err.Error(),
		})
	}

	context.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
	})
}

func getDonations(context *gin.Context) {
	var donations []Donation

	today := time.Now().Format(DATE_FORMAT)
	rows, err := db.Query("SELECT DISTINCT description, name, id FROM donation WHERE claimed = 0 AND date = ? GROUP BY description ORDER BY description", today)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	defer rows.Close()

	for rows.Next() {
		var donation Donation
		err := rows.Scan(&donation.Description, &donation.Name, &donation.Id)
		if err != nil {
			context.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		donations = append(donations, donation)
	}

	context.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
		Data:       donations,
	})

}

func donateMeal(context *gin.Context) {
	var donation Donation
	err := context.BindJSON(&donation)
	if err != nil {
		context.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	_, err = db.Exec("INSERT INTO donation (name, description) VALUES (?, ?)", donation.Name, donation.Description)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	context.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
	})
}

func getTodayMeal(context *gin.Context) {
	var meals []Meal
	today := time.Now().Format(DATE_FORMAT)
	rows, err := db.Query("SELECT id, description, date FROM meal WHERE date = ?", today)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.Id, &meal.Description, &meal.Date)
		if err != nil {
			context.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		meal.Date = time.Now().Format(DATE_FORMAT)
		meals = append(meals, meal)
	}

	context.JSON(http.StatusOK, ApiResult{
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
		context.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	csvString := mealUpload.Csv
	records, err := parseCSV(csvString)
	if err != nil {
		context.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	// Print parsed records
	for _, record := range records {
		if len(record) != 2 {
			context.JSON(http.StatusBadRequest, ApiResult{
				StatusCode: http.StatusBadRequest,
				Error:      "CSV must have 2 columns",
			})
			return
		}
	}

	for _, record := range records {
		date, description := record[0], record[1]
		_, err := db.Exec("INSERT INTO meal (description, date) VALUES (?, ?)", description, date)
		if err != nil {
			context.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
	}

	context.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
	})
}

func parseCSV(input string) ([][]string, error) {
	// Create a new CSV reader
	r := csv.NewReader(strings.NewReader(input))

	// Set the field delimiter to comma
	r.Comma = ','

	// Allow variable number of fields per record
	r.FieldsPerRecord = -1

	// Parse all records
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
		context.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "startDate and endDate are required query parameters",
		})
		return
	}

	var meals []Meal
	rows, err := db.Query("SELECT id, description, date FROM meal WHERE date >= ? AND date <= ? ORDER BY description", startDate, endDate)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	defer rows.Close()

	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.Id, &meal.Description, &meal.Date)
		if err != nil {
			context.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}

		date, err := time.Parse(time.RFC3339, meal.Date)
		if err != nil {
			context.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		meal.Date = date.Format(DATE_FORMAT)
		meals = append(meals, meal)
	}

	context.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
		Data:       meals,
	})

}

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
