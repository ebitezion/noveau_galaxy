package transaction_processing

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Account struct {
	AccountID string
	Balance   decimal.Decimal
}

type Transaction struct {
	TransactionID string
	FromAccount   string
	ToAccount     string
	Amount        decimal.Decimal
}

type TransactionProcessor struct {
	Accounts     map[string]*Account
	Transactions []Transaction
}

func NewTransactionProcessor() *TransactionProcessor {
	return &TransactionProcessor{
		Accounts: make(map[string]*Account),
	}
}
func (tp *TransactionProcessor) AddAccount(accountID string, initialBalance float64) {
	balance := decimal.NewFromFloat(initialBalance)
	account := &Account{
		AccountID: accountID,
		Balance:   balance,
	}
	tp.Accounts[accountID] = account
}

func (tp *TransactionProcessor) Deposit(accountID string, amount float64) error {
	amountDec := decimal.NewFromFloat(amount)
	account, exists := tp.Accounts[accountID]
	if !exists {
		return errors.New("account not found")
	}
	account.Balance = account.Balance.Add(amountDec)
	tp.Transactions = append(tp.Transactions, Transaction{
		TransactionID: fmt.Sprintf("TX%d", len(tp.Transactions)+1),
		ToAccount:     accountID,
		Amount:        amountDec,
	})
	return nil
}

func (tp *TransactionProcessor) Withdraw(accountID string, amount float64) error {
	amountDec := decimal.NewFromFloat(amount)
	account, exists := tp.Accounts[accountID]
	if !exists {
		return errors.New("account not found")
	}
	if account.Balance.LessThan(amountDec) {
		return errors.New("insufficient funds")
	}
	account.Balance = account.Balance.Sub(amountDec)
	tp.Transactions = append(tp.Transactions, Transaction{
		TransactionID: fmt.Sprintf("TX%d", len(tp.Transactions)+1),
		FromAccount:   accountID,
		Amount:        amountDec,
	})
	return nil
}

func (tp *TransactionProcessor) Transfer(fromAccountID, toAccountID string, amount float64) error {
	//amountDec := decimal.NewFromFloat(amount)
	err := tp.Withdraw(fromAccountID, amount)
	if err != nil {
		return err
	}
	err = tp.Deposit(toAccountID, amount)
	if err != nil {
		return err
	}
	return nil
}
