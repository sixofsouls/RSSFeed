package db

import (
	"os"
	"testing"
)

var (
	database, err = ConnectToPostgres(os.Getenv("connString"))
)

func TestDB_WriteData(t *testing.T) {
	//checking if we could connect to the database
	if err != nil {
		t.Errorf("Could not connec to database: %v", err)
	}
	post := []Post{{
		Title:   "Test Title for pkgTest",
		Link:    "test.link",
		Content: "Test description",
		PubTime: 60000,
	}}
	err := database.WriteData(post)
	if err != nil {
		t.Errorf("Unable to write data to database: %v", err)
	}
}

func TestDB_ReadPosts(t *testing.T) {
	//reading array of posts
	amountToShow := 10
	ReadPost, err := database.ReadPosts(amountToShow)
	//readCh := make(chan []Post, 255)
	if err != nil {
		t.Errorf("Failed to read posts: %v", err)
	}
	if len(ReadPost) != amountToShow {
		t.Errorf("Incorrect amount of rows shown. Got %v, wanted %v", len(ReadPost), amountToShow)
	}
}

func TestDB_DeletePost(t *testing.T) {
	//Creating pseudo post with link value same as in TestDB_WriteData test
	postToDelete := Post{Link: "test.link"}
	err = database.DeletePost(postToDelete)
	if err != nil {
		t.Errorf("Could not delete post: %v", err)
	}
}
