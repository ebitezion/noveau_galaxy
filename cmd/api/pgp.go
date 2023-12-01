package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ebitezion/backend-framework/internal/middleware"
)

const pubkey = `-----BEGIN PGP PUBLIC KEY BLOCK-----
...
-----END PGP PUBLIC KEY BLOCK-----`

func Encrypt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Encrypt the JSON body
		armored, err := middleware.EncryptBinaryMessageArmored(pubkey, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Replace the request body with the encrypted message
		r.Body = ioutil.NopCloser(strings.NewReader(string(armored)))

		next.ServeHTTP(w, r)
	})
}
