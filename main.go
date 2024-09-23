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
	r.Run(":8080")
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
CREATE TABLE IF NOT EXISTS  meal (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    date DATE NOT NULL
);
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
	rows, err := db.Query("SELECT * FROM meal WHERE date >= ? AND date <= ?", startDate, endDate)
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
			log.Fatal(err)
		}

		date, err := time.Parse(time.RFC3339, meal.Date)
		if err != nil {
			context.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		meal.Date = date.Format("2006-01-02")
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

type ApiResult struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}
