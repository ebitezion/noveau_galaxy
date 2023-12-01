package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Define a handler to serve static files
	fs := http.FileServer(http.Dir("cmd/web/static"))
	router.Handler("GET", "/static/*filepath", http.StripPrefix("/static/", fs))

	// Log all requests to the server
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request URL NOT FOUND:", r.URL.Path)
		http.NotFound(w, r)
	})

	//RENDERED PAGES
	router.HandlerFunc(http.MethodGet, "/v1/index", app.RenderIndexPage)
	router.HandlerFunc(http.MethodGet, "/v1/loginpage", app.RenderLoginPage)
	router.HandlerFunc(http.MethodGet, "/v1/signuppage", app.RenderSignUpPage)
	router.HandlerFunc(http.MethodGet, "/v1/createAccountPage", app.RenderCreateAccountPage)
	router.HandlerFunc(http.MethodGet, "/v1/balanceEnquiryPage", app.RenderBalanceEnquiry)
	router.HandlerFunc(http.MethodGet, "/v1/accountHistoryPage", app.RenderAccountHistory)
	router.HandlerFunc(http.MethodGet, "/v1/depositPage", app.RenderDepositInitiationPage)
	router.HandlerFunc(http.MethodGet, "/v1/creditPage", app.RenderCreditInitiationPage)
	router.HandlerFunc(http.MethodGet, "/v1/batchTransactionPage", app.RenderBatchTransactionPage)
	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//ACCOUNTS v1
	// router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/createAccount", app.CreateAccount)
	// router.HandlerFunc(http.MethodPost, "/v1/accountDetails", app.AccountDetails)
	// router.HandlerFunc(http.MethodPost, "/v1/accounts", app.RetreiveAccounts)
	// router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	// router.HandlerFunc(http.MethodPost, "/v1/retrieveAccounts", app.RetreiveAccounts)
	// router.HandlerFunc(http.MethodPost, "/v1/beneficiary/new", app.NewBeneficiary)
	// router.HandlerFunc(http.MethodPost, "/v1/beneficiary", app.GetBeneficiaries)
	router.HandlerFunc(http.MethodGet, "/v1/accounts", app.AccountGet)
	router.HandlerFunc(http.MethodPost, "/v1/login", app.AuthLogin)
	router.HandlerFunc(http.MethodPost, "/v1/create", app.AuthCreate)
	router.HandlerFunc(http.MethodPost, "/v1/authindex", app.AuthIndex)
	router.HandlerFunc(http.MethodPost, "/v1/deposit", app.PaymentDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/credit", app.PaymentCreditInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/debit", app.PaymentDebitInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	router.HandlerFunc(http.MethodPost, "/v1/accountHistory", app.AccountHistory)

	//ACCOUNT V2
	router.HandlerFunc(http.MethodPost, "/v1/accounts/create", app.AccountCreate)

	//Currency Exchange
	router.HandlerFunc(http.MethodGet, "/v1/availableCurrencies", app.AvailableCurrenciesHandler)

	return router
}
