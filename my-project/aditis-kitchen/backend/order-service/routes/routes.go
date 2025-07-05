package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madhur1608/aditis-kitchen/order-service/controllers"
)

func Registerorder_^serviceRoutes(router *mux.Router) {
	router.HandleFunc("/health", controllers.HealthCheck).Methods("GET")
	// Add other routes here
}
