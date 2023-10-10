package main

import (
	"fmt"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/ukaccountgen"
	"github.com/ebitezion/backend-framework/internal/validator"
)

func (app *application) CreateAccount(w http.ResponseWriter, r *http.Request) {
	Request := data.AccountBioData{}
	// read the incoming request body
	err := app.readJSON(w, r, &Request)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the request data
	v := validator.New()
	data.ValidateUserInformation(v, &Request)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//create accountno
	accountNumber, err := app.Register(Request)
	if err != nil {
		fmt.Println(err)
	}

	// Return a success response or an error message
	err = app.writeJSON(w, http.StatusOK, envelope{"responseMessage": "Successful", "accountNumber": accountNumber}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// Register simply saves the biodata info of the user and return an account no or an error
func (app *application) Register(biodata data.AccountBioData) (accountNumber string, err error) {
	// Insert biodata to DB
	query := `
    INSERT INTO biodata
    (surname, firstName, homeAddress, city, phoneNumber, dateOfBirth, country)
    VALUES (?, ?, ?, ?, ?, ?, ?)
`

	args := []interface{}{
		biodata.Surname, biodata.FirstName, biodata.HomeAddress, biodata.City,
		biodata.PhoneNumber, biodata.DateOfBirth, biodata.Identity.Country,
	}

	results, err := app.models.AccountModel.Insert(query, args)
	if err != nil {
		return "", err
	}

	bioDataId, err := results.LastInsertId()
	if err != nil {
		return "", err
	}

	//Insert Identity data to DB
	query = `
    INSERT INTO identity
    (bvn, passport, utilityBill, picture, user_id)
    VALUES (?, ?, ?, ?,?)
`

	args = []interface{}{
		biodata.Identity.BVN, biodata.Identity.Passport, biodata.Identity.UtilityBill,
		biodata.Picture, bioDataId,
	}

	// Insert Identity data to DB, use LAST_INSERT_ID() to get the last generated ID
	_, err = app.models.AccountModel.Insert(query, args)
	if err != nil {
		return "", err
	}

	// Generate the account number here (as you've done)
	generator := ukaccountgen.New()
	accountNumber = generator.GenerateUKAccountNumber()

	query = `
    INSERT INTO accounts
    (user_id, internal_act_no)
    VALUES (?, ?)
`

	// Insert accounts data to DB, use LAST_INSERT_ID() to get the last generated ID
	_, err = app.models.AccountModel.Insert(query, []interface{}{bioDataId, accountNumber})
	if err != nil {
		return "", err
	}

	// return the generated account
	return accountNumber, nil
}
