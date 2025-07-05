package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/madhur1608/aditis-kitchen/customer-service/db"
	"github.com/madhur1608/aditis-kitchen/customer-service/models"
	"github.com/madhur1608/aditis-kitchen/customer-service/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterCustomer(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "❌ Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validation
	customer.Name = strings.TrimSpace(customer.Name)
	customer.Mobile = strings.TrimSpace(customer.Mobile)

	if customer.Name == "" || customer.Mobile == "" || len(customer.Mobile) != 10 {
		http.Error(w, "❌ Name and valid 10-digit mobile number required", http.StatusBadRequest)
		return
	}

	if customer.Password == "" {
		http.Error(w, "❌ Password is required", http.StatusBadRequest)
		return
	}

	// Check for duplicate mobile
	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get("Customer::"+customer.Mobile, nil)
	if err == nil {
		http.Error(w, "❌ Customer with this mobile already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "❌ Error securing password", http.StatusInternalServerError)
		return
	}
	customer.Password = string(hashedPassword)
	customer.ID = customer.Mobile

	// Store in Couchbase
	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert("Customer::"+customer.ID, customer, nil)
	if err != nil {
		http.Error(w, "❌ Failed to register customer", http.StatusInternalServerError)
		return
	}

	// ✅ Send mock OTP after registration
	otp := utils.GenerateOTP(customer.Mobile)
	log.Printf("📨 OTP for %s (post-signup): %s", customer.Mobile, otp)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Customer registered successfully. OTP sent.",
	})
}

