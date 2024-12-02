package util

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"

	"github.com/google/uuid"
)

var randomSource = rand.Reader // Global random source

func GenerateShortLink(surveyID uuid.UUID) (string, error) {
	const keyLength = 6
	buffer := make([]byte, keyLength)
	_, err := randomSource.Read(buffer)
	if err != nil {
		return "", errors.New("failed to read random bytes")
	}
	return base64.URLEncoding.EncodeToString(buffer)[:keyLength], nil
}

func GenerateNumericOTP(length int) (string, error) {
	digits := "0123456789"
	otp := ""
	for i := 0; i < length; i++ {
		max := big.NewInt(int64(len(digits)))
		randomIndex, err := rand.Int(randomSource, max)
		if err != nil {
			return "", err
		}
		otp += string(digits[randomIndex.Int64()])
	}
	return otp, nil
}
