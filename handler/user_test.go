package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func resetUserCollection(t *testing.T) {
	t.Helper()
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	coll := client.Database(databaseName).Collection(userCollection)
	if err := coll.Drop(context.Background()); err != nil {
		t.Fatalf("drop pre-test: %v", err)
	}
	t.Cleanup(func() {
		_ = coll.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	})
}

func TestRegisterUser(t *testing.T) {
	resetUserCollection(t)

	body := bytes.NewBufferString(`{"name":"akira","nickname":"sasatake","mail":"a@example.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/user", body)
	rec := httptest.NewRecorder()
	RegisterUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
}

func TestRegisterUser_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(`{not json`))
	rec := httptest.NewRecorder()
	RegisterUser(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestListUser_Empty(t *testing.T) {
	resetUserCollection(t)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()
	ListUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp UsersResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Users) != 0 {
		t.Errorf("got %d users, want 0", len(resp.Users))
	}
}

func TestListUser_AfterRegister(t *testing.T) {
	resetUserCollection(t)

	regReq := httptest.NewRequest(http.MethodPost, "/user",
		bytes.NewBufferString(`{"name":"akira","nickname":"sasatake","mail":"a@example.com"}`))
	regRec := httptest.NewRecorder()
	RegisterUser(regRec, regReq)
	if regRec.Code != http.StatusCreated {
		t.Fatalf("register failed: %d %s", regRec.Code, regRec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/users", nil)
	listRec := httptest.NewRecorder()
	ListUser(listRec, listReq)

	var resp UsersResponse
	if err := json.Unmarshal(listRec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Users) != 1 {
		t.Fatalf("got %d users, want 1", len(resp.Users))
	}
	if resp.Users[0].Name != "akira" {
		t.Errorf("name = %q, want akira", resp.Users[0].Name)
	}
	if resp.Users[0].Nickname != "sasatake" {
		t.Errorf("nickname = %q, want sasatake", resp.Users[0].Nickname)
	}
	if resp.Users[0].Mail != "a@example.com" {
		t.Errorf("mail = %q, want a@example.com", resp.Users[0].Mail)
	}
}
