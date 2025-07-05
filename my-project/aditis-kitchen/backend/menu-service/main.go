package main

import (
	"log"
	"net/http"

	"github.com/madhur1608/aditis-kitchen/menu-service/db"
	"github.com/madhur1608/aditis-kitchen/menu-service/routes"
)

func main() {
	db.InitCouchbase()
	router := routes.RegisterRoutes()
	log.Println("🚀 Menu-Service is running on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
