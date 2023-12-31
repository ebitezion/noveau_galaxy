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
			http.Error(w, "Unauthorized Access", http.StatusInternalServerError)
			return
		}
		token, ok := session.Values["token"].(string)
		if !ok || token == "" {
			log.Println("Token is missing or invalid-----Unauthorized Access")
			http.Error(w, "Unauthorized Access", http.StatusInternalServerError)
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
			http.Error(w, "Unauthorized Access", http.StatusInternalServerError)
			return
		}
		// Token is valid, continue with the next handler
		next.ServeHTTP(w, r)

	})

}

func (app *application) SetJwtSession(w http.ResponseWriter, r *http.Request, accountNumber string) error {

	token, err := app.generateToken(accountNumber)
	if err != nil {
		log.Println("Error generating token:", err)
		return fmt.Errorf("internal Server Error")
	}

	// Set the token in an HTTP-only cookie
	session, err := store.Get(r, "JwtToken")
	if err != nil {
		log.Println("Error creating session ")
		return fmt.Errorf("internal Server Error")

	}

	session.Values["token"] = token
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

// func requestLogger(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         start := time.Now()

//         l := logger.Get()

//         defer func() {
//             l.
//                 Info().
//                 Str("method", r.Method).
//                 Str("url", r.URL.RequestURI()).
//                 Str("user_agent", r.UserAgent()).
//                 Dur("elapsed_ms", time.Since(start)).
//                 Msg("incoming request")
//         }()

//         next.ServeHTTP(w, r)
//     })
// }
