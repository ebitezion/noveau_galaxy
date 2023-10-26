package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()
	fs := http.FileServer((http.Dir("./cmd/web/static")))
	router.Handler("GET", "/", fs)

	//RENDERED PAGES

	router.HandlerFunc(http.MethodGet, "/v1/index", app.RenderIndexPage)
	router.HandlerFunc(http.MethodGet, "/v1/loginpage", app.RenderLoginPage)
	router.HandlerFunc(http.MethodGet, "/v1/signuppage", app.RenderSignUpPage)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//ACCOUNTS v1
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/createAccount", app.CreateAccount)
	router.HandlerFunc(http.MethodPost, "/v1/accountDetails", app.AccountDetails)
	router.HandlerFunc(http.MethodPost, "/v1/accounts", app.RetreiveAccounts)
	router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	router.HandlerFunc(http.MethodPost, "/v1/retrieveAccounts", app.RetreiveAccounts)
	router.HandlerFunc(http.MethodPost, "/v1/beneficiary/new", app.NewBeneficiary)
	router.HandlerFunc(http.MethodPost, "/v1/beneficiary", app.GetBeneficiaries)

	//LOGIN
	router.HandlerFunc(http.MethodPost, "/v1/login", app.AuthLogin)
	router.HandlerFunc(http.MethodPost, "/v1/create", app.AuthCreate)
	router.HandlerFunc(http.MethodPost, "/v1/authindex", app.AuthIndex)
	router.HandlerFunc(http.MethodPost, "/v1/deposit", app.PaymentDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/credit", app.PaymentCreditInitiation)

	//ACCOUNT V2
	router.HandlerFunc(http.MethodPost, "/v1/accounts/create", app.AccountCreate)

	//Currency Exchange
	router.HandlerFunc(http.MethodGet, "/v1/currencyHandler", app.AvailableCurrenciesHandler)

	return router
}
