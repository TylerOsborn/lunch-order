package router

import (
	"lunchorder/handlers"
	"lunchorder/repository"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, mealHandler *handlers.MealHandler, donationHandler *handlers.DonationHandler, donationRequestHandler *handlers.DonationRequestHandler, authHandler *handlers.AuthHandler, mealOrderHandler *handlers.MealOrderHandler, userRepo *repository.UserRepository) {
	// Auth routes
	r.GET("/auth/google/login", authHandler.GoogleLogin)
	r.GET("/auth/google/callback", authHandler.GoogleCallback)
	r.POST("/auth/logout", authHandler.Logout)

	// Protected routes
	api := r.Group("/Api")
	api.Use(handlers.AuthMiddleware(userRepo))
	{
		api.GET("/Me", authHandler.GetMe)

		api.GET("/Meal", mealHandler.HandleGetMeals)
		api.GET("/Meal/Today", mealHandler.HandleGetMealsToday)

		api.POST("/Donation", donationHandler.HandleDonateMeal)
		api.GET("/Donation", donationHandler.HandleGetUnclaimedDonations)

		api.POST("/Donation/Claim", donationHandler.HandleDonationClaim)
		api.GET("/Donation/Claim", donationHandler.HandleGetDonationClaim)

		// Donation request routes
		api.POST("/DonationRequest", donationRequestHandler.HandleCreateDonationRequest)
		api.GET("/DonationRequest", donationRequestHandler.HandleGetPendingDonationRequests)
		api.GET("/DonationRequest/User", donationRequestHandler.HandleGetUserDonationRequests)

		// Meal order routes
		api.POST("/MealOrder", mealOrderHandler.HandleCreateMealOrder)
		api.GET("/MealOrder", mealOrderHandler.HandleGetMealOrder)

		// Admin routes
		admin := api.Group("/")
		admin.Use(handlers.AdminMiddleware())
		{
			admin.POST("/Meal/Upload", mealHandler.HandleMealUpload)
			admin.GET("/Stats/Claims/Summary", donationHandler.HandleGetDonationSummary)
		}
	}
}

func SetupCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
}

func SetupFrontEnd(r *gin.Engine) {
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
