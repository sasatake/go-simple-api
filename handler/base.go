package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const uri = "mongodb://mongo:mongo@localhost:27017/?maxPoolSize=20&w=majority"
const databaseName = "main"
const userCollection = "user"

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

func responseListUsers(w http.ResponseWriter, users []User) {
	response := UsersResponse{Response: Response{Status: http.StatusOK, Message: "list users"}, Users: users}
	responseJsonBytes, _ := json.MarshalIndent(response, "", "  ")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(responseJsonBytes))
}
