// package tests

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/trodix/go-rest-api/api/handlers"
// 	"github.com/trodix/go-rest-api/models"
// )

// func TestCreateUser(t *testing.T) {
// 	var jsonStr = []byte(`{"username":"john","email":"john@example.com"}`)

// 	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(handlers.CreateUser(nil))

// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
// 	}

// 	var user models.User
// 	err = json.NewDecoder(rr.Body).Decode(&user)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if user.Username != "john" {
// 		t.Errorf("expected username 'john', got %v", user.Username)
// 	}
// }
