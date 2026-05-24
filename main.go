package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sasatake/go-simple-api/handler"
)

func main() {
	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/db/ping", handler.Ping)
	http.HandleFunc("/users", handler.ListUser)
	http.HandleFunc("/user", userHandler)

	port := "8080"
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		handler.RegisterUser(w, r)
	case http.MethodDelete:
	default:
		response(w, http.StatusMethodNotAllowed, "method not allowed.")
	}
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func response(w http.ResponseWriter, status int, message string) {
	response := Response{Status: status, Message: message}
	responseJsonBytes, _ := json.MarshalIndent(response, "", "  ")
	w.WriteHeader(status)
	fmt.Fprint(w, string(responseJsonBytes))
}
