package ukaccountgen

import (
	"crypto/rand"
	"math/big"
)

type UKAccountGenerator struct {
}

func New() *UKAccountGenerator {
	return &UKAccountGenerator{}
}

func (u *UKAccountGenerator) generateRandomNumber(length int) string {
	var result string

	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		result += num.String()
	}

	return result
}

func (u *UKAccountGenerator) GenerateUKAccountNumber() string {
	// Generating a random 6-digit account number
	return u.generateRandomNumber(6)
}
