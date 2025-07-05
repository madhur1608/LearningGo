package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/madhur1608/aditis-kitchen/customer-service/controllers"
	"github.com/madhur1608/aditis-kitchen/customer-service/db"
	"github.com/madhur1608/aditis-kitchen/customer-service/models"
	"github.com/madhur1608/aditis-kitchen/customer-service/routes"
	"github.com/madhur1608/aditis-kitchen/customer-service/utils"
)

func TestMain(m *testing.M) {
	utils.EnsureDBConnection()
	os.Exit(m.Run())
}


func setupRouter() *mux.Router {
	utils.EnsureDBConnection()
	router := mux.NewRouter()
	routes.RegisterCustomerRoutes(router)
	return router
}


func TestRoutesExist(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		method string
		path   string
		body   *bytes.Buffer
	}{
		{
			method: "POST",
			path:   "/api/customer",
			body: bytes.NewBufferString(`{
				"name": "Test",
				"mobile": "9999999999",
				"password": "testpass"
			}`),
		},
		{
			method: "POST",
			path:   "/api/customer/login",
			body: bytes.NewBufferString(`{
				"mobile": "9999999999",
				"password": "testpass"
			}`),
		},
		{
			method: "GET",
			path:   "/api/customer?mobile=9999999999",
			body:   nil,
		},
		{
			method: "GET",
			path:   "/health",
			body:   nil,
		},
	}

	for _, tt := range tests {
		var bodyReader *bytes.Buffer
		if tt.body != nil {
			bodyReader = tt.body
		} else {
			bodyReader = bytes.NewBuffer([]byte{}) // safe fallback
		}

		req, err := http.NewRequest(tt.method, tt.path, bodyReader)
		if err != nil {
			t.Errorf("Could not create request for %s %s: %v", tt.method, tt.path, err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code == http.StatusNotFound {
			t.Errorf("Route %s %s not found (404)", tt.method, tt.path)
		}
	}
}


func TestHealthCheckRoute(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.HealthCheck)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Health check route returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateCustomerProfile(t *testing.T) {
	// Insert a test customer first
	testCustomer := models.Customer{
		ID:       "customer_9876543210",
		Name:     "Test User",
		Mobile:   "9876543210",
		Email:    "",
		Address:  "",
		Password: "$2a$10$B5N4VDehCbGWyzSxdeLce.b7VIA98E08ZPL8UEuigH9QYEJrEhAjK", // "password"
	}

	_, err := db.Cluster.Bucket("customers").DefaultCollection().Upsert(testCustomer.ID, testCustomer, nil)
	if err != nil {
		t.Fatalf("Failed to insert test customer: %v", err)
	}

	payload := `{"name": "New Name2", "email": "new2@example.com", "address": "New Address2"}`
	req, err := http.NewRequest("PUT", "/api/customer/profile?mobile=9876543210", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(controllers.UpdateCustomerProfile)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if res["message"] != "✅ Customer profile updated successfully" {
		t.Errorf("Unexpected response message: %v", res["message"])
	}
}
