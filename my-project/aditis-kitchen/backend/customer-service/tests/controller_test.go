package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/madhur1608/aditis-kitchen/customer-service/controllers"
	"github.com/madhur1608/aditis-kitchen/customer-service/utils"
)

func makeJSONRequest(t *testing.T, method, path string, body any, handler http.HandlerFunc) (*httptest.ResponseRecorder, string) {
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

func TestCustomerFullE2EFlow(t *testing.T) {
	utils.EnsureDBConnection()

	mobile := "8446485995"
	password := "Aditi@17"
	// newPassword := "Aditi@16"

	// Step 1: Register Customer
	rr, body := makeJSONRequest(t, "POST", "/register", map[string]string{
		"name":     "Madhur Bajaj",
		"mobile":   mobile,
		"password": password,
	}, controllers.RegisterCustomer)
	if rr.Code != http.StatusCreated {
		t.Fatalf("❌ Registration failed: %d - %s", rr.Code, body)
	}

	// Step 2: Verify OTP sent after signup
	otp := utils.MockOTPs[mobile]
	rr, body = makeJSONRequest(t, "POST", "/verify-otp", map[string]string{
		"mobile": mobile,
		"otp":    otp.OTP,
	}, controllers.VerifyOTP)
	if rr.Code != http.StatusOK {
		t.Fatalf("❌ OTP verification failed: %d - %s", rr.Code, body)
	}

	// Step 3: Login with original password
	rr, body = makeJSONRequest(t, "POST", "/login", map[string]string{
		"mobile":   mobile,
		"password": password,
	}, controllers.LoginCustomer)
	if rr.Code != http.StatusOK {
		t.Fatalf("❌ Login after verification failed: %d - %s", rr.Code, body)
	}

// 	// Step 4: Change password (old -> new)
// 	rr, body = makeJSONRequest(t, "POST", "/change-password", map[string]string{
// 		"mobile":       mobile,
// 		"old_password": password,
// 		"new_password": newPassword,
// 	}, controllers.ChangeCustomerPassword)
// 	if rr.Code != http.StatusOK {
// 		t.Fatalf("❌ Password change failed: %d - %s", rr.Code, body)
// 	}

// 	// Step 5: Forgot password (overwrite back to original password)
// 	rr, body = makeJSONRequest(t, "POST", "/forgot-password", map[string]string{
// 		"mobile":      mobile,
// 		"new_password": password,
// 	}, controllers.ForgotCustomerPassword)
// 	if rr.Code != http.StatusOK {
// 		t.Fatalf("❌ Forgot password failed: %d - %s", rr.Code, body)
// 	}

// 	// Step 6: Request OTP for password reset flow
// 	rr, body = makeJSONRequest(t, "POST", "/request-otp", map[string]string{
// 		"mobile": mobile,
// 	}, controllers.RequestOTP)
// 	if rr.Code != http.StatusOK {
// 		t.Fatalf("❌ OTP request for reset failed: %d - %s", rr.Code, body)
// 	}

// 	// Step 7: Verify OTP again to get reset token
// 	otp = utils.MockOTPs[mobile]
// 	rr, body = makeJSONRequest(t, "POST", "/verify-otp", map[string]string{
// 		"mobile": mobile,
// 		"otp":    otp.OTP,
// 	}, controllers.VerifyOTP)
// 	if rr.Code != http.StatusOK {
// 		t.Fatalf("❌ OTP verification for reset failed: %d - %s", rr.Code, body)
// 	}

// 	tokenData, ok := utils.MockResetTokens[mobile]
// 	if !ok || tokenData.Token == "" {
// 		t.Fatalf("❌ Reset token not returned from utils")
// 	}
// 	resetToken := tokenData.Token

// 	// Step 8: Use reset token to set a new password
// 	rr, body = makeJSONRequest(t, "POST", "/reset-password", map[string]string{
// 		"mobile":       mobile,
// 		"token":        resetToken,
// 		"new_password": newPassword,
// 	}, controllers.ResetPasswordWithToken)
// 	if rr.Code != http.StatusOK {
// 		t.Fatalf("❌ Reset password with token failed: %d - %s", rr.Code, body)
// 	}

// 	// Step 9: Final login with newly set password
// 	rr, body = makeJSONRequest(t, "POST", "/login", map[string]string{
// 		"mobile":   mobile,
// 		"password": newPassword,
// 	}, controllers.LoginCustomer)
// 	if rr.Code != http.StatusOK {
// 		t.Fatalf("❌ Final login with new password failed: %d - %s", rr.Code, body)
// 	}
}
