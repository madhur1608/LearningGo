package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/madhur1608/aditis-kitchen/customer-service/db"
	"github.com/madhur1608/aditis-kitchen/customer-service/routes"
)

func main() {
	db.InitCouchbase() // Initialize Couchbase connection

	router := mux.NewRouter()
	routes.RegisterCustomerRoutes(router)

	// CORS Middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	log.Println("✅ Customer service running on port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
