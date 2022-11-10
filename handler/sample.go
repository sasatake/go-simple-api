package handler

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusNotFound, "not found."))
		return
	}
	name := r.URL.Query().Get("name")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusBadRequest, "set name parameter."))
		return
	}

	message := fmt.Sprintf("Hello %s", name)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, createSimpleResponseJson(http.StatusOK, message))
}

func Ping(w http.ResponseWriter, r *http.Request) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusInternalServerError, err.Error()))
		return
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, createSimpleResponseJson(http.StatusInternalServerError, err.Error()))
			return
		}
	}()

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusInternalServerError, err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, createSimpleResponseJson(http.StatusOK, "Successfully connected and pinged."))
}
