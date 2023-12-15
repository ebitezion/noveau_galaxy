package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/ukaccountgen"
	"github.com/ebitezion/backend-framework/internal/validator"
	Validate "github.com/go-playground/validator/v10"
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

// RetreiveAccounts retrieves all accounts associated with a user
func (app *application) RetreiveAccounts(w http.ResponseWriter, r *http.Request) {
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

	//get user_id
	user_id, err := app.models.AccountModel.GetUserId(Request.AccountNumber)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.writeJSON(w, http.StatusOK, envelope{"message": "No records found"}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//fetch account details
	accounts, err := app.models.AccountModel.GetAccounts(user_id)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.writeJSON(w, http.StatusOK, envelope{"message": "No records found"}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a success response or an error message
	err = app.writeJSON(w, http.StatusOK, envelope{"responseMessage": "Successful", "accounts": accounts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
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
	//get user_id
	userId, err := app.models.AccountModel.GetUserId(Request.UserAccountNumber)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.writeJSON(w, http.StatusOK, envelope{"message": "No records found"}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//create new beneficiary
	err = app.StoreBeneficiary(&Request, userId)
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

	//get user_id
	user_id := "6"

	//fetch account details
	beneficiaries, err := app.models.AccountModel.GetBenefciaries(user_id)
	if err != nil {
		if err == data.ErrRecordNotFound {
			app.writeJSON(w, http.StatusOK, envelope{"message": "No records found"}, nil)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Return a success response or an error message
	err = app.writeJSON(w, http.StatusOK, envelope{"responseMessage": "Successful", "beneficiaries": beneficiaries}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) StoreBeneficiary(Data *data.Beneficiary, UserId string) error {
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
    (bvn, passport, utilityBill, picture, userId)
    VALUES (?, ?, ?, ?,?)`

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
	(userId, accountNumber, type, currencyCode, balance)
	VALUES (?, ?, ?, ?, ?)`

	// Insert accounts data to DB, use LAST_INSERT_ID() to get the last generated ID
	_, err = app.models.AccountModel.Insert(query, []interface{}{bioDataId, accountNumber, accountType, currencyCode, currentBalance})
	if err != nil {
		return "", err
	}

	// return the generated account
	return accountNumber, nil
}

// CreateTransaction stores a transaction that occured in the database
func (app *application) CreateTransaction(transaction *data.Transaction) error {
	//Insert Identity data to DB
	query := `INSERT INTO transactions (sender_account_id, receiver_account_id, amount, currency_code, status, transaction_type, timestamp) VALUES (?, ?, ?, ?, ?, ?,?)`

	args := []interface{}{transaction.SenderAccountID, transaction.ReceiverAccountID, transaction.Amount, transaction.CurrencyCode, transaction.Status, transaction.TransactionType}

	// Insert Identity data to DB, use LAST_INSERT_ID() to get the last generated ID
	_, err := app.models.AccountModel.Insert(query, args)
	if err != nil {
		return err
	}
	return nil
}

// /////////////////////////////////////////////////////
// Version 2.0
// @author Ebite Ogochukwu Zion
// Team Lead
// Provides accounts management features
// ////////////////////////////////////////////////////

type AccountRequest struct {
	Token                   string `json:"token"`
	AccountHolderGivenName  string `json:"accountHolderGivenName" validate:"required,min=6,max=20"`
	AccountHolderFamilyName string `json:"accountHolderFamilyName" validate:"required,min=6,max=20"`
	//AccountHolderDateOfBirth string `json:"accountHolderDateOfBirth" validate:"required,date=2006-01-02"`
	AccountHolderDateOfBirth string `json:"accountHolderDateOfBirth" validate:"required,customDate"`

	AccountHolderIdentificationNum string `json:"accountHolderIdentificationNumber" validate:"required,min=6,max=20"`
	AccountHolderContactNumber2    string `json:"accountHolderContactNumber2"`
	AccountHolderEmailAddress      string `json:"accountHolderEmailAddress" validate:"required,email"`

	AccountHolderContactNumber1 string `json:"accountHolderContactNumber1" validate:"required,len=10"`
	AccountHolderAddressLine1   string `json:"accountHolderAddressLine1" validate:"required,min=6,max=20"`
	AccountHolderAddressLine2   string `json:"accountHolderAddressLine2" validate:"required,min=6,max=20"`
	AccountHolderAddressLine3   string `json:"accountHolderAddressLine3" validate:"required,min=6,max=20"`
	AccountHolderPostalCode     string `json:"accountHolderPostalCode" validate:"required,len=6"`
	AccountNumber               string `json:"accountNumber validate:"required,len=10"`
}

func (app *application) AccountIndex(w http.ResponseWriter, r *http.Request) {
	var req AccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.errorJSON(w, err)
		return
	}

	if err := validateAccountRequest(req); err != nil {
		app.errorJSON(w, err)
		return
	}

	response, err := accounts.ProcessAccount([]string{req.Token, "acmt", "1001"})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

func (app *application) AccountCreate(w http.ResponseWriter, r *http.Request) {

	//_, err := app.getTokenFromHeader(w, r)
	// if err != nil {
	// 	// there was error
	// 	data := envelope{
	// 		"responseCode": "06",
	// 		"status":       "Failed",
	// 		"message":      err.Error(),
	// 	}

	// 	app.writeJSON(w, http.StatusBadRequest, data, nil)
	// 	return
	// }

	var req data.NewAccountRequest
	// read the incoming request body
	err := app.readJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateNewAccountRequestData(v, &req)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	reqSlice := []string{
		"0",
		"acmt",
		"1",
		req.AccountHolderGivenName,
		req.AccountHolderFamilyName,
		req.AccountHolderDateOfBirth,
		req.AccountHolderIdentificationNum,
		req.AccountHolderContactNumber1,
		req.AccountHolderContactNumber2,
		req.AccountHolderEmailAddress,
		req.AccountHolderAddressLine1,
		req.AccountHolderAddressLine2,
		req.AccountHolderAddressLine3,
		req.AccountHolderPostalCode,
		req.Image,
		req.AccountHolderIdentificationType,
		req.Country,
	}

	response, err := accounts.ProcessAccount(reqSlice)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) AccountUpdate(w http.ResponseWriter, r *http.Request) {
	_, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}
	// Parse form data
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}
	accountNumber := r.FormValue("accountNumber")
	accountHolderGivenName := r.FormValue("accountHolderGivenName")
	accountHolderFamilyName := r.FormValue("accountHolderFamilyName")
	accountHolderDateOfBirth := r.FormValue("accountHolderDateOfBirth")
	accountHolderIdentificationNumber := r.FormValue("accountHolderIdentificationNumber")
	accountHolderIdentificationType := r.FormValue("accountHolderIdentificationType")
	accountHolderContactNumber1 := r.FormValue("accountHolderContactNumber1")
	accountHolderContactNumber2 := r.FormValue("accountHolderContactNumber2")
	accountHolderEmailAddress := r.FormValue("accountHolderEmailAddress")
	accountHolderAddressLine1 := r.FormValue("accountHolderAddressLine1")
	accountHolderAddressLine2 := r.FormValue("accountHolderAddressLine2")
	accountHolderAddressLine3 := r.FormValue("accountHolderAddressLine3")
	accountHolderPostalCode := r.FormValue("accountHolderPostalCode")
	accountHolderCountry := r.FormValue("country")

	profileImage := ""

	// Initialize variables with actual data

	req := []string{
		"0",
		"acmt",
		"1007",
		accountHolderGivenName,
		accountHolderFamilyName,
		accountHolderDateOfBirth,
		accountHolderIdentificationNumber,
		accountHolderContactNumber1,
		accountHolderContactNumber2,
		accountHolderEmailAddress,
		accountHolderAddressLine1,
		accountHolderAddressLine2,
		accountHolderAddressLine3,
		accountHolderPostalCode,
		profileImage,
		accountHolderIdentificationType,
		accountHolderCountry,
		accountNumber,
	}

	response, err := accounts.ProcessAccount(req)
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}
	//Response(response, err, w, r)

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)

}
func (app *application) AccountGet(w http.ResponseWriter, r *http.Request) {
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	var req data.User
	// read the incoming request body
	err = app.readJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateUser(v, &req)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}
	response, err := accounts.ProcessAccount([]string{token, "acmt", "1002", req.AccountNumber})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

func (app *application) AccountGetAll(w http.ResponseWriter, r *http.Request) {
	var req AccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		app.errorJSON(w, err)
		return
	}

	if err := validateAccountRequest(req); err != nil {
		app.errorJSON(w, err)
		return
	}

	response, err := accounts.ProcessAccount([]string{req.Token, "acmt", "1000"})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

func (app *application) BalanceEnquiry(w http.ResponseWriter, r *http.Request) {
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	var req data.User
	// read the incoming request body
	err = app.readJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateUser(v, &req)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1003", req.AccountNumber})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

func (app *application) AccountHistory(w http.ResponseWriter, r *http.Request) {
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	var req data.User
	// read the incoming request body
	err = app.readJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateUser(v, &req)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1004", req.AccountNumber})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

// func validateAccountRequest(req AccountRequest) error {
// 	// Perform your validation checks here
// 	// For instance, check if required fields are not empty, validate formats, etc.

// 	if req.AccountHolderGivenName == "" || req.AccountHolderFamilyName == "" {
// 		return errors.New("account holder names are required")
// 	}
// 	// Perform other validations as needed...

// 	return nil
// }

func (app *application) errorJSON(w http.ResponseWriter, err error) {
	data := envelope{
		"responseCode": "06",
		"status":       "Failed",
		"message":      err.Error(),
	}
	app.writeJSON(w, http.StatusBadRequest, data, nil)
}

func validateAccountRequest(req AccountRequest) error {
	validator := Validate.New()
	validator.RegisterValidation("customDate", customDateValidator)
	// Perform validation
	if err := validator.Struct(req); err != nil {
		return err
	}

	return nil
}

func customDateValidator(fl Validate.FieldLevel) bool {
	dateStr := fl.Field().String()

	// Define your date format
	layout := "2006-01-02"
	_, err := time.Parse(layout, dateStr)

	return err == nil
}
