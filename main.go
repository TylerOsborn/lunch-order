package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"lunchorder/handlers"
	"lunchorder/repository"
	"lunchorder/router"
	"lunchorder/service"
	"os"
	"time"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

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
	authHandler := handlers.NewAuthHandler(userRepository)

	// Route setup
	r := gin.Default()
	router.SetupCors(r)
	router.SetupFrontEnd(r)
	router.SetupRoutes(r, mealHandler, donationHandler, donationRequestHandler, authHandler, userRepository)

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

	// Add connection parameters for reliability and timeouts
	finalString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true&timeout=10s&readTimeout=30s&writeTimeout=30s&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		user, password, host, port, database)

	db, err := sqlx.Connect("mysql", finalString)
	if err != nil {
		return nil, err
	}

	// Configure connection pool to prevent connection exhaustion and timeouts
	db.SetMaxOpenConns(25)                  // Maximum number of open connections to the database
	db.SetMaxIdleConns(5)                   // Maximum number of idle connections in the pool
	db.SetConnMaxLifetime(5 * time.Minute)  // Maximum lifetime of a connection (prevents stale connections)
	db.SetConnMaxIdleTime(1 * time.Minute)  // Maximum time a connection can be idle

	return db, nil
}

func initDB(db *sqlx.DB) {
	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver: ", err)
	}

	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatal("Failed to create migration source: ", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "mysql", driver)
	if err != nil {
		log.Fatal("Failed to create migration instance: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations: ", err)
	}
	log.Println("Migrations ran successfully")
}
