package main

import (
	"net/http"
)

// RenderLoginPage renders a html page
func (app *application) RenderIndexPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/index.html"}, nil, "cmd/web/views/index.html", nil)

}

// RenderIndexPage renders a html page
func (app *application) RenderLoginPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/login.html"}, nil, "cmd/web/views/login.html", nil)

}

// RenderSignUpPage renders a html page
func (app *application) RenderSignUpPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/signup.html"}, nil, "cmd/web/views/signup.html", nil)

}

// RenderSignUpPage renders a html page
func (app *application) RenderCreateAccountPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/account_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)

}
func (app *application) RenderDepositInitiationPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/deposit_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)

}
func (app *application) RenderCreditInitiationPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/credit_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)

}
func (app *application) RenderBatchTransactionPage(w http.ResponseWriter, r *http.Request) {

	app.RenderTemplate(w, []string{"cmd/web/views/account_template.html", "cmd/web/views/layout.html"}, nil, "cmd/web/views/layout.html", nil)

}
