package main

import (
	"log"
	"net/http"

	"github.com/sasatake/go-simple-api/handler"
)

func main() {
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/db/ping", handler.Ping)
	http.HandleFunc("/users", handler.ListUser)
	http.HandleFunc("/user", handler.RegisterUser)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
