package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP(length int) string {
	const digits = "0123456789"
	otp := make([]byte, length)

	for i := range otp {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		otp[i] = digits[n.Int64()]
	}
	return string(otp)
}
