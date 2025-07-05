package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/madhur1608/aditis-kitchen/menu-service/db"
	"github.com/madhur1608/aditis-kitchen/menu-service/models"
)

// GetAvailableMeals returns meals available for today (Mon–Fri vs Sat–Sun)
func GetAvailableMeals(w http.ResponseWriter, r *http.Request) {
	today := strings.Title(time.Now().Weekday().String()) // "Monday", etc.

	query := "SELECT m.* FROM `aditiskitchen` m WHERE $1 IN m.days AND m.type = 'meal';"
	rows, err := db.Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{today},
	})
	if err != nil {
		http.Error(w, "❌ Failed to fetch menu", http.StatusInternalServerError)
		return
	}

	var meals []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Row(&item); err == nil {
			meals = append(meals, item)
		}
	}

	json.NewEncoder(w).Encode(meals)
}

// GetMealPlans returns all available meal plans
func GetMealPlans(w http.ResponseWriter, r *http.Request) {
	query := "SELECT p.* FROM `aditiskitchen` p WHERE p.type = 'plan';"
	rows, err := db.Cluster.Query(query, &gocb.QueryOptions{})
	if err != nil {
		http.Error(w, "❌ Failed to fetch plans", http.StatusInternalServerError)
		return
	}

	var plans []models.MealPlan
	for rows.Next() {
		var plan models.MealPlan
		if err := rows.Row(&plan); err == nil {
			plans = append(plans, plan)
		}
	}

	json.NewEncoder(w).Encode(plans)
}
