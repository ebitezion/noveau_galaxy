package transaction_processing

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestTransactionProcessor(t *testing.T) {
	tp := NewTransactionProcessor()

	// Add accounts
	tp.AddAccount("A1", 1000.00)
	tp.AddAccount("A2", 500)

	// Deposit and Withdraw
	err := tp.Deposit("A1", 500)
	assert.NoError(t, err)
	err = tp.Withdraw("A2", 200)
	assert.NoError(t, err)

	// Transfer
	err = tp.Transfer("A1", "A2", 300)
	assert.NoError(t, err)

	// Verify account balances
	assert.Equal(t, decimal.NewFromFloat(700), tp.Accounts["A1"].Balance)
	assert.Equal(t, decimal.NewFromFloat(1000), tp.Accounts["A2"].Balance)

	// Test withdrawal with insufficient funds
	err = tp.Withdraw("A2", 1000)
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	// Test deposit to non-existent account
	err = tp.Deposit("A3", 100)
	assert.Error(t, err)
	assert.Equal(t, "account not found", err.Error())
}
