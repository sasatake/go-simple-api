package handler

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.URL.Path != "/" {
		response(w, http.StatusNotFound, "not found.")
		return
	}
	name := r.URL.Query().Get("name")

	if name == "" {
		response(w, http.StatusBadRequest, "set name parameter.")
		return
	}

	message := fmt.Sprintf("Hello %s", name)
	response(w, http.StatusOK, message)
}

func Ping(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			response(w, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	response(w, http.StatusOK, "Successfully connected and pinged.")
}