func LoginCustomer(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()
	var input struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "❌ Invalid input", http.StatusBadRequest)
		return
	}

	input.Mobile = strings.TrimSpace(input.Mobile)
	input.Password = strings.TrimSpace(input.Password)

	if input.Mobile == "" || input.Password == "" {
		http.Error(w, "❌ Mobile and password are required", http.StatusBadRequest)
		return
	}

	docID := "Customer::"+ input.Mobile
	result, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
	if err != nil {
		http.Error(w, "❌ Invalid mobile or customer not found", http.StatusUnauthorized)
		return
	}

	var customer models.Customer
	if err := result.Content(&customer); err != nil {
		http.Error(w, "❌ Failed to parse customer data", http.StatusInternalServerError)
		return
	}

	if !customer.IsVerified {
		http.Error(w, "❌ Please verify your OTP before login", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.Password))
	if err != nil {
		http.Error(w, "❌ Incorrect password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(input.Mobile)
	if err != nil {
		http.Error(w, "❌ Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Login successful",
		"token":   token,
		"name":    customer.Name,
	})
}

func GetCustomerByMobile(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()
	mobile := r.URL.Query().Get("mobile")
	if mobile == "" {
		http.Error(w, "❌ Mobile number is required", http.StatusBadRequest)
		return
	}

	docID := "Customer::"+ mobile
	result, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
	if err != nil {
		http.Error(w, "❌ Customer not found", http.StatusNotFound)
		return
	}

	var customer models.Customer
	if err := result.Content(&customer); err != nil {
		http.Error(w, "❌ Failed to parse customer data", http.StatusInternalServerError)
		return
	}

	customer.Password = "" // don't send hashed password
	json.NewEncoder(w).Encode(customer)
}

func UpdateCustomerProfile(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()

	mobile := r.URL.Query().Get("mobile")
	if mobile == "" {
		http.Error(w, "❌ Mobile number is required to update profile", http.StatusBadRequest)
		return
	}

	docID := "Customer::"+ mobile
	getResult, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
	if err != nil {
		http.Error(w, "❌ Customer not found", http.StatusNotFound)
		return
	}

	var existingCustomer models.Customer
	if err := getResult.Content(&existingCustomer); err != nil {
		http.Error(w, "❌ Failed to fetch existing customer data", http.StatusInternalServerError)
		return
	}

	var updateData struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "❌ Invalid request payload", http.StatusBadRequest)
		return
	}

	// Apply updates
	if updateData.Name != "" {
		existingCustomer.Name = strings.TrimSpace(updateData.Name)
	}
	if updateData.Email != "" {
		existingCustomer.Email = strings.TrimSpace(updateData.Email)
	}
	if updateData.Address != "" {
		existingCustomer.Address = strings.TrimSpace(updateData.Address)
	}

	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Replace(docID, existingCustomer, nil)
	if err != nil {
		http.Error(w, "❌ Failed to update customer profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "✅ Customer profile updated successfully"})
}

func ChangeCustomerPassword(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()

	var input struct {
		Mobile       string `json:"mobile"`
		OldPassword  string `json:"old_password"`
		NewPassword  string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Mobile == "" || input.OldPassword == "" || input.NewPassword == "" {
		http.Error(w, "❌ Invalid input. All fields are required.", http.StatusBadRequest)
		return
	}

	docID := "Customer::"+ strings.TrimSpace(input.Mobile)
	result, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
	if err != nil {
		http.Error(w, "❌ Customer not found", http.StatusNotFound)
		return
	}

	var customer models.Customer
	if err := result.Content(&customer); err != nil {
		http.Error(w, "❌ Failed to parse customer data", http.StatusInternalServerError)
		return
	}

	// Check old password
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(input.OldPassword)); err != nil {
		http.Error(w, "❌ Old password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "❌ Failed to hash new password", http.StatusInternalServerError)
		return
	}
	customer.Password = string(hashedNewPassword)

	// Update in DB
	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert(docID, customer, nil)
	if err != nil {
		http.Error(w, "❌ Failed to update password", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Password updated successfully",
	})
}

func ForgotCustomerPassword(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()

	var input struct {
		Mobile      string `json:"mobile"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Mobile == "" || input.NewPassword == "" {
		http.Error(w, "❌ Mobile and new password are required", http.StatusBadRequest)
		return
	}

	docID := "Customer::"+ strings.TrimSpace(input.Mobile)
	result, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
	if err != nil {
		http.Error(w, "❌ Customer not found", http.StatusNotFound)
		return
	}

	var customer models.Customer
	if err := result.Content(&customer); err != nil {
		http.Error(w, "❌ Failed to read customer data", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "❌ Failed to hash new password", http.StatusInternalServerError)
		return
	}

	customer.Password = string(hashedPassword)

	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert(docID, customer, nil)
	if err != nil {
		http.Error(w, "❌ Failed to update password", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Password reset successful",
	})
}

func RequestOTP(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()
	var input struct {
		Mobile string `json:"mobile"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Mobile == "" {
		http.Error(w, "❌ Mobile number required", http.StatusBadRequest)
		return
	}

	otp := utils.GenerateOTP(input.Mobile)

	// TODO: Send via SMS/Email (mocking now)
	log.Printf("📨 OTP for %s: %s", input.Mobile, otp)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ OTP sent",
	})
}

func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Mobile string `json:"mobile"`
		OTP    string `json:"otp"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Mobile == "" || input.OTP == "" {
		http.Error(w, "❌ Mobile and OTP required", http.StatusBadRequest)
		return
	}

	if utils.VerifyOTP(input.Mobile, input.OTP) {
		docID := "Customer::"+ input.Mobile
		result, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
		if err != nil {
			http.Error(w, "❌ Customer not found", http.StatusNotFound)
			return
		}
	
		var customer models.Customer
		if err := result.Content(&customer); err != nil {
			http.Error(w, "❌ Failed to read customer", http.StatusInternalServerError)
			return
		}
	
		token := uuid.New().String()
		utils.SetResetToken(input.Mobile, token)
	
		if !customer.IsVerified {
			customer.IsVerified = true
			_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Replace(docID, customer, nil)
			if err != nil {
				http.Error(w, "❌ Failed to verify customer", http.StatusInternalServerError)
				return
			}
	
			json.NewEncoder(w).Encode(map[string]string{
				"message": "✅ Customer verified successfully",
				"token":   token,
			})
		} else {
			json.NewEncoder(w).Encode(map[string]string{
				"message": "✅ OTP verified (already verified customer)",
				"token":   token,
			})
		}
	}
	
}

func ResetPasswordWithToken(w http.ResponseWriter, r *http.Request) {
	utils.EnsureDBConnection()
	var input struct {
		Mobile      string `json:"mobile"`
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil || input.Mobile == "" || input.Token == "" || input.NewPassword == "" {
		http.Error(w, "❌ All fields required", http.StatusBadRequest)
		return
	}

	if !utils.ValidateResetToken(input.Mobile, input.Token) {
		http.Error(w, "❌ Invalid or expired reset token", http.StatusUnauthorized)
		return
	}

	// Fetch and update customer
	docID := "Customer::"+ input.Mobile
	result, err := db.Cluster.Bucket("aditiskitchen").DefaultCollection().Get(docID, nil)
	if err != nil {
		http.Error(w, "❌ Customer not found", http.StatusNotFound)
		return
	}

	var customer models.Customer
	if err := result.Content(&customer); err != nil {
		http.Error(w, "❌ Parse error", http.StatusInternalServerError)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	customer.Password = string(hashedPassword)

	_, err = db.Cluster.Bucket("aditiskitchen").DefaultCollection().Upsert(docID, customer, nil)
	if err != nil {
		http.Error(w, "❌ Failed to update password", http.StatusInternalServerError)
		return
	}

	utils.ClearResetToken(input.Mobile)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Password reset successful",
	})
}