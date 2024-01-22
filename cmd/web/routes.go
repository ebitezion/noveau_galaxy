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
	router.HandlerFunc(http.MethodGet, "/v1/signuppage", (app.RenderSignUpPage))
	router.HandlerFunc(http.MethodGet, "/v1/index", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderIndexPage)))
	router.HandlerFunc(http.MethodGet, "/v1/createAccountPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCreateAccountPage)))
	router.HandlerFunc(http.MethodGet, "/v1/blockAccountPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBlockAccountPage)))
	router.HandlerFunc(http.MethodGet, "/v1/unblockAccountPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderUnblockAccountPage)))
	router.HandlerFunc(http.MethodGet, "/v1/balanceEnquiryPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBalanceEnquiry)))
	router.HandlerFunc(http.MethodGet, "/v1/accountHistoryPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderAccountHistory)))
	router.HandlerFunc(http.MethodGet, "/v1/depositPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderDepositInitiationPage)))
	router.HandlerFunc(http.MethodGet, "/v1/creditPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCreditInitiationPage)))
	router.HandlerFunc(http.MethodGet, "/v1/batchTransactionPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBatchTransactionPage)))
	router.HandlerFunc(http.MethodGet, "/v1/allAccountsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderAllAccountsPage)))
	router.HandlerFunc(http.MethodGet, "/v1/allTransactionsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderTransactionsPage)))
	router.HandlerFunc(http.MethodGet, "/v1/inflowPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderInflowPage)))
	router.HandlerFunc(http.MethodGet, "/v1/outflowPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderOutflowPage)))
	router.HandlerFunc(http.MethodGet, "/v1/businessPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBusinessesPage)))
	router.HandlerFunc(http.MethodGet, "/v1/partnersPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderPartnersPage)))
	router.HandlerFunc(http.MethodGet, "/v1/kycPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderKycPage)))
	router.HandlerFunc(http.MethodGet, "/v1/currencyConversionPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCurrencyConversionPage)))
	router.HandlerFunc(http.MethodGet, "/v1/teamsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderTeamsPage)))
	router.HandlerFunc(http.MethodGet, "/v1/rolesPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderRolesPage)))
	router.HandlerFunc(http.MethodGet, "/v1/systemLogsPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderSystemLogsPage)))
	router.HandlerFunc(http.MethodGet, "/v1/usAccountPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderUsAccountPage)))
	router.HandlerFunc(http.MethodGet, "/v1/ukAccountPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderUkAccountPage)))
	router.HandlerFunc(http.MethodGet, "/v1/cashPickupPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderCashPickupPage)))
	router.HandlerFunc(http.MethodGet, "/v1/allCashPickupPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderAllCashPickupPage)))
	router.HandlerFunc(http.MethodGet, "/v1/userCashPickupPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderUserCashPickupPage)))
	router.HandlerFunc(http.MethodGet, "/v1/withdrawalPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderWithdrawalPage)))
	router.HandlerFunc(http.MethodGet, "/v1/transferApprovalPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderApproveTransferPage)))
	router.HandlerFunc(http.MethodGet, "/v1/withdrawalApprovalPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderApproveWithdrawalPage)))
	router.HandlerFunc(http.MethodGet, "/v1/depositApprovalPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderApproveDepositPage)))
	router.HandlerFunc(http.MethodGet, "/v1/cashPickupApprovalPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderApproveCashPickupPage)))
	router.HandlerFunc(http.MethodGet, "/v1/beneficiariesPage", app.AuthenticationMiddleware(http.HandlerFunc(app.RenderBeneficiariesPage)))

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

	router.HandlerFunc(http.MethodGet, "/v1/accounts", app.AccountGet)
	router.HandlerFunc(http.MethodPost, "/v1/login", app.AuthLogin)
	router.HandlerFunc(http.MethodPost, "/v1/create", app.AuthCreate)
	router.HandlerFunc(http.MethodPost, "/v1/authindex", app.AuthIndex)
	router.HandlerFunc(http.MethodGet, "/v1/signout", app.StaffSignOutProcess)
	router.HandlerFunc(http.MethodPost, "/v1/deposit", app.PaymentDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/credit", app.PaymentCreditInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/fullAccessTransfer", app.FullAccessTransferInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/fullAccessWithdrawal", app.FullAccessWithdrawalInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/fullAccessDeposit", app.FullAccessDepositInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/debit", app.PaymentDebitInitiation)
	router.HandlerFunc(http.MethodPost, "/v1/balanceEnquiry", app.BalanceEnquiry)
	router.HandlerFunc(http.MethodPost, "/v1/accountHistory", app.AccountHistory)
	router.HandlerFunc(http.MethodGet, "/v1/allTransactions", app.AllTransactions)
	router.HandlerFunc(http.MethodPost, "/v1/cashPickup", app.CashPickup)
	// router.HandlerFunc(http.MethodGet, "/v1/cashPickup/all", app.AllCashPickup)
	// router.HandlerFunc(http.MethodPost, "/v1/cashPickup/user", app.UserCashPickup)

	//@TODO i have to update the frontend to call the backend then it should be able to download
	//the updated pdf /  excel sheet
	router.HandlerFunc(http.MethodGet, "/v1/pdfTransactions", app.PdfTransactions)
	router.HandlerFunc(http.MethodGet, "/v1/excelTransactions", app.ExcelTransactions)

	//ACCOUNT V2
	router.HandlerFunc(http.MethodPost, "/v1/accounts/create", app.AccountCreate)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/create/uk", app.AccountCreateUs)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/create/us", app.AccountCreateUk)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/special", app.AccountCreateSpecial)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/update", app.AccountUpdate)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/block", app.BlockAccount)
	router.HandlerFunc(http.MethodPost, "/v1/accounts/unblock", app.UnblockAccount)
	router.HandlerFunc(http.MethodPost, "/v1/beneficiary", app.GetBeneficiaries)

	//Currency Exchange
	router.HandlerFunc(http.MethodGet, "/v1/availableCurrencies", app.AvailableCurrenciesHandler)

	//token
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	//router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	//authorize our API with this
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requirePermission("account:read", app.healthcheckHandler))

	return router
}
