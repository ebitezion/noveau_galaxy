package accounts

/*
@TODO Fix DB repetition
*/

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ebitezion/backend-framework/internal/configuration"
	"github.com/twinj/uuid"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
}

func loadDatabase() (db *sql.DB, err error) {
	// Test connection with ping
	err = Config.Db.Ping()
	if err != nil {
		return
	}

	return
}

func createAccount(accountDetails *AccountDetails, accountHolderDetails *AccountHolderDetails) (err error) {
	// Convert variables

	err = doCreateAccount(accountDetails)
	if err != nil {
		return errors.New("accounts.createAccount: " + err.Error())
	}

	err = doCreateAccountMeta(
		accountHolderDetails, accountDetails)
	if err != nil {
		return errors.New("accounts.createAccount: " + err.Error())
	}

	return
}

func deleteAccount(accountDetails *AccountDetails, accountHolderDetails *AccountHolderDetails) (err error) {
	err = doDeleteAccount(accountDetails)
	if err != nil {
		return errors.New("accounts.deleteAccount: " + err.Error())
	}

	err = doDeleteAccountMeta(accountHolderDetails, accountDetails)
	if err != nil {
		return errors.New("accounts.deleteAccount: " + err.Error())
	}

	return
}

func doCreateAccount(accountDetails *AccountDetails) (err error) {
	// Create account
	insertStatement := "INSERT INTO accounts (`accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}

	// Prepare statement for inserting data
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// Generate account number
	newUuid := uuid.NewV4()
	accountDetails.AccountNumber = newUuid.String()

	_, err = stmtIns.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName, accountDetails.AccountBalance, accountDetails.Overdraft, accountDetails.AvailableBalance)
	if err != nil {
		return errors.New("accounts.doCreateAccount: " + err.Error())
	}
	return
}

func doDeleteAccount(accountDetails *AccountDetails) (err error) {
	// Create account
	deleteStatement := "DELETE FROM accounts WHERE `accountNumber` = ? AND `bankNumber` = ? AND `accountHolderName` = ? "
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccount: " + err.Error())
	}

	// Prepare statement for inserting data
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	_, err = stmtDel.Exec(accountDetails.AccountNumber, accountDetails.BankNumber, accountDetails.AccountHolderName)
	if err != nil {
		return errors.New("accounts.doDeleteAccount: " + err.Error())
	}
	// Can use db.RowsAffected() to check
	return
}

func doCreateAccountMeta(accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	// Create account meta
	insertStatement := "INSERT INTO accounts_meta (`accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode`) "
	insertStatement += "VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("accounts.doCreateAccountMeta: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtIns.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode)

	if err != nil {
		return errors.New("accounts.doCreateAccountMeta: " + err.Error())
	}

	return
}

func doDeleteAccountMeta(accountHolderDetails *AccountHolderDetails, accountDetails *AccountDetails) (err error) {
	// Create account meta
	deleteStatement := "DELETE FROM accounts_meta WHERE `accountNumber` = ? AND `bankNumber` = ? AND `accountHolderGivenName` = ? AND `accountHolderFamilyName` = ? AND `accountHolderDateOfBirth` = ? AND `accountHolderIdentificationNumber` = ? AND `accountHolderContactNumber1` = ? AND `accountHolderContactNumber2` = ? AND `accountHolderEmailAddress` = ? AND `accountHolderAddressLine1` = ? AND `accountHolderAddressLine2` = ? AND `accountHolderAddressLine3` = ? AND `accountHolderPostalCode` = ? "
	stmtDel, err := Config.Db.Prepare(deleteStatement)
	if err != nil {
		return errors.New("accounts.doDeleteAccountMeta: " + err.Error())
	}
	defer stmtDel.Close() // Close the statement when we leave main() / the program terminates

	accountHolderDetails.AccountNumber = accountDetails.AccountNumber

	_, err = stmtDel.Exec(accountHolderDetails.AccountNumber, accountHolderDetails.BankNumber, accountHolderDetails.GivenName, accountHolderDetails.FamilyName, accountHolderDetails.DateOfBirth, accountHolderDetails.IdentificationNumber, accountHolderDetails.ContactNumber1, accountHolderDetails.ContactNumber2, accountHolderDetails.EmailAddress, accountHolderDetails.AddressLine1, accountHolderDetails.AddressLine2, accountHolderDetails.AddressLine3,
		accountHolderDetails.PostalCode)

	if err != nil {
		return errors.New("accounts.doDeleteAccountMeta: " + err.Error())
	}

	return
}

