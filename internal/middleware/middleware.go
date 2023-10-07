package middleware

import (
	"net/http"
	"strings"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	// Define the allowed IP address
	allowedIP := "192.168.1.100"

	// Get the remote address from the request
	remoteAddr := strings.Split(r.RemoteAddr, ":")[0]

	// Check if the remote address matches the allowed IP
	if remoteAddr != allowedIP {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Your handler logic here
	w.Write([]byte("Welcome!"))
}
