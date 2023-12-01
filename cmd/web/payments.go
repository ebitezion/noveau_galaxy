package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/notifications"
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

	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign

	sendersAccountNumber := r.FormValue("sendersAccountNumber")
	receiversAccountNumber := r.FormValue("receiversAccountNumber")
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := r.FormValue("Amount")

	response, err := payments.ProcessPAIN([]string{token, "pain", "1", sendersDetails, receiversDetails, amount, "CR"})
	if err != nil {
		fmt.Println(err)
	}
	//send notification
	err = Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	if err != nil {
		fmt.Println(err)
	}

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}
func (app *application) PaymentDebitInitiation(w http.ResponseWriter, r *http.Request) {

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

	//for credit only  receivers account number and sender account number is required
	//which is the number before the @ sign

	sendersAccountNumber := r.FormValue("sendersAccountNumber")
	receiversAccountNumber := r.FormValue("receiversAccountNumber")
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := r.FormValue("Amount")

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
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}

	//receivers account number and bank number
	//for deposit only credit is receivers account number is required
	//which is the number before the @ sign

	sendersAccountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")
	sendersBankNumber := os.Getenv("DEPOSIT_BANK_NUMBER")

	receiversAccountNumber := r.FormValue("receiversAccountNumber")
	amount := r.FormValue("Amount")

	sendersDetails := fmt.Sprintf("%s@%s", sendersAccountNumber, sendersBankNumber)
	receiversDetails := receiversAccountNumber + "@"

	//senders account number and bank number

	response, err := payments.ProcessPAIN([]string{token, "pain", "1000", sendersDetails, receiversDetails, amount, "CR"})
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

// BatchTransaction is used by the banking application to process different transactions at the same time
func (app *application) BatchTransaction() {

}

// this function is how to use the notification package
func Notification(token string, sendersAccountNumber string, receiversAccountNumber string, amount string) error {
	sender, err := accounts.FetchAccountMeta(sendersAccountNumber)
	if err != nil {
		return err
	}

	receiver, err := accounts.FetchAccountMeta(receiversAccountNumber)
	if err != nil {
		return err
	}

	users := [2]*accounts.AccountHolderDetails{sender, receiver}

	for i := range users {

		ns := notifications.NotificationService{}
		User := notifications.User{
			ID:       1,
			Username: "adeoluwa",
			Email:    "akanbiadenugba699@gmail.com",
			Phone:    users[i].ContactNumber1,
		}

		notification := notifications.Notification{
			User:    User,
			Message: fmt.Sprintf("Amount of %s was transfered from %s to %s", amount, sender.AccountNumber, receiver.AccountNumber),
		}
		notifications.SendNotification(ns, notification)
	}

	return nil
}
