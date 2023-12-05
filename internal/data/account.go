package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type NewAccountRequest struct {
	AccountHolderGivenName          string `json:"accountHolderGivenName"`
	AccountHolderFamilyName         string `json:"accountHolderFamilyName"`
	AccountHolderDateOfBirth        string `json:"accountHolderDateOfBirth"`
	AccountHolderIdentificationNum  string `json:"accountHolderIdentificationNumber"`
	AccountHolderIdentificationType string `json:"accountHolderIdentificationType"`
	AccountHolderContactNumber1     string `json:"accountHolderContactNumber1"`
	AccountHolderContactNumber2     string `json:"accountHolderContactNumber2"`
	AccountHolderEmailAddress       string `json:"accountHolderEmailAddress"`
	AccountHolderAddressLine1       string `json:"accountHolderAddressLine1"`
	AccountHolderAddressLine2       string `json:"accountHolderAddressLine2"`
	AccountHolderAddressLine3       string `json:"accountHolderAddressLine3"`
	AccountHolderPostalCode         string `json:"accountHolderPostalCode"`
	AccountNumber                   string `json:"accountNumber"`
	Image                           string `json:"image"`
	Country                         string `json:"country"`
}

type AccountID struct {
	AccountID string `json:"accountID"`
}
type AuthLoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type AuthCreateData struct {
	Username      string `json:"username"`
	AccountNumber string `json:"accountNumber"`
	Password      string `json:"password"`
}
type PaymentInitiationData struct {
	SendersAccountNumber   string `json:"sendersAccountNumber"`
	ReceiversAccountNumber string `json:"receiversAccountNumber"`
	Amount                 string `json:"amount"`
}
type DepositInitiationData struct {
	AccountNumber string `json:"accountNumber"`
	Amount        string `json:"amount"`
}

type AccountDetails struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	AccountNumber string `json:"accountNumber"`
	AccountType   string `json:"accountType"`
	LedgerBalance string `json:"ledgerBalance"`
	CurrencyCode  string `json:"currencyCode"`
	PhoneNumber   string `json:"phoneNumber"`
	BVN           string `json:"bvn"`
}
type BalanceEnquiry struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	AccountNumber string `json:"accountNumber"`
	LedgerBalance string `json:"ledgerBalance"`
	CurrencyCode  string `json:"currencyCode"`
}

type AccountBioData struct {
	Surname     string `json:"surname"`
	FirstName   string `json:"firstName"`
	HomeAddress string `json:"homeAddress"`
	City        string `json:"city"`
	PhoneNumber string `json:"phoneNumber"`
	Picture     string `json:"picture"`
	DateOfBirth string `json:"dateOfBirth"`
	Identity    struct {
		BVN         string `json:"bvn"`
		Passport    string `json:"passport"`
		UtilityBill string `json:"utilityBill"`
		Country     string `json:"country"`
	} `json:"identity"`
}

type User struct {
	AccountNumber string `json:"accountNumber"`
}
type Beneficiary struct {
	UserAccountNumber string `json:"userAccountNumber,omitempty"`
	BeneficiaryID     int64  `json:"beneficiaryId,omitempty"`
	UserID            int64  `json:"userId,omitempty"`
	FullName          string `json:"fullName"`
	BankName          string `json:"bankName"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankRoutingNumber string `json:"bankRoutingNumber"`
	SwiftCode         string `json:"swiftCode"`
}

type Transaction struct {
	TransactionID     int64   `json:"transactionId"`
	SenderAccountID   int64   `json:"senderAccountId"`
	ReceiverAccountID int64   `json:"receiverAccountId"`
	Amount            float64 `json:"amount"`
	CurrencyCode      string  `json:"currencyCode"`
	Status            string  `json:"status"`
	TransactionType   string  `json:"transactionType"`
	Timestamp         string  `json:"timestamp"`
}
type Account struct {
	AccountNumber string  `json:"accountNumber"`
	Type          string  `json:"type"`
	CurrencyCode  string  `json:"currencyCode"`
	Balance       float64 `json:"balance"`
	SortCode      string  `json:"sortCode"`
	SwiftCode     string  `json:"swiftCode"`
	IBAN          string  `json:"iban"`
	RoutingNumber string  `json:"routingNumber"`
	Other         string  `json:"other"`
}

// connection to DB resources
type AccountModel struct {
	DB *sql.DB
}

// GetUserId gets the userid from a particular account_number
func (a AccountModel) GetUserId(accountNumber string) (UserId string, err error) {

	query := `SELECT userId FROM accounts WHERE accountNumber = ? `
	//To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if len(accountNumber) < 1 {
		return "", ErrRecordNotFound
	}

	// Declare a Users struct to hold the data returned by the query.
	var userid string

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err = a.DB.QueryRowContext(ctx, query, accountNumber).Scan(&userid)
	// Handle any errors. If there was no matching referralcode found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.

	if err != nil {
		// Check specifically for the ErrRecordNotFound error.
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrRecordNotFound
		}
		// Handle other errors.
		return "", err
	}

	// Otherwise, return a pointer to the referrer struct.
	return userid, nil
}

// Insert method for inserting a new record in a table.
func (a AccountModel) Insert(Query string, Args []interface{}) (sql.Result, error) {
	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the database query using the provided context and arguments.
	result, err := a.DB.ExecContext(ctx, Query, Args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Update method for updating a specific record from a table.
func (a AccountModel) Update(Args []interface{}, Query string) error {

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.DB.ExecContext(ctx, Query, Args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Get method for fetching a specific record from the users table.
func (a AccountModel) GetAccountDetails(accountNumber string) (*AccountDetails, error) {
	//To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if len(accountNumber) < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT b.firstName, b.surname, b.phoneNumber, a.account_number, a.type, a.balance, a.currency_code, c.bvn
    FROM v1accounts AS a
    INNER JOIN biodata AS b ON a.user_id = b.id
    INNER JOIN identity AS c ON a.user_id = c.user_id
    WHERE a.account_number = ?;
`

	// Declare a Users struct to hold the data returned by the query.
	var accountDetails AccountDetails

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.
	err := a.DB.QueryRowContext(ctx, query, accountNumber).Scan(&accountDetails.FirstName, &accountDetails.LastName, &accountDetails.PhoneNumber, &accountDetails.AccountNumber, &accountDetails.AccountType, &accountDetails.LedgerBalance, &accountDetails.CurrencyCode, &accountDetails.BVN)

	// Handle any errors. If there was no matching referralcode found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.

	if err != nil {
		// Check specifically for the ErrRecordNotFound error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		// Handle other errors.
		return nil, err
	}

	// Otherwise, return a pointer to the referrer struct.
	return &accountDetails, nil
}

