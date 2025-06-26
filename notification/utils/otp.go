package utils

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	otp := ""
	for i := 0; i < length; i++ {
		otp += strconv.Itoa(rand.Intn(10))
	}
	return otp
}
