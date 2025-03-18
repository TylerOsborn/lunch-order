package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"lunchorder/constants"
	"lunchorder/models"
	"lunchorder/repository"
	"lunchorder/service"
	"lunchorder/utils"
)

var db *gorm.DB

var (
	donationRepository *repository.DonationRepository
	mealRepository     *repository.MealRepository
	userRepository     *repository.UserRepository
)

var (
	donationService *service.DonationService
	mealService     *service.MealService
)

func main() {
	var err error

	// Configure logging to stdout for fly.io
	log.SetOutput(os.Stdout)
	log.Println("Starting lunch-order application")

	loadEnvironmentVariables()
	log.Println("Environment variables loaded")

	err = getDBConfig()
	if err != nil {
		log.Fatal("Failed to configure database: ", err)
		return
	}
	log.Println("Database configuration successful")

	initDB(db)
	log.Println("Database initialized")

	mealRepository = repository.NewMealRepository(db)
	userRepository = repository.NewUserRepository(db)
	donationRepository = repository.NewDonationRepository(db, userRepository)
	log.Println("Repositories initialized")

	donationService = service.NewDonationService(donationRepository, mealRepository, userRepository)
	mealService = service.NewMealService(mealRepository)
	log.Println("Services initialized")

	// Route setup
	r := gin.Default()
	setupCors(r)
	setupFrontEnd(r)
	setupRoutes(r)
	log.Println("Routes configured")

	// Start server
	log.Println("Starting server on port 8080")
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Server failed to start: ", err)
		return
	}
}

func loadEnvironmentVariables() error {
	err := godotenv.Load()
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

	// Log database configuration (excluding password)
	log.Printf("Database configuration check - User: %v, Host: %v, Port: %v, Database: %v",
		foundUser, foundHost, foundPort, foundDatabase)

	if !foundUser || !foundPassword || !foundHost || !foundPort || !foundDatabase {
		log.Println("Missing required database environment variables")
		if !foundUser {
			log.Println("MYSQL_USER not set")
		}
		if !foundPassword {
			log.Println("MYSQL_PASSWORD not set")
		}
		if !foundHost {
			log.Println("MYSQL_HOST not set")
		}
		if !foundPort {
			log.Println("MYSQL_PORT not set")
		}
		if !foundDatabase {
			log.Println("MYSQL_DATABASE not set")
		}
		panic("Missing required environment variables")
	}

	finalString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	log.Printf("Connecting to MySQL at %s:%s/%s", host, port, database)

	db, err = gorm.Open(mysql.Open(finalString), &gorm.Config{})
	if err != nil {
		log.Printf("Database connection error: %v", err)
		panic(err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get DB instance: %v", err)
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return err
	}

	log.Println("Database connection successful")
	return nil
}

func setupCors(r *gin.Engine) {
	log.Println("Setting up CORS configuration")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	log.Println("CORS configuration complete")
}

func initDB(db *gorm.DB) {
	log.Println("Starting database migrations")

	log.Println("Migrating Meal table")
	err := db.Migrator().AutoMigrate(&repository.Meal{})
	if err != nil {
		log.Printf("Failed to migrate Meal table: %v", err)
		panic(err)
	}

	log.Println("Migrating User table")
	err = db.Migrator().AutoMigrate(&repository.User{})
	if err != nil {
		log.Printf("Failed to migrate User table: %v", err)
		panic(err)
	}

	log.Println("Migrating Donation table")
	err = db.Migrator().AutoMigrate(&repository.Donation{})
	if err != nil {
		log.Printf("Failed to migrate Donation table: %v", err)
		panic(err)
	}

	log.Println("All migrations completed successfully")
}

func setupFrontEnd(r *gin.Engine) {
	log.Println("Setting up frontend static file serving")

	// Check if frontend files exist
	_, err := os.Stat("./frontend/dist/assets")
	if err != nil {
		log.Printf("Warning: Frontend assets directory not found: %v", err)
	}

	_, err = os.Stat("./frontend/dist/index.html")
	if err != nil {
		log.Printf("Warning: Frontend index.html not found: %v", err)
	}

	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/vite.svg", "./frontend/dist/vite.svg")
	log.Println("Static file routes configured")

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		log.Printf("NoRoute handler for path: %s", path)
		if !strings.HasPrefix(path, "/api") && !strings.Contains(path, ".") {
			log.Printf("Serving index.html for path: %s", path)
			c.File("./frontend/dist/index.html")
		} else {
			c.Next()
		}
	})

	log.Println("Frontend setup complete")
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
