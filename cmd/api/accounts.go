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
	data.ValidateAccountBioData(v, &Request)
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

// AccountDetails gets the account details of a user rovided a valid accountNumber
func (app *application) AccountDetails(w http.ResponseWriter, r *http.Request) {
	Request := data.User{}
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

	//fetch account details
	result, err := app.models.AccountModel.GetAccountDetails(Request.AccountNumber)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.writeJSON(w, http.StatusOK, envelope{"message": "No records found"}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a success response or an error message
	err = app.writeJSON(w, http.StatusOK, envelope{"responseMessage": "Successful", "accountDetails": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// BalanceEnquiry gets the balance details of a user provided a valid accountNumber
func (app *application) BalanceEnquiry(w http.ResponseWriter, r *http.Request) {
	Request := data.User{}
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

	//fetch account details
	result, err := app.models.AccountModel.GetBalanceDetails(Request.AccountNumber)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.writeJSON(w, http.StatusOK, envelope{"message": "No records found"}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a success response or an error message
	err = app.writeJSON(w, http.StatusOK, envelope{"responseMessage": "Successful", "accountDetails": result}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// RetreiveAccounts retrieves all accounts associated with a user
func (app *application) RetreiveAccounts(w http.ResponseWriter, r *http.Request) {

}

// NewBeneficiary creates a new beneficiary tied to an account
func (app *application) NewBeneficiary(w http.ResponseWriter, r *http.Request) {
	Request := data.Beneficiary{}
	// read the incoming request body
	err := app.readJSON(w, r, &Request)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the request data
	v := validator.New()
	data.ValidateBeneficiaryData(v, &Request)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	//get user_id of account_number

	//create new beneficiary
	err = app.StoreBeneficiary(&Request, 6)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Return a success response or an error message
	err = app.writeJSON(w, http.StatusOK, envelope{"responseMessage": "Beneficiary Successful Created"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// GetBeneficiaries gets all the beneficiaries tied to an account
func (app *application) GetBeneficiaries(w http.ResponseWriter, r *http.Request) {

}
func (app *application) StoreBeneficiary(Data *data.Beneficiary, UserId int) error {
	//Insert Identity data to DB
	query := `
     INSERT INTO beneficiaries ( user_id, full_name, bank_name, bank_account_number, bank_routing_number, swift_code)
    VALUES (?, ?, ?, ?,?,?)
`

	args := []interface{}{UserId, Data.FullName, Data.BankName, Data.BankAccountNumber, Data.BankRoutingNumber, Data.SwiftCode}

	// Insert Identity data to DB, use LAST_INSERT_ID() to get the last generated ID
	_, err := app.models.AccountModel.Insert(query, args)
	if err != nil {
		return err
	}
	return nil
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
	accountType := data.InternalAccount
	currencyCode := data.Nigeria
	currentBalance := 0.0

	query = `
    INSERT INTO accounts
	(user_id, account_number, type, currency_code, balance)
	VALUES (?, ?, ?, ?, ?)`

	// Insert accounts data to DB, use LAST_INSERT_ID() to get the last generated ID
	_, err = app.models.AccountModel.Insert(query, []interface{}{bioDataId, accountNumber, accountType, currencyCode, currentBalance})
	if err != nil {
		return "", err
	}

	// return the generated account
	return accountNumber, nil
}