func getAccountDetails(id string) (accountDetails AccountDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", id)
	if err != nil {
		return AccountDetails{}, errors.New("accounts.getAccountDetails: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.AccountHolderName, &accountDetails.AccountBalance, &accountDetails.Overdraft, &accountDetails.AvailableBalance)
		if err != nil {
			break
		}
		count++
	}

	if count == 0 {
		return AccountDetails{}, errors.New("accounts.getAccountDetails: Account not found")
	}

	if count > 1 {
		//@TODO: Allow user to have multiple accounts
		return AccountDetails{}, errors.New("accounts.getAccountDetails: More than one account found")
	}

	return
}

func getAccountMeta(id string) (accountDetails AccountHolderDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderGivenName`, `accountHolderFamilyName`, `accountHolderDateOfBirth`, `accountHolderIdentificationNumber`, `accountHolderContactNumber1`, `accountHolderContactNumber2`, `accountHolderEmailAddress`, `accountHolderAddressLine1`, `accountHolderAddressLine2`, `accountHolderAddressLine3`, `accountHolderPostalCode` FROM `accounts_meta` WHERE `accountNumber` = ?", id)
	if err != nil {
		return AccountHolderDetails{}, errors.New("accounts.getAccountMeta: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&accountDetails.AccountNumber, &accountDetails.BankNumber, &accountDetails.GivenName, &accountDetails.FamilyName, &accountDetails.DateOfBirth, &accountDetails.IdentificationNumber, &accountDetails.ContactNumber1, &accountDetails.ContactNumber2, &accountDetails.EmailAddress, &accountDetails.AddressLine1, &accountDetails.AddressLine2,
			&accountDetails.AddressLine3, &accountDetails.PostalCode); err != nil {
			//@TODO Throw error
			break
		}
		count++
	}

	if count == 0 {
		return AccountHolderDetails{}, errors.New("accounts.getAccountMeta: Account not found")
	}

	if count > 1 {
		//@TODO: Allow user to have multiple accounts
		return AccountHolderDetails{}, errors.New("accounts.getAccountMeta: More than one account found")
	}

	return
}

func getAllAccountDetails() (allAccounts []AccountDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName` FROM `accounts`")
	if err != nil {
		return []AccountDetails{}, errors.New("accounts.getAllAccountDetails: Error with select query: " + err.Error())
	}
	defer rows.Close()

	count := 0
	allAccounts = make([]AccountDetails, 0)

	for rows.Next() {
		accountDetailsSingle := AccountDetails{}
		if err := rows.Scan(&accountDetailsSingle.AccountNumber, &accountDetailsSingle.BankNumber, &accountDetailsSingle.AccountHolderName); err != nil {
			break
		}

		allAccounts = append(allAccounts, accountDetailsSingle)
		count++
	}

	return
}

func getSingleAccountDetail(accountNumber string) (account AccountDetails, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber`, `bankNumber`, `accountHolderName`, `accountBalance`, `overdraft`, `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return AccountDetails{}, errors.New("accounts.getSingleAccountDetail: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&account.AccountNumber, &account.BankNumber, &account.AccountHolderName, &account.AccountBalance, &account.Overdraft, &account.AvailableBalance); err != nil {
			break
		}

		count++
	}

	return
}

func getSingleAccountNumberByID(userID string) (accountID string, err error) {
	rows, err := Config.Db.Query("SELECT `accountNumber` FROM `accounts_meta` WHERE `accountHolderIdentificationNumber` = ?", userID)
	if err != nil {
		return "", errors.New("accounts.getSingleAccountNumberByID: " + err.Error())
	}
	defer rows.Close()

	count := 0
	// @TODO Right now this will return the latest account only, if there are two accounts
	for rows.Next() {
		if err := rows.Scan(&accountID); err != nil {
			break
		}
		count++
	}

	if count == 0 {
		return "", errors.New("accounts.getSingleAccountNumberByID: Account not found")
	}

	return
}

// Get method for fetching a specific record from the users table.
func GetBalanceDetails(accountNumber string) (BalanceEnquiry, error) {

	query := "SELECT `accountHolderName`, `accountNumber`, `accountBalance` FROM `accounts` WHERE `accountNumber` = ?"

	// Declare a Users struct to hold the data returned by the query.
	var BalanceEnquiry BalanceEnquiry

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryRowContext() method to execute the query, passing in the context
	// with the deadline as the first argument.
	err := Config.Db.QueryRowContext(ctx, query, accountNumber).Scan(&BalanceEnquiry.AccountHolderName, &BalanceEnquiry.AccountNumber, &BalanceEnquiry.LedgerBalance)

	// Handle any errors. If there was no matching referralcode found, Scan() will return
	// a sql.ErrNoRows error. We check for this and return our custom ErrRecordNotFound
	// error instead.

	if err != nil {
		// Check specifically for the ErrRecordNotFound error.
		if errors.Is(err, sql.ErrNoRows) {
			return BalanceEnquiry, errors.New("accounts.getSingleAccountNumberByID: Account not found")
		}
		// Handle other errors.
		return BalanceEnquiry, err
	}

	// Otherwise, return a pointer to the referrer struct.
	return BalanceEnquiry, nil
}

// Get method for fetching all records from the transactions table for a specific account number.
func GetAccountHistory(accountNumber string) ([]Transaction, error) {
	query := "SELECT transaction, type, senderAccountNumber, senderBankNumber, receiverAccountNumber, receiverBankNumber, transactionAmount, feeAmount, timestamp FROM transactions WHERE senderAccountNumber = ?"

	var transactions []Transaction // Slice to hold multiple transaction records.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryContext method to execute the query, passing in the context.
	rows, err := Config.Db.QueryContext(ctx, query, accountNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and scan each row into a Transaction struct.
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.Transaction, &t.Type, &t.SenderAccountNumber, &t.SenderBankNumber, &t.ReceiverAccountNumber, &t.ReceiverBankNumber, &t.TransactionAmount, &t.FeeAmount, &t.Timestamp)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
