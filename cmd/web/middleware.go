package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gorilla/sessions"
)

var SECRETKEY = []byte(os.Getenv("JWT_KEY"))

func (app *application) generateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SECRETKEY)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (app *application) AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "JwtToken")
		if err != nil {
			log.Println("Error getting  session-----Unauthorized Access")
			http.Redirect(w, r, "/v1/loginpage", http.StatusSeeOther)
			return
		}
		token, ok := session.Values["token"].(string)
		if !ok || token == "" {
			log.Println("Token is missing or invalid-----Unauthorized Access")
			http.Redirect(w, r, "/v1/loginpage", http.StatusSeeOther)
			return
		}
		// validate jwt
		ParsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method and return the secret key used to sign the token
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return SECRETKEY, nil

		})
		if err != nil || !ParsedToken.Valid {
			log.Println("Token is invalid-----Unauthorized Access")
			http.Redirect(w, r, "/v1/loginpage", http.StatusSeeOther)
			return
		}
		// Token is valid, continue with the next handler
		next.ServeHTTP(w, r)

	})

}

func (app *application) SetJwtSession(w http.ResponseWriter, r *http.Request, accountNumber string, fullname string) error {

	token, err := app.generateToken(accountNumber)
	if err != nil {
		log.Println("Error generating token:", err)
		return fmt.Errorf("internal Server Error")
	}

	// Set the token in an HTTP-only cookie
	//@TODO change it to being stored on redis
	session, err := store.Get(r, "JwtToken")
	if err != nil {
		log.Println("Error creating session ")
		return fmt.Errorf("internal Server Error")

	}

	session.Values["token"] = token
	session.Values["fullname"] = fullname
	session.Options = &sessions.Options{
		MaxAge:   24 * 3600,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	}

	err = session.Save(r, w)
	if err != nil {
		log.Println("Error saving session")
		return fmt.Errorf("internal Server Error")

	}
	return nil
}
