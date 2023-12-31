package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ebitezion/backend-framework/internal/accounts"
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
	query := `INSERT INTO transactions (sender_account_id, receiver_account_id, amount, currency_code, status, transaction_type, ) VALUES (?, ?, ?, ?, ?, ?)`

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
func (app *application) AccountIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Account Index")
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1001"})
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
	}
	//Response(response, err, w, r)
	app.logger.Println(response)
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)

}

func (app *application) AccountCreate(w http.ResponseWriter, r *http.Request) {

	_, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "07",
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

	profileImage, err := accounts.ImageToBase64FromRequest(r)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize variables with actual data

	req := []string{
		"0",
		"acmt",
		"1",
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
	}

	_, err = accounts.ProcessAccount(req)
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
		"message":      "Account Created Successfully",
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) AccountCreateSpecial(w http.ResponseWriter, r *http.Request) {

	_, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "07",
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
	accountHolderName := r.FormValue("accountHolderName")
	purpose := r.FormValue("purpose")
	creator := r.FormValue("creator")

	// Initialize variables with actual data

	req := []string{
		"0",
		"acmt",
		"1009",
		accountHolderName,
		creator,
		purpose,
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
			"responseCode": "07",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	//	vars := mux.Vars(r)

	accountId := r.FormValue("accountId")

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1002", accountId})
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
	app.logger.Println(response)
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
			"responseCode": "07",
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

	profileImage, err := ImagetoHexacimal(r)
	if err != nil {
		log.Println(err)
		return
	}

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
func (app *application) AccountGetAll(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Account GetAll")
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)

	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1000"})
	app.logger.Println(response)
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

// BalanceEnquiry gets the balance details of a user provided a valid accountNumber
func (app *application) BalanceEnquiry(w http.ResponseWriter, r *http.Request) {

	token, err := app.getTokenFromHeader(w, r)
	fmt.Println(token)
	if err != nil {

		//there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return

	}

	accountNumber := r.FormValue("accountNumber")

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1003", accountNumber})

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

// AccountHistory retrieves the account history of a user
func (app *application) AccountHistory(w http.ResponseWriter, r *http.Request) {

	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return

	}
	accountNumber := r.FormValue("accountNumber")

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1004", accountNumber})

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

func (app *application) AllTransactions(w http.ResponseWriter, r *http.Request) {
	token, err := app.getTokenFromHeader(w, r)

	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1008"})
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) ExcelTransactions(w http.ResponseWriter, r *http.Request) {
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1008"})
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	path, err := createExcelSheet(response)
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

	// Open the file for reading
	file, err := os.Open(path)
	if err != nil {
		// Handle file open error
		app.writeJSON(w, http.StatusInternalServerError, envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}, nil)
		return
	}
	defer file.Close()

	// Set the appropriate headers for the response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=transactions.xlsx")

	// Copy the file content to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		// Handle file copy error
		app.writeJSON(w, http.StatusInternalServerError, envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}, nil)
		return
	}

}
func (app *application) PdfTransactions(w http.ResponseWriter, r *http.Request) {
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		// there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	response, err := accounts.ProcessAccount([]string{token, "acmt", "1008"})
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}
	path, err := createPdf(response)

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

	fmt.Println(path)

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) BlockAccount(w http.ResponseWriter, r *http.Request) {

	_, err := app.getTokenFromHeader(w, r)

	if err != nil {

		//there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return

	}

	accountNumber := r.FormValue("accountNumber")
	response, err := accounts.ProcessAccount([]string{"", "acmt", "1010", accountNumber})
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

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) UnblockAccount(w http.ResponseWriter, r *http.Request) {
	_, err := app.getTokenFromHeader(w, r)

	if err != nil {

		//there was error
		data := envelope{
			"responseCode": "07",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return

	}

	accountNumber := r.FormValue("accountNumber")
	response, err := accounts.ProcessAccount([]string{"", "acmt", "1011", accountNumber})
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

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
