package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method.

	//ACCOUNTS
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/createAccount", app.CreateAccount)
	router.HandlerFunc(http.MethodPost, "/v1/accountDetails", app.AccountDetails)
	router.HandlerFunc(http.MethodPost, "/v1/accounts", app.RetreiveAccounts)
	router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	router.HandlerFunc(http.MethodPost, "/v1/retrieveAccounts", app.RetreiveAccounts)
	router.HandlerFunc(http.MethodPost, "/v1/beneficiary/new", app.NewBeneficiary)
	router.HandlerFunc(http.MethodPost, "/v1/beneficiary", app.GetBeneficiaries)

	//Currency Exchange
	router.HandlerFunc(http.MethodGet, "v1/", app.AvailableCurrenciesHandler)

	// Return the httprouter instance.
	return router
}
