package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func Ping(w http.ResponseWriter, r *http.Request) {

	const uri = "mongodb://mongo:mongo@localhost:27017/?maxPoolSize=20&w=majority"

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, createResponseJson(http.StatusInternalServerError, err.Error()))
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, createResponseJson(http.StatusInternalServerError, err.Error()))
			return
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, createResponseJson(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, createResponseJson(http.StatusOK, "Successfully connected and pinged."))
}
