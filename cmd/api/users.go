package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/validator"
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
	AuthLoginData := data.AuthLoginData{}
	// read the incoming request body
	err := app.readJSON(w, r, &AuthLoginData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateAuthLoginData(v, &AuthLoginData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}
	accountNumber, _, _, err := accounts.FetchAuthDetails(AuthLoginData.Username)

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

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", accountNumber, AuthLoginData.Password})
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

// Get token
func (app *application) AuthLoginExternal(w http.ResponseWriter, r *http.Request) {
	AuthLoginData := data.AuthLoginData{}
	// read the incoming request body
	err := app.readJSON(w, r, &AuthLoginData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateAuthLoginData(v, &AuthLoginData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}
	accountNumber, _, _, err := accounts.FetchAuthDetails(AuthLoginData.Username)

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

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", accountNumber, AuthLoginData.Password})
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
func (app *application) AuthCreateExternal(w http.ResponseWriter, r *http.Request) {
	AuthCreateData := data.AuthLoginData{}
	// read the incoming request body
	err := app.readJSON(w, r, &AuthCreateData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateAuthLoginData(v, &AuthCreateData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "5", AuthCreateData.Password, AuthCreateData.Username})

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
	AuthCreateData := data.AuthCreateData{}
	// read the incoming request body
	err := app.readJSON(w, r, &AuthCreateData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateAuthCreateData(v, &AuthCreateData)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "3", AuthCreateData.AccountNumber, AuthCreateData.Password, AuthCreateData.Username})

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
