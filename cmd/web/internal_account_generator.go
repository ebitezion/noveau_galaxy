package main

import (
	"net/http"

	"github.com/ebitezion/backend-framework/internal/ukaccountgen"
)

// GenerateUKAccountNumberHandler to create an internal account no address system compliant with the UK
func (app *application) GenerateUKAccountNumberHandler(w http.ResponseWriter, r *http.Request) {
	generator := ukaccountgen.New()
	accountNumber := generator.GenerateUKAccountNumber()

	env := envelope{
		"responseCode":          "00",
		"status":                "success",
		"internalAccountNumber": accountNumber,
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
