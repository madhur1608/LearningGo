package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/madhur1608/aditis-kitchen/menu-service/controllers"
)

// Helper for making test requests
func makeCustomerRequest(t *testing.T, method, path string, body interface{}, handler http.HandlerFunc) (*httptest.ResponseRecorder, string) {
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
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr, rr.Body.String()
}

func TestSubscribeToPlan(t *testing.T) {
	// db.InitCouchbase()

	start := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	input := map[string]string{
		"mobile":     "8446485995",
		"plan_id":    "weekly_meal_plan",
		"start_date": start,
		"meal_type":  "lunch", // Required
	}
	rr, body := makeCustomerRequest(t, "POST", "/subscribe", input, controllers.SubscribeToPlan)

	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated {
		t.Fatalf("❌ TestSubscribeToPlan failed: %d - %s", rr.Code, body)
	}
	t.Logf("✅ TestSubscribeToPlan passed: %s", body)
}

func TestGetCurrentSubscription(t *testing.T) {
	// db.InitCouchbase()

	rr, body := makeCustomerRequest(t, "GET", "/subscription?mobile=8446485995", nil, controllers.GetCurrentSubscription)

	if rr.Code != http.StatusOK {
		t.Fatalf("❌ TestGetCurrentSubscription failed: %d - %s", rr.Code, body)
	}
	t.Logf("✅ TestGetCurrentSubscription passed: %s", body)
}

func TestPreventPastDateSubscription(t *testing.T) {
	// db.InitCouchbase()

	past := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	input := map[string]string{
		"mobile":     "9999999999",
		"plan_id":    "weekly_plan_lunch",
		"start_date": past,
		"meal_type":  "lunch", // Required
	}
	rr, body := makeCustomerRequest(t, "POST", "/subscribe", input, controllers.SubscribeToPlan)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("❌ TestPreventPastDateSubscription failed (expected 400): %d - %s", rr.Code, body)
	}
	t.Logf("✅ TestPreventPastDateSubscription passed")
}

func TestExtendActiveSubscription(t *testing.T) {
	// db.InitCouchbase()

	start := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	input := map[string]string{
		"mobile":     "8446485995",
		"plan_id":    "weekly_meal_plan",
		"start_date": start,
		"meal_type":  "lunch",
	}
	rr, body := makeCustomerRequest(t, "POST", "/subscribe", input, controllers.SubscribeToPlan)

	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated {
		t.Fatalf("❌ TestExtendActiveSubscription failed: %d - %s", rr.Code, body)
	}
	t.Logf("✅ TestExtendActiveSubscription passed: %s", body)
}
