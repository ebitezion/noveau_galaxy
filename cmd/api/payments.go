package main

import (
	"fmt"
	"net/http"
	"os"

	cashpickup "github.com/ebitezion/backend-framework/internal/cash_pickup"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/payments"
	"github.com/ebitezion/backend-framework/internal/validator"
)

func (app *application) FullAccessTransferInitiation(w http.ResponseWriter, r *http.Request) {
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
	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign

	PaymentInitiationData := data.PaymentInitiationData{}
	// read the incoming request body
	err = app.readJSON(w, r, &PaymentInitiationData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateCreditInitiationData(v, &PaymentInitiationData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	sendersAccountNumber := PaymentInitiationData.SendersAccountNumber
	receiversAccountNumber := PaymentInitiationData.ReceiversAccountNumber
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := PaymentInitiationData.Amount
	initiator := sendersAccountNumber

	response, err := payments.ProcessPAIN([]string{token, "pain", "13", sendersDetails, receiversDetails, amount, "CR", initiator})

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

	//send notification
	err = app.Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	if err != nil {
		fmt.Println(err)
	}
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response + "Transfer Made Successfully",
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

func (app *application) FullAccessDepositInitiation(w http.ResponseWriter, r *http.Request) {

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
	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign
	DepositInitiationData := data.DepositInitiationData{}
	// read the incoming request body
	err = app.readJSON(w, r, &DepositInitiationData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateDepositInitiationData(v, &DepositInitiationData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}
	sendersAccountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")
	receiversAccountNumber := DepositInitiationData.AccountNumber
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := DepositInitiationData.Amount
	initiator := sendersAccountNumber

	response, err := payments.ProcessPAIN([]string{token, "pain", "14", sendersDetails, receiversDetails, amount, "DR", initiator})

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

	// //send notification
	// err = app.Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response + "Deposit Made Sucessfully",
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) PaymentCreditInitiation(w http.ResponseWriter, r *http.Request) {

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

	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign
	PaymentInitiationData := data.PaymentInitiationData{}
	// read the incoming request body
	err = app.readJSON(w, r, &PaymentInitiationData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateCreditInitiationData(v, &PaymentInitiationData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	sendersAccountNumber := PaymentInitiationData.SendersAccountNumber
	receiversAccountNumber := PaymentInitiationData.ReceiversAccountNumber
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := PaymentInitiationData.Amount

	response, err := payments.ProcessPAIN([]string{token, "pain", "1", sendersDetails, receiversDetails, amount, "CR"})

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

	//send notification
	// err = Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response + "Credit Made Successfully",
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) PaymentDebitInitiation(w http.ResponseWriter, r *http.Request) {

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
	PaymentInitiationData := data.PaymentInitiationData{}
	// read the incoming request body
	err = app.readJSON(w, r, &PaymentInitiationData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateCreditInitiationData(v, &PaymentInitiationData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	sendersAccountNumber := PaymentInitiationData.SendersAccountNumber
	receiversAccountNumber := PaymentInitiationData.ReceiversAccountNumber
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := PaymentInitiationData.Amount
	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign

	fmt.Println(sendersAccountNumber, receiversAccountNumber)

	response, err := payments.ProcessPAIN([]string{token, "pain", "1", sendersDetails, receiversDetails, amount, "DR"})
	if err != nil {
		fmt.Println(err)
	}
	//send notification

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) PaymentDepositInitiation(w http.ResponseWriter, r *http.Request) {
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
	DepositInitiationData := data.DepositInitiationData{}
	// read the incoming request body
	err = app.readJSON(w, r, &DepositInitiationData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateDepositInitiationData(v, &DepositInitiationData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}
	sendersAccountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")

	receiversAccountNumber := DepositInitiationData.AccountNumber

	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := DepositInitiationData.Amount

	response, err := payments.ProcessPAIN([]string{token, "pain", "1", sendersDetails, receiversDetails, amount, "CR"})
	if err != nil {
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}
	//send notification
	// err = Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response + "Deposit Made Sucessfully",
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

// BatchTransaction is used by the banking application to process different transactions at the same time
func (app *application) BatchTransaction() {

}

func (app *application) CashPickup(w http.ResponseWriter, r *http.Request) {
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
	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign
	CashPickupData := cashpickup.CashPickup{}
	// read the incoming request body
	err = app.readJSON(w, r, &CashPickupData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateCashPickupData(v, &CashPickupData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	result, err := cashpickup.NewCashPickup(CashPickupData)

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

	// //send notification
	// err = app.Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data := envelope{
		"responseCode":     "00",
		"status":           "CashPickup Created Successfully",
		"reference_number": result,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
