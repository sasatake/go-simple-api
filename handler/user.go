package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type UsersResponse struct {
	Response Response `json:"response"`
	Users    []User   `json:"users"`
}

type User struct {
	Id       bson.ObjectID `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Nickname string        `json:"nickname" bson:"nickname"`
	Mail     string        `json:"mail" bson:"mail"`
}

func ListUser(w http.ResponseWriter, r *http.Request) {

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

	collection := client.Database(databaseName).Collection(userCollection)
	cursor, err := collection.Find(context.TODO(), options.Find())
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}
	var users []User
	if err = cursor.All(context.TODO(), &users); err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseListUsers(w, users)
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response(w, http.StatusBadRequest, "bad request body.")
		return
	}

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

	collection := client.Database(databaseName).Collection(userCollection)
	doc := bson.D{
		{Key: "name", Value: user.Name},
		{Key: "nickname", Value: user.Nickname},
		{Key: "mail", Value: user.Mail},
	}
	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		response(w, http.StatusInternalServerError, err.Error())
		return
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		response(w, http.StatusCreated, oid.Hex())
	}
}
