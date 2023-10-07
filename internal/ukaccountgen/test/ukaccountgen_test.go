package test

// import (
// 	"testing"

// 	"github.com/ebitezion/internal/ukaccountgen"
// )

// func TestGenerateUKAccountNumber(t *testing.T) {
// 	generator := ukaccountgen.New()
// 	accountNumber := generator.GenerateUKAccountNumber()

// 	// Check if the generated account number is exactly 6 characters long
// 	if len(accountNumber) != 6 {
// 		t.Errorf("Generated account number is not 6 characters long. Got: %s", accountNumber)
// 	}

// 	// Check if the generated account number consists only of digits
// 	for _, char := range accountNumber {
// 		if char < '0' || char > '9' {
// 			t.Errorf("Generated account number contains non-digit characters. Got: %s", accountNumber)
// 			break
// 		}
// 	}
// }
