package main

import (
	"database/sql"
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
	fmt.Println("Starting server...")
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
CREATE TABLE IF NOT EXISTS meal_type (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS  meal (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type_id INTEGER NOT NULL,
    date DATE NOT NULL,
    FOREIGN KEY (type_id) REFERENCES meal_type (id)
);

CREATE TABLE IF NOT EXISTS donation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    meal_id INTEGER NOT NULL,
    donor_name TEXT NOT NULL,
    claimed BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (meal_id) REFERENCES meal (id)
);

CREATE TABLE IF NOT EXISTS claim (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    donation_id INTEGER,
    claimer_name TEXT NOT NULL,
    fulfilled BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (donation_id) REFERENCES donation (id)
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
	r.GET("/Api/MealType", getMealTypes)
	r.POST("/Api/MealType", createMealType)

	r.GET("/Api/Meal", getMeals)

	r.GET("/Api/Menu", getMenu)
	//r.POST("/Api/Meal", postMeal)
}

func getMeals(g *gin.Context) {
	startDate := g.Query("startDate")
	endDate := g.Query("endDate")

	if startDate == "" || endDate == "" {
		g.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "startDate and endDate are required query parameters",
		})
		return
	}

	var meals []Meal
	rows, err := db.Query("SELECT * FROM meal WHERE date >= ? AND date <= ?", startDate, endDate)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	defer rows.Close()

	for rows.Next() {
		var meal Meal
		err := rows.Scan(&meal.Id, &meal.TypeId, &meal.Date)
		if err != nil {
			log.Fatal(err)
		}

		date, err := time.Parse(time.RFC3339, meal.Date)
		if err != nil {
			g.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		meal.Date = date.Format("2006-01-02")
		meals = append(meals, meal)
	}

	g.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
		Data:       meals,
	})

}

func getMealTypes(g *gin.Context) {
	var mealTypes []MealType
	rows, err := db.Query("SELECT * FROM meal_type")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var mealType MealType
		err := rows.Scan(&mealType.Id, &mealType.Description)
		if err != nil {
			g.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
		}
		mealTypes = append(mealTypes, mealType)
	}

	g.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
		Data:       mealTypes,
	})

}

func createMealType(g *gin.Context) {
	MealType := MealType{}
	err := g.BindJSON(&MealType)
	if err != nil {
		g.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      err.Error(),
		})
		return
	}

	stmt, err := db.Prepare("INSERT INTO meal_type(description) VALUES(?)")
	if err != nil {
		g.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	result, err := stmt.Exec(MealType.Description)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
		return
	}

	id, _ := result.LastInsertId()
	MealType.Id = int(id)

	g.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
		Data:       MealType,
	})
}

func getMenu(g *gin.Context) {
	startDate := g.Query("startDate")

	if startDate == "" {
		g.JSON(http.StatusBadRequest, ApiResult{
			StatusCode: http.StatusBadRequest,
			Error:      "startDate is a required query parameter",
		})
		return
	}

	var meals []MenuMeal

	sqlQuery := "SELECT strftime('%Y-%m-%d', date) AS date, description FROM meal INNER JOIN meal_type ON meal.type_id = meal_type.id WHERE date >= ? AND date <= DATE(?, '+5 day')"

	rows, err := db.Query(sqlQuery, startDate, startDate)
	if err != nil {
		g.JSON(http.StatusInternalServerError, ApiResult{
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		})
	}

	defer rows.Close()

	for rows.Next() {
		var menuDay MenuMeal
		err := rows.Scan(&menuDay.Date, &menuDay.Description)
		if err != nil {
			g.JSON(http.StatusInternalServerError, ApiResult{
				StatusCode: http.StatusInternalServerError,
				Error:      err.Error(),
			})
			return
		}
		meals = append(meals, menuDay)
	}

	g.JSON(http.StatusOK, ApiResult{
		StatusCode: http.StatusOK,
		Data:       meals,
	})
}

type MealType struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type Meal struct {
	Id     int    `json:"id"`
	TypeId int    `json:"typeId"`
	Date   string `json:"date"`
}

type Donation struct {
	Id        int    `json:"id"`
	MealId    int    `json:"mealId"`
	DonorName string `json:"donorName"`
	Claimed   bool   `json:"claimed"`
}

type Claim struct {
	Id          int    `json:"id"`
	DonationId  int    `json:"donationId"`
	ClaimerName string `json:"claimerName"`
	Fulfilled   bool   `json:"fulfilled"`
}

type ApiResult struct {
	StatusCode int         `json:"statusCode"`
	Error      string      `json:"error"`
	Data       interface{} `json:"data"`
}

type MenuMeal struct {
	Date        string `json:"date"`
	Description string `json:"description"`
}
