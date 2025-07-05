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

func getISTNow() time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	return time.Now().In(loc)
}

func SubscribeToPlan(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Mobile    string `json:"mobile"`
		PlanID    string `json:"plan_id"`
		StartDate string `json:"start_date"` // Format: YYYY-MM-DD
		MealType  string `json:"meal_type"`  // lunch or dinner
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Mobile == "" || input.PlanID == "" || input.StartDate == "" || input.MealType == "" {
		http.Error(w, "❌ Invalid input", http.StatusBadRequest)
		return
	}

	input.MealType = strings.ToLower(input.MealType)
	if input.MealType != "lunch" && input.MealType != "dinner" {
		http.Error(w, "❌ Meal type must be lunch or dinner", http.StatusBadRequest)
		return
	}

	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := getISTNow()

	// Parse and validate start date
	startDateParsed, err := time.ParseInLocation("2006-01-02", input.StartDate, loc)
	if err != nil {
		http.Error(w, "❌ Invalid start date format", http.StatusBadRequest)
		return
	}
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	if !startDateParsed.After(today) {
		http.Error(w, "❌ Start date must be from tomorrow onward", http.StatusBadRequest)
		return
	}

	// Set midnight timestamp to startDate
	startDate := time.Date(startDateParsed.Year(), startDateParsed.Month(), startDateParsed.Day(), 0, 0, 0, 0, loc)

	// Fetch plan
	res, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get("Plan::"+input.PlanID, nil)
	if err != nil {
		http.Error(w, "❌ Invalid plan selected", http.StatusNotFound)
		return
	}

	var plan models.MealPlan
	if err := res.Content(&plan); err != nil {
		http.Error(w, "❌ Failed to parse plan", http.StatusInternalServerError)
		return
	}

	customerID := strings.TrimSpace(input.Mobile)
	subscriptionID := customerID + "_" + input.MealType

	// Check for existing active subscription
	query := `
		SELECT META(s).id, s.* FROM aditiskitchen s 
		WHERE s.type = "subscription" AND s.customer_id = $1 AND s.meal_type = $2
		ORDER BY s.created_at DESC LIMIT 1;
	`
	rows, err := db.Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID, input.MealType},
	})
	if err != nil {
		http.Error(w, "❌ Failed to query existing subscriptions", http.StatusInternalServerError)
		return
	}

	if rows.Next() {
		var existing models.Subscription
		var meta struct {
			ID string `json:"id"`
		}
		if err := rows.Row(&meta); err == nil {
			_ = rows.Row(&existing)
			endDate, _ := time.Parse("2006-01-02", existing.EndDate)
			if endDate.After(today) || endDate.Equal(today) {
				// Extend existing
				existing.EndDate = endDate.AddDate(0, 0, plan.TotalMeals).Format("2006-01-02")
				existing.TotalMeals += plan.TotalMeals
				existing.UpdatedAt = now.Format(time.RFC3339)

				_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Replace(meta.ID, existing, nil)
				if err != nil {
					http.Error(w, "❌ Failed to extend subscription", http.StatusInternalServerError)
					return
				}
				json.NewEncoder(w).Encode(map[string]string{
					"message": "✅ Subscription extended successfully",
					"sub_id":  meta.ID,
				})
				return
			}
		}
	}

	// New subscription
	sub := models.Subscription{
		ID:         subscriptionID,
		CustomerID: customerID,
		PlanID:     plan.ID,
		PlanName:   plan.Name,
		StartDate:  startDate.Format("2006-01-02"),
		EndDate:    startDate.AddDate(0, 0, plan.TotalMeals).Format("2006-01-02"),
		TotalMeals: plan.TotalMeals,
		PlanPrice:  plan.PlanPrice,
		Status:     "active",
		MealType:   input.MealType,
		Type:       "subscription",
		CreatedAt:  now.Format(time.RFC3339),
	}

	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert("Subscription::" + sub.PlanID +"_"+ sub.ID , sub, nil)
	if err != nil {
		http.Error(w, "❌ Failed to store subscription", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Subscribed successfully",
		"sub_id":  sub.ID,
	})
}

func GetCurrentSubscription(w http.ResponseWriter, r *http.Request) {
	mobile := r.URL.Query().Get("mobile")
	mealType := r.URL.Query().Get("meal_type")
	if mobile == "" || mealType == "" {
		http.Error(w, "❌ Mobile and meal_type are required", http.StatusBadRequest)
		return
	}
	customerID := "customer_" + strings.TrimSpace(mobile)

	query := `
		SELECT s.* FROM aditiskitchen s 
		WHERE s.type = "subscription" AND s.customer_id = $1 AND s.meal_type = $2
		ORDER BY s.created_at DESC 
		LIMIT 1;
	`
	rows, err := db.Cluster.Query(query, &gocb.QueryOptions{
		PositionalParameters: []interface{}{customerID, mealType},
	})
	if err != nil {
		http.Error(w, "❌ Failed to fetch subscription", http.StatusInternalServerError)
		return
	}

	if rows.Next() {
		var sub models.Subscription
		err := rows.Row(&sub)
		if err != nil {
			http.Error(w, "❌ Failed to parse subscription", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sub)
		return
	}

	http.Error(w, "❌ No subscription found", http.StatusNotFound)
}
