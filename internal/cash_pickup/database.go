package cashpickup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ebitezion/backend-framework/internal/configuration"
	"github.com/shopspring/decimal"
)

var Config configuration.Configuration

func SetConfig(config *configuration.Configuration) {
	Config = *config
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

func checkAccountBalanceBalance(accountNumber string) (balance decimal.Decimal, err error) {
	rows, err := Config.Db.Query("SELECT `availableBalance` FROM `accounts` WHERE `accountNumber` = ?", accountNumber)
	if err != nil {
		return decimal.NewFromFloat(0.), errors.New("payments.checkBalance: " + err.Error())
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&balance); err != nil {
			return decimal.NewFromFloat(0.), errors.New("payments.checkBalance: Could not retrieve account details. " + err.Error())
		}
		count++
	}

	if count > 1 {
		return decimal.NewFromFloat(0.), errors.New("payments.checkBalance: More than one account found with uuid")
	}

	return
}
func saveCashPickup(data CashPickup) (err error) {
	// Prepare statement for inserting data
	insertStatement := "INSERT INTO cash_pickup (`sendersAccountNumber`, `firstName`, `lastName`, `currency`, `reason`, `amount`, `charge`,`bvn`,`nin`,`dob`,`email`,`address`,`image`,`state`,`reference_number`,`sex`,`country`,`phone`) "
	insertStatement += "VALUES(?, ?, ?, ?, ?, ?, ?, ?,?,?,?,?,?,?,?,?,?,?)"
	stmtIns, err := Config.Db.Prepare(insertStatement)
	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// The feePerc is a percentage, convert to amount

	transactionAmountDecimal, err := decimal.NewFromString(data.Amount)
	if err != nil {
		return errors.New("Could not convert transaction amount to decimal. " + err.Error())
	}
	feeAmount := data.Charge.Mul(transactionAmountDecimal)
	_, err = stmtIns.Exec(data.SendersAccountNumber, data.FirstName, data.LastName, data.Currency, data.Reason, data.Amount, feeAmount, data.BVN, data.NIN, data.DOB, data.Email, data.Address, data.Image, data.State, data.ReferenceNumber, data.Sex, data.Country, data.Phone)

	if err != nil {
		return errors.New("payments.savePainTransaction: " + err.Error())
	}

	return
}

func GetAllCashPickups() ([]CashPickup, error) {
	query := "SELECT `sendersAccountNumber`, `firstName`, `lastName`, `status`, `sex`, `currency`, `reason`, `amount`, `charge`, `timestamp`, `bvn`, `nin`, `dob`, `email`, `address`, `image`, `state`, `reference_number`, `country`, `phone` FROM cash_pickup"

	var CashPickups []CashPickup // Slice to hold multiple transaction records.

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use the QueryContext method to execute the query, passing in the context.
	rows, err := Config.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate through the result set and scan each row into a Transaction struct.
	for rows.Next() {
		var data CashPickup
		err := rows.Scan(&data.SendersAccountNumber, &data.FirstName, &data.LastName, &data.Status, &data.Sex, &data.Currency, &data.Reason, &data.Amount, &data.Charge, &data.Timestamp, &data.BVN, &data.NIN, &data.DOB, &data.Email, &data.Address, &data.Image, &data.State, &data.ReferenceNumber, &data.Country, &data.Phone)
		if err != nil {
			return nil, err
		}
		CashPickups = append(CashPickups, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return CashPickups, nil
}
func GetUsersCashPickups(accountNumber string) ([]CashPickup, error) {
	query := "SELECT `sendersAccountNumber`, `firstName`, `lastName`, `status`, `sex`, `currency`, `reason`, `amount`, `charge`, `timestamp`, `bvn`, `nin`, `dob`, `email`, `address`, `image`, `state`, `reference_number`, `country`, `phone` FROM cash_pickup WHERE sendersAccountNumber = ?"
	var CashPickups []CashPickup // Slice to hold multiple transaction records.

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
		var data CashPickup
		err := rows.Scan(&data.SendersAccountNumber, &data.FirstName, &data.LastName, &data.Status, &data.Sex, &data.Currency, &data.Reason, &data.Amount, &data.Charge, &data.Timestamp, &data.BVN, &data.NIN, &data.DOB, &data.Email, &data.Address, &data.Image, &data.State, &data.ReferenceNumber, &data.Country, &data.Phone)
		if err != nil {
			return nil, err
		}
		CashPickups = append(CashPickups, data)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return CashPickups, nil
}
