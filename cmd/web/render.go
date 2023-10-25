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
