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

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//ACCOUNTS v1
	 router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// router.HandlerFunc(http.MethodPost, "/v1/createAccount", app.CreateAccount)
	// router.HandlerFunc(http.MethodPost, "/v1/accountDetails", app.AccountDetails)
	// router.HandlerFunc(http.MethodPost, "/v1/accounts", app.RetreiveAccounts)
	// router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	// router.HandlerFunc(http.MethodPost, "/v1/retrieveAccounts", app.RetreiveAccounts)
	// router.HandlerFunc(http.MethodPost, "/v1/beneficiary/new", app.NewBeneficiary)
	// router.HandlerFunc(http.MethodPost, "/v1/beneficiary", app.GetBeneficiaries)

	//authentication
	router.HandlerFunc(http.MethodPost, "/v1/api/login", app.AuthLogin)
	router.HandlerFunc(http.MethodPost, "/v1/api/create", app.AuthCreate)
	router.HandlerFunc(http.MethodPost, "/v1/authindex", app.AuthIndex)

	//Transactions and account management
	router.HandlerFunc(http.MethodPost, "/v1/api/deposit", app.PaymentDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/api/credit", app.PaymentCreditInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/api/debit", app.PaymentDebitInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/api/fullAccessTransfer", app.FullAccessTransferInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/api/fullAccessDeposit", app.FullAccessDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/api/balanceEnquiry", app.BalanceEnquiry)
	router.HandlerFunc(http.MethodPost, "/v1/api/accountHistory", app.AccountHistory)
	router.HandlerFunc(http.MethodGet, "/v1/api/allTransactions", app.AllTransactions)
	router.HandlerFunc(http.MethodGet, "/v1/api/pdfTransactions", app.PdfTransactions)
	router.HandlerFunc(http.MethodGet, "/v1/api/excelTransactions", app.ExcelTransactions)
	router.HandlerFunc(http.MethodPost, "/v1/api/proofOfAddress", app.ProofOfAddress)
	router.HandlerFunc(http.MethodPost, "/v1/api/cashPickup", app.CashPickup)

	//ACCOUNT V2
	router.HandlerFunc(http.MethodPost, "/v1/api/accounts/create", app.AccountCreate)
	router.HandlerFunc(http.MethodPost, "/v1/api/accounts/update", app.AccountUpdate)
	router.HandlerFunc(http.MethodGet, "/v1/api/accounts", app.AccountGet)
	router.HandlerFunc(http.MethodPost, "/v1/api/accounts/block", app.BlockAccount)
	router.HandlerFunc(http.MethodPost, "/v1/api/accounts/unblock", app.UnblockAccount)
	router.HandlerFunc(http.MethodPost, "/v1/api/beneficiary/new", app.NewBeneficiary)
	router.HandlerFunc(http.MethodPost, "/v1/api/beneficiary", app.GetBeneficiaries)

	router.HandlerFunc(http.MethodPost,"v1/api/role_based",app.PaymentCreditInitiation2)

	//Currency Exchange
	router.HandlerFunc(http.MethodGet, "/v1/availableCurrencies", app.AvailableCurrenciesHandler)

	return router
}
