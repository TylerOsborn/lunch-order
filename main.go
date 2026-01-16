package main

import (
	_ "embed"
	"fmt"
	"log"
	"lunchorder/handlers"
	"lunchorder/queries"
	"lunchorder/repository"
	"lunchorder/router"
	"lunchorder/service"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	err = loadEnvironmentVariables(err)

	db, err := getDBConfig()
	if err != nil {
		log.Fatal(err)
	}
	initDB(db)

	// Repositories
	mealRepository := repository.NewMealRepository(db)
	userRepository := repository.NewUserRepository(db)
	donationRepository := repository.NewDonationRepository(db, userRepository)
	donationRequestRepository := repository.NewDonationRequestRepository(db, userRepository, donationRepository)

	// Services
	donationService := service.NewDonationService(donationRepository, mealRepository, userRepository)
	mealService := service.NewMealService(mealRepository)
	donationRequestService := service.NewDonationRequestService(donationRequestRepository, donationRepository, userRepository)

	// Handlers
	mealHandler := handlers.NewMealHandler(mealService)
	donationHandler := handlers.NewDonationHandler(donationService, donationRequestService)
	donationRequestHandler := handlers.NewDonationRequestHandler(donationRequestService)

	// Route setup
	r := gin.Default()
	router.SetupCors(r)
	router.SetupFrontEnd(r)
	router.SetupRoutes(r, mealHandler, donationHandler, donationRequestHandler)

	// Start server
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}

func loadEnvironmentVariables(err error) error {
	err = godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, reading from environment")
	}
	return err
}

func getDBConfig() (*sqlx.DB, error) {
	user, foundUser := os.LookupEnv("MYSQL_USER")
	password, foundPassword := os.LookupEnv("MYSQL_PASSWORD")
	host, foundHost := os.LookupEnv("MYSQL_HOST")
	port, foundPort := os.LookupEnv("MYSQL_PORT")
	database, foundDatabase := os.LookupEnv("MYSQL_DATABASE")

	if !foundUser || !foundPassword || !foundHost || !foundPort || !foundDatabase {
		return nil, fmt.Errorf("missing required environment variables")
	}

	finalString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	db, err := sqlx.Connect("mysql", finalString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func initDB(db *sqlx.DB) {
	_, err := db.Exec(queries.Schema)
	if err != nil {
		log.Fatal("Failed to execute schema migration: ", err)
	}
}
