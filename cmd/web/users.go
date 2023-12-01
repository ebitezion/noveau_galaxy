package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/appauth"
)

func (app *application) getTokenFromHeader(w http.ResponseWriter, r *http.Request) (token string, err error) {
	// Get token from header
	token = r.Header.Get("X-Auth-Token")

	if token == "" {
		app.badRequestResponse(w, r, errors.New("could not retrieve token from headers"))
		return "", errors.New("could not retrieve token from headers")
	}
	fmt.Println(err)

	// Check token
	err = appauth.CheckToken(token)
	if err != nil {
		return "", errors.New("token invalid")
	}

	return
}

// Extend token
func (app *application) AuthIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Extend token")
	token, err := app.getTokenFromHeader(w, r)
	if err != nil {
		//there was error
		data := envelope{
			"responseCode": "06",
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
	accountNumber := r.FormValue("accountNumber")
	password := r.FormValue("password")
	username := r.FormValue("username")

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", accountNumber, password, username})

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

// Create auth account
func (app *application) AuthCreate(w http.ResponseWriter, r *http.Request) {

	accountNumber := r.FormValue("accountNumber")
	password := r.FormValue("password")
	username := r.FormValue("username")
	fmt.Println(accountNumber, password, username)
	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "3", accountNumber, password, username})
	fmt.Println(response, err)
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
			"responseCode": "06",
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
