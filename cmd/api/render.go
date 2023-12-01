package main

import (
	"net/http"
)

// RenderIndexPage renders a HTML page
func (app *application) RenderIndexPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/index.html"}, nil, "cmd/web/views/index.html", nil)
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
	app.RenderTemplate(w, []string{"cmd/web/views/account_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)
}

// RenderDepositInitiationPage renders a HTML page
func (app *application) RenderDepositInitiationPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/deposit_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)
}

// RenderCreditInitiationPage renders a HTML page
func (app *application) RenderCreditInitiationPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/credit_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)
}

// RenderBatchTransactionPage renders a HTML page
func (app *application) RenderBatchTransactionPage(w http.ResponseWriter, r *http.Request) {
	app.RenderTemplate(w, []string{"cmd/web/views/batch_transaction.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)
}
