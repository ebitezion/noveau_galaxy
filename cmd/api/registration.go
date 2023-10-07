package main

import (
	"fmt"
	"net/http"

	accountcreation "github.com/ebitezion/backend-framework/internal/account_creation"
	"github.com/ebitezion/backend-framework/internal/data"
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
	data.ValidateUserInformation(v, &Request)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//create accountno
	accountno, err := accountcreation.Register(Request)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(accountno)

	// Return a success response or an error message
}
