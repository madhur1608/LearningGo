package routes

import (
	"github.com/gorilla/mux"
	"github.com/madhur1608/aditis-kitchen/customer-service/controllers"
)

func RegisterCustomerRoutes(router *mux.Router) {
	router.HandleFunc("/api/customer", controllers.RegisterCustomer).Methods("POST")
	router.HandleFunc("/api/customer/login", controllers.LoginCustomer).Methods("POST")
	router.HandleFunc("/api/customer", controllers.GetCustomerByMobile).Methods("GET")
	router.HandleFunc("/health", controllers.HealthCheck).Methods("GET")
	router.HandleFunc("/api/customer/profile", controllers.UpdateCustomerProfile).Methods("PUT")
	router.HandleFunc("/api/customer/change-password", controllers.ChangeCustomerPassword).Methods("PUT")
	router.HandleFunc("/api/customer/forgot-password", controllers.ForgotCustomerPassword).Methods("PUT")
	router.HandleFunc("/api/customer/request-otp", controllers.RequestOTP).Methods("POST")
	router.HandleFunc("/api/customer/verify-otp", controllers.VerifyOTP).Methods("POST")
	router.HandleFunc("/api/customer/reset-password", controllers.ResetPasswordWithToken).Methods("PUT")
}
