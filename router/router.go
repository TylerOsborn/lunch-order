package router

import (
	"lunchorder/handlers"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, mealHandler *handlers.MealHandler, donationHandler *handlers.DonationHandler, donationRequestHandler *handlers.DonationRequestHandler) {
	r.GET("/Api/Meal", mealHandler.HandleGetMeals)

	r.POST("/Api/Meal/Upload", mealHandler.HandleMealUpload)
	r.GET("/Api/Meal/Today", mealHandler.HandleGetMealsToday)

	r.POST("/Api/Donation", donationHandler.HandleDonateMeal)
	r.GET("/Api/Donation", donationHandler.HandleGetUnclaimedDonations)

	r.POST("/Api/Donation/Claim", donationHandler.HandleDonationClaim)
	r.GET("/Api/Donation/Claim", donationHandler.HandleGetDonationClaim)

	r.GET("/Api/Stats/Claims/Summary", donationHandler.HandleGetDonationSummary)

	// Donation request routes
	r.POST("/Api/DonationRequest", donationRequestHandler.HandleCreateDonationRequest)
	r.GET("/Api/DonationRequest", donationRequestHandler.HandleGetPendingDonationRequests)
	r.GET("/Api/DonationRequest/User", donationRequestHandler.HandleGetUserDonationRequests)
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
