package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersResponse struct {
	Response Response `json:"response"`
	Users    []User   `json:"users"`
}

type User struct {
	Id       string `json:"id" bson:"_id"`
	Name     string `json:"name" bson:"name"`
	Nickname string `json:"nickname" bson:"nickname"`
	Mail     string `json:"mail" bson:"mail"`
}

func createListUsersResponseJson(users []User) string {
	response := UsersResponse{Response: Response{Status: http.StatusOK, Message: "list users"}, Users: users}
	responseJsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(responseJsonBytes)
}

func ListUser(w http.ResponseWriter, r *http.Request) {

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

	collection := client.Database(databaseName).Collection(userCollection)
	cursor, err := collection.Find(context.TODO(), options.Find())
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, createListUsersResponseJson(users))
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusMethodNotAllowed, "method not allowed."))
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusBadRequest, "bad request body."))
		return
	}

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

	collection := client.Database(databaseName).Collection(userCollection)
	doc := bson.D{
		{Key: "name", Value: user.Name},
		{Key: "nickname", Value: user.Nickname},
		{Key: "mail", Value: user.Mail},
	}
	result, err := collection.InsertOne(context.TODO(), doc)

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, createSimpleResponseJson(http.StatusCreated, oid.Hex()))
	}
}
