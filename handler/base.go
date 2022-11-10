package handler

import (
	"encoding/json"
	"fmt"
)

const uri = "mongodb://mongo:mongo@localhost:27017/?maxPoolSize=20&w=majority"
const databaseName = "main"
const userCollection = "user"

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func createSimpleResponseJson(status int, message string) string {
	response := Response{Status: status, Message: message}
	responseJsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(responseJsonBytes)
}
