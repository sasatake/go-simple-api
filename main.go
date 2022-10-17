package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/db/ping", Ping)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
