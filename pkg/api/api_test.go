package api

import (
	"NewsFeed/pkg/db"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAPI_newsHandler(t *testing.T) {
	database, _ := db.ConnectToPostgres(os.Getenv("connString"))
	api := New(database)
	// New request. Want to get 5 news
	request := httptest.NewRequest(http.MethodGet, "/news/5", nil)
	recorder := httptest.NewRecorder()
	//Sending request. Reading response.
	api.r.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Errorf("Invalid response. Wanted %d, got %d", recorder.Code, http.StatusOK)
	}
	//Reading server response data
	body, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Unable to read server response: %v", err)
	}
	//Decoding server response data
	var data []db.Post
	err = json.Unmarshal(body, &data)
	if err != nil {
		t.Errorf("Unable to decode server response: %v", err)
	}
	//Checking amount of news returned
	if len(data) != 5 {
		t.Errorf("Invalid amount of news returned. Wanted %v, got %v", 5, len(data))
	}
}
