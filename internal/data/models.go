package data

import (
	"database/sql"
	"errors"
)

var (
	// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
	// looking up a movie that doesn't exist in our database.
	ErrRecordNotFound = errors.New("record not found")

	// Define a custom ErrEditConflict error. We'll return this from our Update() method
	// when there is a data race.
	ErrEditConflict = errors.New("edit conflict")
)

type CurrencyCode string

const (
	Nigeria CurrencyCode = "NGN"
	Uk      CurrencyCode = "GBP"
	Euro    CurrencyCode = "EUR"
)

type AccountType string

const (
	InternalAccount AccountType = "internal"
	UkAccount       AccountType = "uk"
	EuroAccount     AccountType = "euro"
)

type StatusEnum string

const (
	StatusPending   StatusEnum = "pending"
	StatusCompleted StatusEnum = "completed"
	StatusFailed    StatusEnum = "failed"
)

type TypeEnum string

const (
	TypeDeposit    TypeEnum = "deposit"
	TypeWithdrawal TypeEnum = "withdrawal"
	TypeTransfer   TypeEnum = "transfer"
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.

type Models struct {
	AccountModel *AccountModel
	UserModel    *UserModel
	Tokens       *TokenModel
	Permissions  *PermissionModel
}

// For ease of use, we also add a New() method which returns a Models struct containing
// the intitialized savings_accountModel.
func NewModels(db *sql.DB) Models {
	return Models{
		UserModel:    &UserModel{DB: db},
		AccountModel: &AccountModel{DB: db},
		Tokens:       &TokenModel{DB: db},
		Permissions:  &PermissionModel{DB: db},
	}
}
