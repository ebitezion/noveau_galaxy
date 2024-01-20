package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ebitezion/backend-framework/internal/accounts"
	"github.com/ebitezion/backend-framework/internal/appauth"
	"github.com/ebitezion/backend-framework/internal/data"
	"github.com/ebitezion/backend-framework/internal/notifications"
	"github.com/ebitezion/backend-framework/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Create an anonymous struct to hold the expected data from the request body.
	var input struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
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
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//TODO 4: Keep LOGS of successful Registration as background task

	// Call the Send() method on our Mailer, passing in the user's email address,
	// name of the template file, and the Users struct containing the new user's data.

	// Launch a goroutine which runs an anonymous function that sends the welcome email.

	app.background(func() {
		// Run a deferred function which uses recover() to catch any panic, and log an
		// error message instead of terminating the application.
		// As there are now multiple pieces of data that we want to pass to our email
		// templates, we create a map to act as a 'holding structure' for the data. This
		// contains the plaintext version of the activation token for the user, along
		// with their ID.
		data := map[string]interface{}{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			// Importantly, if there is an error sending the email then we use the
			// app.logger.PrintError() helper to manage it, instead of the
			// app.serverErrorResponse() helper like before.
			app.logger.Println(err, nil)
		}
	})
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
	fmt.Println("Get token")
	password := r.FormValue("password")
	email := r.FormValue("email")
	//get accountnumber
	accountNumber, fullname, _, err := accounts.FetchAuthDetails(email)
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
	token, err := appauth.ProcessAppAuth([]string{"0", "appauth", "2", email, password})
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
	email := r.FormValue("email")

	response, err := appauth.ProcessAppAuth([]string{"0", "appauth", "3", accountNumber, password, email})

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

	// user := r.FormValue("Users")
	// password := r.FormValue("Password")

	user := app.readString(r.PostForm, "Users", "")
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
