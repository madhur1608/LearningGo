package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/madhur1608/aditis-kitchen/menu-service/controllers"
	"github.com/madhur1608/aditis-kitchen/menu-service/db"
	"github.com/madhur1608/aditis-kitchen/menu-service/middleware"
)

func TestMain(m *testing.M) {
	err := db.InitCouchbase()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Couchbase: %v", err)
	}
	os.Exit(m.Run())
}

func makeAdminRequest(t *testing.T, method, path string, body interface{}, token string, handler http.HandlerFunc) (*httptest.ResponseRecorder, string) {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			t.Fatalf("Failed to encode body: %v", err)
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()

	// Wrap handler with admin middleware
	middleware.RequireAdminAuth(handler).ServeHTTP(rr, req)
	return rr, rr.Body.String()
}

func TestAdminE2EFlow(t *testing.T) {
	// Generate token
	token, err := middleware.GenerateToken("aditiskitchen17@gmail.com")
	if err != nil {
		t.Fatalf("❌ Failed to generate admin token: %v", err)
	}

	// 1. Create a meal
	createMeal := map[string]interface{}{
		"name":        "Full Meal",
		"description": "4 Fulka, Sabji, Daal, Rice, Salad",
		"price":       110,
		"days":        []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
	}
	rr, body := makeAdminRequest(t, "POST", "/admin/meals", createMeal, token, controllers.CreateMeal)
	if rr.Code != http.StatusCreated {
		t.Fatalf("❌ CreateMeal failed: %d - %s", rr.Code, body)
	}

	// // 2. Update the meal
	// updateMeal := map[string]interface{}{
	// 	"id":          "test_meal1",
	// 	"description": "2 Fulka, Sabji, Rice, Salad",
	// 	"price":       80,
	// 	"days":        []string{"Monday", "Tuesday"},
	// }
	// rr, body = makeAdminRequest(t, "PUT", "/admin/meals/test_meal", updateMeal, token, controllers.UpdateMeal)
	// if rr.Code != http.StatusOK {
	// 	t.Fatalf("❌ UpdateMeal failed: %d - %s", rr.Code, body)
	// }

	// 3. Create a plan
	createPlan := map[string]interface{}{
		"name":        "Weekly Meal Plan",
		"plan_price":       700,
		"total_meals": 7,
	}
	rr, body = makeAdminRequest(t, "POST", "/admin/plans", createPlan, token, controllers.CreateMealPlan)
	if rr.Code != http.StatusCreated {
		t.Fatalf("❌ CreateMealPlan failed: %d - %s", rr.Code, body)
	}

	// // 4. Update the plan
	// updatePlan := map[string]interface{}{
	// 	"id":       "test_plan1",
	// 	"price":    550,
	// 	"duration": 6,
	// }
	// rr, body = makeAdminRequest(t, "PUT", "/admin/plans/test_plan", updatePlan, token, controllers.UpdateMealPlan)
	// if rr.Code != http.StatusOK {
	// 	t.Fatalf("❌ UpdateMealPlan failed: %d - %s", rr.Code, body)
	// }

	// 5. Get all subscriptions
	rr, body = makeAdminRequest(t, "GET", "/admin/subscriptions", nil, token, controllers.GetAllSubscriptions)
	if rr.Code != http.StatusOK {
		t.Fatalf("❌ GetAllSubscriptions failed: %d - %s", rr.Code, body)
	}
}