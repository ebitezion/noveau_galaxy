package main

import (
	"fmt"
	"net/http"

	"github.com/ebitezion/backend-framework/internal/accounts"
)

type BalanceEnquiryPageData struct {
	Data *accounts.BalanceEnquiry
	name string
}
type AllAccountPageData struct {
	Accounts []accounts.AccountDetails
}
type AllTransactionsPageData struct {
	Transactions []accounts.Transaction
}

// RenderIndexPage renders a HTML page
func (app *application) RenderIndexPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/index.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/index.html", nil)
}

// RenderLoginPage renders a HTML page
func (app *application) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/login.html"}, nil, "cmd/web/views/login.html", nil)
}

// RenderSignUpPage renders a HTML page
func (app *application) RenderSignUpPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/signup.html"}, nil, "cmd/web/views/signup.html", nil)
}

// RenderCreateAccountPage renders a HTML page
func (app *application) RenderCreateAccountPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/createAccount.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/createAccount.html", nil)
}

// RenderDepositInitiationPage renders a HTML page
func (app *application) RenderDepositInitiationPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/deposit.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/deposit.html", nil)
}

// RenderCreditInitiationPage renders a HTML page
func (app *application) RenderCreditInitiationPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/credit.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/credit.html", nil)
}

// RenderBatchTransactionPage renders a HTML page
func (app *application) RenderBatchTransactionPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/batch_transaction.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)
}

func (app *application) RenderBalanceEnquiry(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/balanceEnquiry.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/balanceEnquiry.html", nil)
}

func (app *application) RenderAccountHistory(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/accountHistory.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/accountHistory.html", nil)
}

func (app *application) RenderBusinessesPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/businesses.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/businesses.html", nil)
}
func (app *application) RenderPartnersPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/partners.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/partners.html", nil)
}
func (app *application) RenderKycPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/kyc.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/kyc.html", nil)
}
func (app *application) RenderCurrencyConversionPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/currencyConverter.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/currencyConverter.html", nil)
}
func (app *application) RenderTeamsPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/teams.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/teams.html", nil)
}
func (app *application) RenderRolesPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/roles.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/roles.html", nil)
}
func (app *application) RenderSystemLogsPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/systemLogs.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, nil, "cmd/web/views/systemLogs.html", nil)
}

func (app *application) RenderTransactionsPage(w http.ResponseWriter, r *http.Request) {
	//get all transactions
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1008"})
	if err != nil {
		fmt.Println(err)
	}
	Transactions, ok := data.([]accounts.Transaction)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}

	pageData := AllTransactionsPageData{
		Transactions: Transactions,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/allTransactions.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/allTransactions.html", nil)
}

// RenderBatchTransactionPage renders a HTML page
func (app *application) RenderAllAccountsPage(w http.ResponseWriter, r *http.Request) {
	data, err := accounts.ProcessAccount([]string{"", "acmt", "1000"})
	if err != nil {
		fmt.Println(err)
		return
	}
	accountDetails, ok := data.([]accounts.AccountDetails)
	if !ok {
		fmt.Println("Failed to convert to []AccountDetails")
		return
	}

	pageData := AllAccountPageData{
		Accounts: accountDetails,
	}
	app.RenderTemplate(w, []string{"cmd/web/views/allAccounts.html", "cmd/web/views/header.html", "cmd/web/views/footer.html"}, pageData, "cmd/web/views/allAccounts.html", nil)
}
