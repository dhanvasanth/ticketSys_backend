// project/internal/utils/otp.go
package utils

import (
    "crypto/rand"
    "fmt"
    "math/big"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
    // Generate 6-digit OTP
    max := big.NewInt(1000000) // 10^6
    min := big.NewInt(100000)  // 10^5
    
    n, err := rand.Int(rand.Reader, new(big.Int).Sub(max, min))
    if err != nil {
        return "", fmt.Errorf("failed to generate random number: %w", err)
    }
    
    otp := new(big.Int).Add(n, min)
    return fmt.Sprintf("%06d", otp.Int64()), nil
}

// ValidateOTP validates the format of OTP
func ValidateOTP(otp string) bool {
    if len(otp) != 6 {
        return false
    }
    
    for _, char := range otp {
        if char < '0' || char > '9' {
            return false
        }
    }
    
    return true
}