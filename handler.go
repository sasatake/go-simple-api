package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func createResponseJson(status int, message string) string {
	response := Response{Status: status, Message: message}
	responseJsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(responseJsonBytes)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, createResponseJson(http.StatusNotFound, "not found."))
		return
	}
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, createResponseJson(http.StatusBadRequest, "set name parameter."))
		return
	}

	message := fmt.Sprintf("Hello %s", name)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, createResponseJson(http.StatusOK, message))
}
