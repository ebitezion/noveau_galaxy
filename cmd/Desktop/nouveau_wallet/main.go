package main

/*

 */

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/index", func(http.ResponseWriter, *http.Request) {

	})

	log.Print("Listening on :5050...")
	err := http.ListenAndServe(":5050", nil)
	if err != nil {
		log.Fatal(err)
	}
}
