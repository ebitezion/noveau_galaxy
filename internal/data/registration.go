package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UserInformation struct {
	Surname     string `json:"surname"`
	FirstName   string `json:"firstName"`
	HomeAddress string `json:"homeAddress"`
	City        string `json:"city"`
	PhoneNumber string `json:"phoneNumber"`
	Identity    struct {
		BVN         string `json:"bvn"`
		Passport    string `json:"passport"`
		UtilityBill string `json:"utilityBill"`
		Country     string `json:"country"`
	} `json:"identity"`
	Picture string `json:"picture"`
}

type RegistrationModel struct {
	DB *sql.DB
}

// CheckIfValueExists checks if a given value is in the specified table and returns a boolean
func (r RegistrationModel) CheckIfValueExists(Query string, Args []interface{}) (bool, error) {

	// Declare a variable to store the count.
	var count int

	// Use the context.WithTimeout() function to create a context.Context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() to execute the query and get the count.
	err := r.DB.QueryRowContext(ctx, Query, Args...).Scan(&count)
	if err != nil {
		return false, err
	}

	// If the count is greater than 0, the invitee is already in the database.
	return count > 0, nil
}

// Insert method for inserting a new record in the  table.
func (r RegistrationModel) Insert(Query string, Args []interface{}) error {

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the database query using the provided context and arguments.
	_, err := r.DB.ExecContext(ctx, Query, Args...)
	if err != nil {
		return err
	}
	return nil
}

// Update method for updating a specific record from a table.
func (r RegistrationModel) Update(Args []interface{}, Query string) error {

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.DB.ExecContext(ctx, Query, Args...)
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
