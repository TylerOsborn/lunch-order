package main

import (
	"embed"
	"fmt"
	"log"
	"lunchorder/handlers"
	"lunchorder/repository"
	"lunchorder/router"
	"lunchorder/service"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
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

func init() {
	log.Println("Building frontend...")

	// Check if node_modules directory exists
	if _, err := os.Stat(filepath.Join("frontend", "node_modules")); os.IsNotExist(err) {
		log.Println("node_modules not found, running npm install...")
		cmd := exec.Command("npm", "install")
		cmd.Dir = "frontend"
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Failed to run npm install: %v\n%s", err, output)
		}
		log.Println("npm install completed successfully")
	} else {
		log.Println("node_modules found, skipping npm install")
	}

	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = "frontend"

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to build frontend: %v\n%s", err, output)
	}

	log.Println("Frontend built successfully")
	log.Println(string(output))
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

	finalString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, password, host, port, database)
	db, err := sqlx.Connect("mysql", finalString)
	if err != nil {
		return nil, err
	}
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
