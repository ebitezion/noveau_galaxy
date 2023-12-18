package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/notifications"
	"github.com/ebitezion/backend-framework/internal/payments"
	"github.com/ebitezion/backend-framework/internal/rbac_2"
)

type RolePrivileges map[rbac_2.Role][]rbac_2.Privilege

func (app *application) FullAccessCreditInitiation(w http.ResponseWriter, r *http.Request) {
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

	sendersAccountNumber := r.FormValue("sendersAccountNumber")
	receiversAccountNumber := r.FormValue("receiversAccountNumber")
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := r.FormValue("Amount")

	response, err := payments.ProcessPAIN([]string{token, "pain", "13", sendersDetails, receiversDetails, amount, "CR"})

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

	sendersAccountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")
	receiversAccountNumber := r.FormValue("receiversAccountNumber")
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := r.FormValue("Amount")

	// Initialize RBAC system and define roles with associated privileges
	rbac := rbac_2.NewRBACWithDB()
	rbac.AddRole("Admin", []rbac_2.Privilege{"privilege_for_painType_14"})
	//rbac.AddRole("Subadmin", []rbac_2.Privilege{"privilege_1", "privilege_2"})
	// Assuming your payments.ProcessPAIN_2 function signature matches
	response, err := payments.ProcessPAIN_2([]string{token, "pain", "14", sendersDetails, receiversDetails, amount, "DR"}, rbac, "username")
	
	//response, err := payments.ProcessPAIN([]string{token, "pain", "14", sendersDetails, receiversDetails, amount, "DR"})

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
	// err = Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response + "Deposit Made Successfully",
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

	sendersAccountNumber := r.FormValue("sendersAccountNumber")
	receiversAccountNumber := r.FormValue("receiversAccountNumber")
	sendersDetails := sendersAccountNumber + "@"
	receiversDetails := receiversAccountNumber + "@"
	amount := r.FormValue("Amount")

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

	// //send notification
	// err = Notification(token, sendersAccountNumber, receiversAccountNumber, amount)
	// if err != nil {
	// 	fmt.Println(err)
	// }

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
			"responseCode": "07",
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

	sendersAccountNumber := os.Getenv("DEPOSIT_ACCOUNT_NUMBER")

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
