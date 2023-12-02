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
