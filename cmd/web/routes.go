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
	router.HandlerFunc(http.MethodGet, "/v1/loginpage", app.RenderLoginPage)
	router.HandlerFunc(http.MethodGet, "/v1/signuppage", app.RenderSignUpPage)
	router.HandlerFunc(http.MethodGet, "/v1/index", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderIndexPage)))
	router.HandlerFunc(http.MethodGet, "/v1/createAccountPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCreateAccountPage)))
	router.HandlerFunc(http.MethodGet, "/v1/balanceEnquiryPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBalanceEnquiry)))
	router.HandlerFunc(http.MethodGet, "/v1/accountHistoryPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderAccountHistory)))
	router.HandlerFunc(http.MethodGet, "/v1/depositPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderDepositInitiationPage)))
	router.HandlerFunc(http.MethodGet, "/v1/creditPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCreditInitiationPage)))
	router.HandlerFunc(http.MethodGet, "/v1/batchTransactionPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBatchTransactionPage)))
	router.HandlerFunc(http.MethodGet, "/v1/allAccountsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderAllAccountsPage)))
	router.HandlerFunc(http.MethodGet, "/v1/allTransactionsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderTransactionsPage)))

	router.HandlerFunc(http.MethodGet, "/v1/businessPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBusinessesPage)))
	router.HandlerFunc(http.MethodGet, "/v1/partnersPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderPartnersPage)))
	router.HandlerFunc(http.MethodGet, "/v1/kycPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderKycPage)))
	router.HandlerFunc(http.MethodGet, "/v1/currencyConversionPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCurrencyConversionPage)))
	router.HandlerFunc(http.MethodGet, "/v1/teamsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderTeamsPage)))
	router.HandlerFunc(http.MethodGet, "/v1/rolesPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderRolesPage)))
	router.HandlerFunc(http.MethodGet, "/v1/systemLogsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderSystemLogsPage)))
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
	router.HandlerFunc(http.MethodGet, "/v1/signout", app.StaffSignOutProcess)
	router.HandlerFunc(http.MethodPost, "/v1/deposit", app.PaymentDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/credit", app.PaymentCreditInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/fullAccessCredit", app.FullAccessCreditInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/fullAccessDeposit", app.FullAccessDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/debit", app.PaymentDebitInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	router.HandlerFunc(http.MethodPost, "/v1/accountHistory", app.AccountHistory)
	router.HandlerFunc(http.MethodGet, "/v1/allTransactions", app.AllTransactions)
	//@TODO i have to update the frontend to call the backend then it should be able to download
	//the updated pdf /  excel sheet
	router.HandlerFunc(http.MethodGet, "/v1/pdfTransactions", app.PdfTransactions)
	router.HandlerFunc(http.MethodGet, "/v1/excelTransactions", app.ExcelTransactions)

	//ACCOUNT V2
	router.HandlerFunc(http.MethodPost, "/v1/accounts/create", app.AccountCreate)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/special", app.AccountCreateSpecial)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/update", app.AccountUpdate)
	//Currency Exchange
	router.HandlerFunc(http.MethodGet, "/v1/availableCurrencies", app.AvailableCurrenciesHandler)

	return router
}
