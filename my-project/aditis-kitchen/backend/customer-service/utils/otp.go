package utils

import (
	"math/rand"
	"sync"
	"time"
)

type OTPData struct {
	OTP       string
	ExpiresAt time.Time
}

type TokenData struct {
	Token     string
	ExpiresAt time.Time
}

var MockOTPs = make(map[string]OTPData)
var MockResetTokens = make(map[string]TokenData)
var mu sync.Mutex

func GenerateOTP(mobile string) string {
	mu.Lock()
	defer mu.Unlock()

	otp := string(rune(100000 + rand.Intn(899999)))
	MockOTPs[mobile] = OTPData{
		OTP:       otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	return otp
}

func VerifyOTP(mobile, otp string) bool {
	mu.Lock()
	defer mu.Unlock()

	data, exists := MockOTPs[mobile]
	if !exists || time.Now().After(data.ExpiresAt) || data.OTP != otp {
		return false
	}
	delete(MockOTPs, mobile)
	return true
}

func SetResetToken(mobile, token string) {
	mu.Lock()
	defer mu.Unlock()
	MockResetTokens[mobile] = TokenData{
		Token:     token,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

func ValidateResetToken(mobile, token string) bool {
	mu.Lock()
	defer mu.Unlock()
	data, exists := MockResetTokens[mobile]
	if !exists || time.Now().After(data.ExpiresAt) || data.Token != token {
		return false
	}
	return true
}

func ClearResetToken(mobile string) {
	mu.Lock()
	defer mu.Unlock()
	delete(MockResetTokens, mobile)
}