// Get method for fetching a specific record from the users table.
func (a AccountModel) GetBalanceDetails(accountNumber string) (*BalanceEnquiry, error) {
	//To avoid making an unnecessary database call, we take a shortcut
	// and return an ErrRecordNotFound error straight away.
	if len(accountNumber) < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
    SELECT b.firstname, b.surname, a.account_number, a.balance, a.currency_code
    FROM v1accounts AS a
    INNER JOIN biodata AS b ON a.user_id = b.id
    INNER JOIN identity AS c ON a.user_id = c.user_id
    WHERE a.account_number = ?;
`

	// Declare a Users struct to hold the data returned by the query.
	var BalanceEnquiry BalanceEnquiry

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.
	err := a.DB.QueryRowContext(ctx, query, accountNumber).Scan(&BalanceEnquiry.FirstName, &BalanceEnquiry.LastName, &BalanceEnquiry.AccountNumber, &BalanceEnquiry.LedgerBalance, &BalanceEnquiry.CurrencyCode)

	// Handle any errors. If there was no matching referralcode found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.

	if err != nil {
		// Check specifically for the ErrRecordNotFound error.
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		// Handle other errors.
		return nil, err
	}

	// Otherwise, return a pointer to the referrer struct.
	return &BalanceEnquiry, nil
}

// GetBenefciaries  gets all the beneficiaries of a user
func (a AccountModel) GetBenefciaries(UserId string) ([]Beneficiary, error) {

	query := `SELECT  full_name, bank_name, bank_account_number, bank_routing_number, swift_code 	FROM beneficiaries WHERE user_id  = ?`
	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the Query() method to execute the query, passing in the context
	// with the deadline as the first argument.
	rows, err := a.DB.QueryContext(ctx, query, UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Beneficiaries []Beneficiary

	//loop through to store each row gotten in the Beneficiary array
	for rows.Next() {
		var B Beneficiary
		err := rows.Scan(&B.FullName, &B.BankName, &B.BankAccountNumber, &B.BankRoutingNumber, &B.SwiftCode)
		if err != nil {
			return nil, err
		}
		Beneficiaries = append(Beneficiaries, B)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(Beneficiaries) == 0 {
		return nil, ErrRecordNotFound
	}

	return Beneficiaries, nil
}

// GetAccounts  gets all the accounts of a user
func (a AccountModel) GetAccounts(UserId string) ([]Account, error) {

	query := "SELECT account_number, type, currency_code, balance, sort_code, swift_code, iban, routing_number, other FROM v1accounts WHERE user_id = ?"

	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the Query() method to execute the query, passing in the context
	// with the deadline as the first argument.
	rows, err := a.DB.QueryContext(ctx, query, UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Accounts []Account

	//loop through to store each row gotten in the Beneficiary array
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.AccountNumber,
			&account.Type,
			&account.CurrencyCode,
			&account.Balance,
			&account.SortCode,
			&account.SwiftCode,
			&account.IBAN,
			&account.RoutingNumber,
			&account.Other)
		if err != nil {
			return nil, err
		}
		Accounts = append(Accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(Accounts) == 0 {
		return nil, ErrRecordNotFound
	}

	return Accounts, nil
}
