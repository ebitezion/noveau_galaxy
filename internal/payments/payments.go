package payments

/*
PAIN transactions are as follows

Payments initiation:
1 - CustomerCreditTransferInitiationV06
2 - CustomerPaymentStatusReportV06
7 - CustomerPaymentReversalV05
8 - CustomerDirectDebitInitiationV05
9-  CustomerDebitTransferInitiationV05

Payments mandates:
9 - MandateInitiationRequestV04
10 - MandateAmendmentRequestV04
11 - MandateCancellationRequestV04
12 - MandateAcceptanceReportV04


@author adenugba adeoluwa 1st december
13- FullAccessTransferInitiation
14- FullAccessDepositInitiation


#### Custom payments
1000 - CustomerDepositInitiation (@FIXME Will need to implement this properly, for now we use it to demonstrate functionality)

*/

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/rbac_2"
	"github.com/shopspring/decimal"
)

const TRANSACTION_FEE = 0.0001 // 0.01%

// @TODO Have this struct not repeat in payments and accounts
type AccountHolder struct {
	AccountNumber string
	BankNumber    string
}

type PAINTrans struct {
	PainType  int64
	Sender    AccountHolder
	Receiver  AccountHolder
	Amount    decimal.Decimal
	Fee       decimal.Decimal
	Narration string
	Initiator string
}
type TransactionBatch struct {
	Transactions []Transaction
}
type CashPickup struct {
	SendersAccountNumber string    `json:"sendersAccountNumber"`
	FirstName            string    `json:"firstName"`
	LastName             string    `json:"lastName"`
	Status               string    `json:"status"`
	Currency             string    `json:"currency"`
	Reason               string    `json:"reason"`
	Amount               float64   `json:"amount"`
	Charge               float64   `json:"charge"`
	Timestamp            string    `json:"timestamp"`
	BVN                  string    `json:"bvn"`
	NIN                  string    `json:"nin"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type Transaction struct {
	// Define transaction fields
}

var transactionMutex sync.Mutex

func ProcessPAIN_2(data []string, rbac *rbac_2.RBAC, username string) (result string, err error) {
	// Lock the mutex before accessing/modifying shared resources
	transactionMutex.Lock()
	defer transactionMutex.Unlock() // Ensure the mutex is always unlocked

	// There must be at least 3 elements
	if len(data) < 3 {
		return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
	}

	// Get type
	painType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return "", errors.New("payments.ProcessPAIN: Could not get type of PAIN transaction. " + err.Error())
	}

	// Check RBAC permissions before processing each painType
	var requiredPrivilege rbac_2.Privilege
	//requiredPrivilege := ""

	switch painType {
	case 1:
		requiredPrivilege = "privilege_for_painType_1"
	case 9:
		requiredPrivilege = "privilege_for_painType_9"
	case 13:
		requiredPrivilege = "privilege_for_painType_13"
	case 14:
		requiredPrivilege = "privilege_for_painType_14"
	case 1000:
		requiredPrivilege = "privilege_for_painType_1000"
	}

	if requiredPrivilege != "" && !rbac.CheckPermission(username, requiredPrivilege) {
		// User doesn't have required privilege for this painType
		return "", errors.New("payments.ProcessPAIN: User does not have the required privilege")
	}

	switch painType {
	case 1:
		if !rbac.CheckPermission(username, "privilege_for_painType_1") {
			return "", errors.New("payments.ProcessPAIN: User does not have the required privilege for painType 1")
		}
		// Process for painType 1...
	case 9:
		if !rbac.CheckPermission(username, "privilege_for_painType_9") {
			return "", errors.New("payments.ProcessPAIN: User does not have the required privilege for painType 9")
		}
		// Process for painType 9...
	case 14:
		if !rbac.CheckPermission(username, "access_case_14") {
			return "", errors.New("payments.ProcessPAIN: User does not have the required privilege for case 14")
		}
		// Process for painType 14...
		//There must be at least 4 elements
		//token~pain~type~amount
		if len(data) < 5 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
		}

		result, err = painFullAccessDepositInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	// Other painTypes and their associated permissions...
	default:
		return "", errors.New("payments.ProcessPAIN: Invalid painType")
	}

	// Continue processing based on painType...

	return result, nil
}

func ProcessPAIN(data []string) (result string, err error) {
	// Lock the mutex before accessing/modifying shared resources
	transactionMutex.Lock()
	defer transactionMutex.Unlock() // Ensure the mutex is always unlocked

	//There must be at least 3 elements
	if len(data) < 3 {
		return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
	}

	// Get type
	painType, err := strconv.ParseInt(data[2], 10, 64)
	if err != nil {
		return "", errors.New("payments.ProcessPAIN: Could not get type of PAIN transaction. " + err.Error())
	}

	switch painType {
	case 1:
		//There must be at least 6 elements
		if len(data) < 6 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
		}

		result, err = painCreditTransferInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	case 9:
		//There must be at least 6 elements
		if len(data) < 6 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
		}

		result, err = painDebitTransferInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	case 13:
		//There must be at least 4 elements
		//token~pain~type~amount
		if len(data) < 5 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
		}
		result, err = painFullAccessTransferInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break
	case 14:
		//There must be at least 4 elements
		//token~pain~type~amount
		if len(data) < 5 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
		}

		result, err = painFullAccessDepositInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break

	case 1000:
		//There must be at least 4 elements
		//token~pain~type~amount
		if len(data) < 5 {
			return "", errors.New("payments.ProcessPAIN: Not all data is present. Run pain~help to check for needed PAIN data")
		}
		result, err = customerDepositInitiation(painType, data)
		if err != nil {
			return "", errors.New("payments.ProcessPAIN: " + err.Error())
		}
		break

	}

	return
}

func painCreditTransferInitiation(painType int64, data []string) (result string, err error) {

	// Validate input
	sender, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	receiver, err := parseAccountHolder(data[4])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	//check if receivers accounts is valid

	exists, err := CheckIfAccountNumberExists(receiver.AccountNumber)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	// Check the result.
	if !exists {
		return "", errors.New("payments.painCreditTransferInitiation: " + "Receivers Account Not valid")
	}
	////check if senders accounts is valid
	exists, err = CheckIfAccountIsActive(sender.AccountNumber)

	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	// Check the result.
	if !exists {
		return "", errors.New("payments.painCreditTransferInitiation: " + "Senders Account Not valid")
	}

	Narration := data[6]
	Initiator := data[7]
	transaction := PAINTrans{painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), Narration, Initiator}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable, err := checkBalance(transaction.Sender)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	// Comparing decimals results in -1 if <
	if balanceAvailable.Cmp(transaction.Amount) == -1 {
		return "", errors.New("payments.painCreditTransferInitiation: Insufficient funds available")
	}

	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	return
}
func painFullAccessDepositInitiation(painType int64, data []string) (result string, err error) {

	// Validate input
	sender, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}
	//validate sender
	exists, err := CheckIfAccountNumberExists(sender.AccountNumber)
	// Check the result.
	if !exists {
		return "", errors.New("payments.CustomerDepositInitiation: " + "Account Not valid")
	}
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	//validate receiver
	receiver, err := parseAccountHolder(data[4])
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}
	exists, err = CheckIfAccountNumberExists(receiver.AccountNumber)
	// Check the result.
	if !exists {
		return "", errors.New("payments.CustomerDepositInitiation: " + "Account Not valid")
	}
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	Narration := data[6]
	Initiator := data[7]
	transaction := PAINTrans{painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), Narration, Initiator}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable, err := checkBalance(transaction.Sender)
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}
	// Comparing decimals results in -1 if <
	if balanceAvailable.Cmp(transaction.Amount) == -1 {
		return "", errors.New("payments.CustomerDepositInitiation: Insufficient funds available")
	}

	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	return
}
func painFullAccessTransferInitiation(painType int64, data []string) (result string, err error) {

	// Validate input
	sender, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.painFullAccessTransferInitiation: " + err.Error())
	}
	receiver, err := parseAccountHolder(data[4])
	if err != nil {
		return "", errors.New("payments.painFullAccessTransferInitiation: " + err.Error())
	}

	//check if receivers accounts is valid

	exists, err := CheckIfAccountNumberExists(receiver.AccountNumber)
	fmt.Println(exists)
	if err != nil {
		return "", errors.New("payments.painFullAccessTransferInitiation: " + err.Error())
	}

	// Check the result.
	if !exists {
		return "", errors.New("payments.painFullAccessTransferInitiation: " + "Receivers Account Not valid")
	}
	//check if senders accounts is valid
	exists, err = CheckIfAccountIsActive(sender.AccountNumber)
	fmt.Println(exists)
	if err != nil {
		return "", errors.New("payments.painFullAccessTransferInitiation: " + err.Error())
	}
	// Check the result.
	if !exists {
		return "", errors.New("payments.painFullAccessTransferInitiation: " + "Senders Account Not valid")
	}

	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.painFullAccessTransferInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	Narration := data[6]
	Initiator := data[7]
	transaction := PAINTrans{painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), Narration, Initiator}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable, err := checkBalance(transaction.Sender)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	// Comparing decimals results in -1 if <
	if balanceAvailable.Cmp(transaction.Amount) == -1 {
		return "", errors.New("payments.painCreditTransferInitiation: Insufficient funds available")
	}

	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	return
}
func painDebitTransferInitiation(painType int64, data []string) (result string, err error) {

	// Validate input
	sender, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	receiver, err := parseAccountHolder(data[4])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	// Check if sender valid
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	if tokenUser != sender.AccountNumber {
		return "", errors.New("payments.painCreditTransferInitiation: Sender not valid")
	}
	Narration := data[6]
	Initiator := data[7]

	transaction := PAINTrans{painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), Narration, Initiator}

	// Checks for transaction (avail balance, accounts open, etc)
	balanceAvailable, err := checkBalance(transaction.Sender)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}
	// Comparing decimals results in -1 if <
	if balanceAvailable.Cmp(transaction.Amount) == -1 {
		return "", errors.New("payments.painCreditTransferInitiation: Insufficient funds available")
	}

	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.painCreditTransferInitiation: " + err.Error())
	}

	return
}
func processPAINTransaction(transaction PAINTrans) (result string, err error) {
	// Test: pain~1~1b2ca241-0373-4610-abad-da7b06c50a7b@~181ac0ae-45cb-461d-b740-15ce33e4612f@~20

	// Save in transaction table
	err = savePainTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.processPAINTransaction: " + err.Error())
	}

	// Amend sender and receiver accounts
	// Amend bank's account with fee addition
	err = updateAccounts(transaction)
	if err != nil {
		return "", errors.New("payments.processPAINTransaction: " + err.Error())
	}

	return
}
func processExternalPAINTransaction(transaction PAINTrans) (result string, err error) {
	// Test: pain~1~1b2ca241-0373-4610-abad-da7b06c50a7b@~181ac0ae-45cb-461d-b740-15ce33e4612f@~20

	//external api to actually transfer the money

	// verification of payment
	// Save in transaction table
	err = savePainTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.processPAINTransaction: " + err.Error())
	}

	// Amend sender and receiver accounts
	// Amend bank's account with fee addition
	err = updateAccounts(transaction)
	if err != nil {
		return "", errors.New("payments.processPAINTransaction: " + err.Error())
	}

	return
}

func parseAccountHolder(account string) (accountHolder AccountHolder, err error) {
	accountStr := strings.Split(account, "@")

	if len(accountStr) < 2 {
		return AccountHolder{}, errors.New("payments.parseAccountHolder: Not all details present")
	}

	accountHolder = AccountHolder{accountStr[0], accountStr[1]}

	return
}

func customerDepositInitiation(painType int64, data []string) (result string, err error) {
	// Validate input
	// Sender is bank

	//senders account number and bank number
	sender, err := parseAccountHolder(data[3])
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	receiver, err := parseAccountHolder(data[4])
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	exists, err := CheckIfAccountNumberExists(receiver.AccountNumber)

	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}
	// Check the result.
	if !exists {
		return "", errors.New("payments.CustomerDepositInitiation: " + "Account Not valid")
	}
	trAmt := strings.TrimRight(data[5], "\x00")
	transactionAmountDecimal, err := decimal.NewFromString(trAmt)
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: Could not convert transaction amount to decimal. " + err.Error())
	}

	// Check if sender valid
	tokenUser, err := appauth.GetUserFromToken(data[0])
	if err != nil {
		return "", errors.New("payments.customerDepositInitiation: " + err.Error())
	}
	if tokenUser != sender.AccountNumber {
		return "", errors.New("payments.customerDepositInitiation: Sender not valid")
	}
	Narration := data[6]
	Initiator := data[7]
	// Issue deposit
	// @TODO This flow show be fixed. Maybe have banks approve deposits before initiation, or
	// immediate approval below a certain amount subject to rate limiting
	transaction := PAINTrans{painType, sender, receiver, transactionAmountDecimal, decimal.NewFromFloat(TRANSACTION_FEE), Narration, Initiator}
	// Save transaction
	result, err = processPAINTransaction(transaction)
	if err != nil {
		return "", errors.New("payments.CustomerDepositInitiation: " + err.Error())
	}

	return
}

func ProcessTransactionBatch(batch TransactionBatch) error {
	for _, transaction := range batch.Transactions {
		// Process each transaction
		if err := processSingleTransaction(transaction); err != nil {
			return err
		}
	}
	return nil
}

func processSingleTransaction(transaction Transaction) error {
	// Logic to handle a single transaction
	// ...
	return nil
}

// CheckIfValueExists checks if a given value is in the specified table and returns a boolean
func CheckIfAccountNumberExists(accountNumber string) (bool, error) {
	query := "SELECT COUNT(*) FROM accounts WHERE accountNumber = ?;"
	// Declare a variable to store the count.
	var count int

	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() to execute the query and get the count.
	err := Config.Db.QueryRowContext(ctx, query, accountNumber).Scan(&count)
	if err != nil {
		// Print the error message for debugging purposes.
		fmt.Println("Error executing query:", err)
		return false, err
	}

	// If the count is greater than 0, the value exists in the database.
	return count > 0, nil
}

// CheckIfAccountIsActive checks if a given value is in the specified table and returns a boolean
func CheckIfAccountIsActive(accountNumber string) (bool, error) {

	query := "SELECT COUNT(*) FROM accounts WHERE accountNumber = ? AND status = 'Active';"
	// Declare a variable to store the count.
	var count int

	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() to execute the query and get the count.
	err := Config.Db.QueryRowContext(ctx, query, accountNumber).Scan(&count)
	if err != nil {
		// Print the error message for debugging purposes.
		fmt.Println("Error executing query:", err)
		return false, err
	}

	// If the count is greater than 0, the value exists in the database.
	return count > 0, nil
}
