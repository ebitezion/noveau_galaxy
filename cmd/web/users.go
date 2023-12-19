package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/notifications"
)

func (app *application) getTokenFromHeader(w http.ResponseWriter, r *http.Request) (token string, err error) {
	// Get token from header
	token = r.Header.Get("X-Auth-Token")

	if token == "" {

		return "", errors.New("could not retrieve token from headers")

	}
	// Check token
	err = appauth.CheckToken(token)

	if err != nil {
		return "", errors.New("token invalid")
	}

	return token, nil
}

// Extend token
func (app *application) AuthIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Extend token")
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

	//Extend token
	response, err := appauth.ProcessAppAuth([]string{token, "appauth", "1"})
	if err != nil {
		app.logger.Println(err)
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return

	}

	app.logger.Println(response)
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}

	app.writeJSON(w, http.StatusOK, data, nil)

}

// Get token
func (app *application) AuthLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get token")
	password := r.FormValue("password")
	username := r.FormValue("username")
	//get accountnumber
	accountNumber, fullname, role, err := accounts.FetchAuthDetails(username)
	fmt.Println(role)
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

	//get token
	token, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", accountNumber, password})
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

	//set jwt token
	err = app.SetJwtSession(w, r, accountNumber, fullname)
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
	// err = app.LoginNotification(token, accountNumber)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//send success response
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      token,
	}
	app.writeJSON(w, http.StatusOK, data, nil)

}

// Create auth account
func (app *application) AuthCreate(w http.ResponseWriter, r *http.Request) {

	accountNumber := r.FormValue("accountNumber")
	password := r.FormValue("password")
	username := r.FormValue("username")

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "3", accountNumber, password, username})

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
	app.logger.Println(response)
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)
}

// Remove auth account
func (app *application) AuthRemove(w http.ResponseWriter, r *http.Request) {
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

	// user := r.FormValue("User")
	// password := r.FormValue("Password")

	user := app.readString(r.PostForm, "User", "")
	password := app.readString(r.PostForm, "Password", "")

	response, err := appauth.ProcessAppAuth([]string{token, "appauth", "4", user, password})
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err,
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
	}
	app.logger.Println(response)
	data := envelope{
		"responseCode": "00",
		"status":       "Success",
		"message":      response,
	}
	app.writeJSON(w, http.StatusOK, data, nil)

}

// this signs a user from there session
func (app *application) StaffSignOutProcess(w http.ResponseWriter, r *http.Request) {
	app.DeleteSession(w, r, "JwtToken")
	http.Redirect(w, r, "/v1/loginpage", http.StatusFound)
	log.Println("logout successfully")

}

// this deletes a users session token
func (app *application) DeleteSession(w http.ResponseWriter, r *http.Request, sessionname string) {
	// Get the session cookie
	session, err := store.Get(r, sessionname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the MaxAge to a negative value to delete the session cookie
	session.Options.MaxAge = -1

	// Save the session to apply the changes
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// this function is how to use the notification package
func (app *application) LoginNotification(token string, sendersAccountNumber string) error {
	sender, err := accounts.FetchAccountMeta(sendersAccountNumber)
	if err != nil {
		return err
	}

	ns := notifications.NotificationService{}
	User := notifications.User{
		ID:       1,
		Username: "adeoluwa",
		Email:    "akanbiadenugba699@gmail.com",
		Phone:    sender.ContactNumber1,
	}

	notification := notifications.Notification{
		User:    User,
		Message: fmt.Sprintf("Login Notification."),
	}
	notifications.SendNotification(ns, notification)

	return nil
}
