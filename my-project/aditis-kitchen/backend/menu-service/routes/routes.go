package routes

import (
	"github.com/gorilla/mux"
	"github.com/madhur1608/aditis-kitchen/menu-service/controllers"
	"github.com/madhur1608/aditis-kitchen/menu-service/middleware"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Health
	router.HandleFunc("/health", controllers.HealthCheck).Methods("GET")

	// Public Meal APIs
	router.HandleFunc("/meals/today", controllers.GetAvailableMeals).Methods("GET")
	router.HandleFunc("/meal-plans", controllers.GetMealPlans).Methods("GET")

	// Subscription APIs
	router.HandleFunc("/subscribe", controllers.SubscribeToPlan).Methods("POST")
	router.HandleFunc("/subscription", controllers.GetCurrentSubscription).Methods("GET")

	// Apply middleware for admin routes
	admin := router.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.RequireAdminAuth)
	admin.HandleFunc("/meals", controllers.CreateMeal).Methods("POST")
	admin.HandleFunc("/plans", controllers.CreateMealPlan).Methods("POST")
	admin.HandleFunc("/subscriptions", controllers.GetAllSubscriptions).Methods("GET")
	admin.HandleFunc("/meals/update", controllers.UpdateMeal).Methods("PUT")
	admin.HandleFunc("/plans/update", controllers.UpdateMealPlan).Methods("PUT")

	return router
}