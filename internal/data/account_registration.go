package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

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

// cconnection to DB resources
type AccountModel struct {
	DB *sql.DB
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
