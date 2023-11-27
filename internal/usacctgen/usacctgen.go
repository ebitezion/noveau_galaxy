package usacctgen

import (
	"crypto/rand"
	"math/big"
)

type USAccountGenerator struct {
}

func NewUSAccountGenerator() *USAccountGenerator {
	return &USAccountGenerator{}
}

func (u *USAccountGenerator) generateRandomNumber(length int) string {
	var result string

	for i := 0; i < length; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		result += num.String()
	}

	return result
}

func (u *USAccountGenerator) GenerateUSAccountNumber() string {
	//routingNumber := u.generateRandomNumber(9)  // Example: 9-digit routing number
	accountNumber := u.generateRandomNumber(10) // Example: 10-digit account number
	return accountNumber                        //routingNumber + accountNumber
}

//Test
// func main() {
// 	usAccountGenerator := NewUSAccountGenerator()
// 	usAccount := usAccountGenerator.GenerateUSAccountNumber()
// 	fmt.Println("Generated US Account Number:", usAccount)
// }
