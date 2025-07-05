package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/madhur1608/aditis-kitchen/menu-service/db"
	"github.com/madhur1608/aditis-kitchen/menu-service/models"
)

// CreateMeal handles creation of a new meal item
func CreateMeal(w http.ResponseWriter, r *http.Request) {
	var meal models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&meal); err != nil {
		http.Error(w, "❌ Invalid input", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(meal.Name) == "" || meal.Price <= 0 || len(meal.Days) == 0 {
		http.Error(w, "❌ Name, price, and days are required", http.StatusBadRequest)
		return
	}	

	meal.ID = strings.ToLower(strings.ReplaceAll(meal.Name, " ", "_"))
	meal.Type = "meal"

	_, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert("Meal::"+meal.ID, meal, nil)
	if err != nil {
		http.Error(w, "❌ Could not store meal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Meal created successfully",
		"meal_id": meal.ID,
	})
}

// CreateMealPlan handles creation of a new meal plan
func CreateMealPlan(w http.ResponseWriter, r *http.Request) {
	var plan models.MealPlan
	if err := json.NewDecoder(r.Body).Decode(&plan); err != nil {
		http.Error(w, "❌ Invalid input", http.StatusBadRequest)
		return
	}
	plan.ID = strings.ToLower(strings.ReplaceAll(plan.Name, " ", "_"))
	plan.Type = "plan"

	_, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert("Plan::"+plan.ID, plan, nil)
	if err != nil {
		http.Error(w, "❌ Could not store meal plan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Meal plan created successfully",
		"plan_id": plan.ID,
	})
}

// GetAllSubscriptions returns all subscriptions for review
func GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	query := `SELECT s.* FROM aditiskitchen s WHERE s.type = "subscription";`
	rows, err := db.Cluster.Query(query, nil)
	if err != nil {
		http.Error(w, "❌ Failed to query subscriptions"+err.Error(), http.StatusInternalServerError)
		return
	}

	var subscriptions []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Row(&sub); err == nil {
			subscriptions = append(subscriptions, sub)
		}
	}

	json.NewEncoder(w).Encode(subscriptions)
}

// UpdateMeal allows admin to update an existing meal
func UpdateMeal(w http.ResponseWriter, r *http.Request) {
	var updated models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil || strings.TrimSpace(updated.ID) == "" {
		http.Error(w, "❌ Invalid input or missing meal ID", http.StatusBadRequest)
		return
	}

	updated.Type = "meal"

	_, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Replace(updated.ID, updated, nil)
	if err != nil {
		http.Error(w, "❌ Failed to update meal"+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Meal updated successfully",
	})
}

// UpdateMealPlan allows admin to update an existing meal plan
func UpdateMealPlan(w http.ResponseWriter, r *http.Request) {
	var updated models.MealPlan
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil || strings.TrimSpace(updated.ID) == "" {
		http.Error(w, "❌ Invalid input or missing plan ID", http.StatusBadRequest)
		return
	}

	updated.Type = "plan"

	_, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Replace(updated.ID, updated, nil)
	if err != nil {
		http.Error(w, "❌ Failed to update meal plan"+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Meal plan updated successfully",
	})
}
