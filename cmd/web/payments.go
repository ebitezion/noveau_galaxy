package main

import (
	"fmt"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/payments"
)

func (app *application) PaymentCreditInitiation(w http.ResponseWriter, r *http.Request) {

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

	// senderDetails := r.FormValue("SenderDetails")
	// recipientDetails := r.FormValue("RecipientDetails")
	sendersDetails := "befcedbd-0f53-48f4-8219-60477bffb9d6@"
	receiversDetails := "befcedbd-0f53-48f4-8219-60477bffb9d6@"
	amount := r.FormValue("Amount")

	response, err := payments.ProcessPAIN([]string{token, "pain", "1", sendersDetails, receiversDetails, amount})
	if err != nil {
		fmt.Println(err)
	}
	app.logger.Println(response)
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
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	//accountDetails := r.FormValue("AccountDetails")
	//receivers account number and bank number
	sendersDetails := "befcedbd-0f53-48f4-8219-60477bffb9d6@a0299975-b8e2-4358-8f1a-911ee12dbaac"
	receiversDetails := "befcedbd-0f53-48f4-8219-60477bffb9d6@"
	//accountDetails := fmt.Sprintf("%s@%s", sendersDetails, receiversDetails)

	//senders account number and bank number
	amount := r.FormValue("Amount")

	response, err := payments.ProcessPAIN([]string{token, "pain", "1000", sendersDetails, receiversDetails, amount})
	if err != nil {
		fmt.Println(err)
	}
	app.logger.Println(response)
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
