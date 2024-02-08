package cashpickup

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

const TRANSACTION_FEE = 0.0001 // 0.01%

type CashPickup struct {
	SendersAccountNumber string          `json:"sendersAccountNumber"`
	FirstName            string          `json:"firstName"`
	LastName             string          `json:"lastName"`
	Status               string          `json:"status"`
	Sex                  string          `json:"sex"`
	Currency             string          `json:"currency"`
	Reason               string          `json:"reason"`
	Amount               string          `json:"amount"`
	Charge               decimal.Decimal `json:"charge"`
	BVN                  string          `json:"bvn"`
	NIN                  string          `json:"nin"`
	DOB                  string          `json:"dob"`
	Image                string          `json:"image"`
	Address              string          `json:"address"`
	State                string          `json:"state"`
	Email                string          `json:"email"`
	Country              string          `json:"country"`
	Phone                string          `json:"phone"`
	ReferenceNumber      string          `json:"referenceNumber"`
	Timestamp            string          `json:"timestamp"`
	UpdatedAt            string          `json:"updated_at"`
}

func NewCashPickup(cashPickupData CashPickup) (result string, err error) {
	transactionAmountDecimal, err := decimal.NewFromString(cashPickupData.Amount)
	if err != nil {
		return "", errors.New("Could not convert transaction amount to decimal. " + err.Error())
	}
	//check if senders accounts is valid
	exists, err := CheckIfAccountIsActive(cashPickupData.SendersAccountNumber)

	if err != nil {
		return "", errors.New("payments.painNewCashPickup: " + err.Error())
	}
	// Check the result.
	if !exists {
		return "", errors.New("payments.painNewCashPickup: " + "Senders Account Not valid")
	}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable, err := checkAccountBalanceBalance(cashPickupData.SendersAccountNumber)
	if err != nil {
		return "", errors.New("payments.painNewCashPickup: " + err.Error())
	}

	// Comparing decimals results in -1 if <
	if balanceAvailable.Cmp(transactionAmountDecimal) == -1 {
		return "", errors.New("payments.painNewCashPickup: Insufficient funds available")
	}
	reference, err := generateRandomNumber(8)
	if err != nil {
		return "", errors.New("payments.painNewCashPickup: Reference Error")
	}
	cashPickupData.ReferenceNumber = strconv.Itoa(reference)
	cashPickupData.Status = "pending"
	charge, err := decimal.NewFromString("50.00")
	if err != nil {
		return "", errors.New("payments.painNewCashPickup: Conversion Error")
	}
	cashPickupData.Charge = charge

	err = saveCashPickup(cashPickupData)
	if err != nil {
		return "", errors.New("payments.painNewCashPickup: " + err.Error())
	}

	return cashPickupData.ReferenceNumber, nil
}

func fetchAllCashPickups() (result *[]CashPickup, err error) {

	data, err := GetAllCashPickups()
	if err != nil {
		return nil, errors.New("accounts.fetchAllCashPickups: " + err.Error())
	}

	return &data, nil
}
func fetchUsersCashPickups(accountNumber string) (result *[]CashPickup, err error) {

	if accountNumber == "" {
		return nil, errors.New("accounts.fetchUsersCashPickups: Account number not present")
	}
	////check if senders accounts is valid
	exists, err := CheckIfAccountIsActive(accountNumber)

	if err != nil {
		return nil, errors.New("payments.fetchUsersCashPickups: " + err.Error())
	}
	// Check the result.
	if !exists {
		return nil, errors.New("payments.fetchUsersCashPickups: " + "Senders Account Not valid")
	}
	data, err := GetUsersCashPickups(accountNumber)
	if err != nil {
		return nil, errors.New("accounts.fetchUsersCashPickups: " + err.Error())
	}

	return &data, nil
}

// generateRandomNumber gives a random number of a given length
func generateRandomNumber(length int) (int, error) {
	if length < 1 {
		return 0, fmt.Errorf("length should be at least 1")
	}

	// Calculate the minimum and maximum values for the specified length
	min := int(math.Pow10(length - 1))
	max := int(math.Pow10(length)) - 1

	if min >= max {
		return 0, fmt.Errorf("invalid length")
	}

	// Initialize the random number generator with a seed based on the current time
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between min and max (inclusive)
	return rand.Intn(max-min+1) + min, nil
}
