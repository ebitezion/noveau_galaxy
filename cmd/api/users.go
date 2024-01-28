package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name        string `json:"name"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		PhoneNumber string `json:"phoneNumber"`
	}
	// Parse the request body into the anonymous struct.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.Users{
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		Activated: false,
	}
	// Use the Password.Set() method to generate and store the hashed and plaintext
	// passwords.
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	// Validate the user struct and return the error messages to the client if any of
	// the checks fail.
	if data.ValidateUsers(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Insert the user data into the database.
	err = app.models.UserModel.Insert(user)
	if err != nil {
		//TODO 3: Keep LOGS of failed Registration as background task
		switch {
		// If we get a ErrDuplicateEmail error, use the v.AddError() method to manually
		// add a message to the validator instance, and then call our
		// failedValidationResponse() helper.
		case errors.Is(err, data.ErrDuplicateEmailOrUsername):
			v.AddError("email", "a user with this email/username address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Add the "accounts:read" permission for the new user.
	err = app.models.Permissions.AddForUser(user.ID, "account:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// After the user record has been created in the database, generate a new activation
	// token for the user.
	_, err = app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//TODO 4: Keep LOGS of successful Registration as background task

	// Call the Send() method on our Mailer, passing in the user's email address,
	// name of the template file, and the Users struct containing the new user's data.

	// Launch a goroutine which runs an anonymous function that sends the welcome email.

	// app.background(func() {
	// 	// Run a deferred function which uses recover() to catch any panic, and log an
	// 	// error message instead of terminating the application.
	// 	// As there are now multiple pieces of data that we want to pass to our email
	// 	// templates, we create a map to act as a 'holding structure' for the data. This
	// 	// contains the plaintext version of the activation token for the user, along
	// 	// with their ID.
	// 	data := map[string]interface{}{
	// 		"activationToken": token.Plaintext,
	// 		"userID":          user.ID,
	// 	}
	// 	err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
	// 	if err != nil {
	// 		// Importantly, if there is an error sending the email then we use the
	// 		// app.logger.PrintError() helper to manage it, instead of the
	// 		// app.serverErrorResponse() helper like before.

	// 		app.logger.Println(err)

	// 	}
	// })
	// Write a JSON response containing the user data along with a 201 Created status
	// code.
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
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

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", AuthLoginData.Email, AuthLoginData.Password})
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

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", AuthLoginData.Email, AuthLoginData.Password})
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

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "5", AuthCreateData.Password, AuthCreateData.Email})

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
	//generate verification Token
	token, err := app.generateRandomNumber(6)
	if err != nil {
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}
	tokenstr := strconv.Itoa(token)
	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "6", AuthCreateData.Password, AuthCreateData.Email, AuthCreateData.Phone, tokenstr})

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
	//send notification
	err = app.VerificationNotification(token, AuthCreateData.Email)
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

func (app *application) VerifyToken(w http.ResponseWriter, r *http.Request) {
	Token := data.VerifyToken{}
	// read the incoming request body
	err := app.readJSON(w, r, &Token)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateTokenData(v, &Token)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "7", Token.Email, Token.Token})

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
func (app *application) GenerateToken(w http.ResponseWriter, r *http.Request) {
	Email := data.Email{}
	// read the incoming request body
	err := app.readJSON(w, r, &Email)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Validate the user ID
	v := validator.New()
	data.ValidateEmail(v, Email.Email)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		//log to db

		return
	}

	//generate verification Token
	token, err := app.generateRandomNumber(6)
	if err != nil {
		data := envelope{
			"responseCode": "06",
			"status":       "Failed",
			"message":      err.Error(),
		}

		app.writeJSON(w, http.StatusBadRequest, data, nil)
		return
	}
	tokenstr := strconv.Itoa(token)
	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "8", Email.Email, tokenstr})

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
	//send notification
	err = app.VerificationNotification(token, Email.Email)
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
