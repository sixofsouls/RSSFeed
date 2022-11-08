package main

import (
	"NewsFeed/pkg/api"
	"NewsFeed/pkg/db"
	"NewsFeed/pkg/parser"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	Rss            []string `json:"rss"`
	Request_period int      `json:"request_period"`
}

func main() {
	//Creating log file
	logfile, err := os.OpenFile("./errorsLog", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}
	//Close errorsLog file
	defer func(logfile *os.File) {
		err := logfile.Close()
		if err != nil {
			log.Println(err)
		}
	}(logfile)
	//Storing all log messages in that file, so server can run in the background
	//log.SetOutput(logfile)

	//Connect to the database
	database, err := db.ConnectToPostgres(os.Getenv("connString"))
	API := api.New(database)

	//Read configuration from config.json
	configFile, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Unable to read config file,", err)
	}
	config := Config{}
	_ = json.Unmarshal([]byte(configFile), &config)

	// Launching server
	go func() {
		log.Println(time.Now().Format(time.RFC1123), " - Server start")
		err := http.ListenAndServe(":80", API.Router())
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Scanning for errors
	go func() {
		for {
			select {
			case gotErr := <-db.ErrCh:
				log.Printf("\nError: %v", gotErr)
			}
		}
	}()

	//Launching goroutines for each RSS feed and setting up update interval
	var wg sync.WaitGroup
	for {
		for _, value := range config.Rss {
			wg.Add(1)
			value := value
			go func(val string) {
				defer wg.Done()
				log.Println(time.Now().Format(time.RFC1123), " - Reading RSS Feed:", value)
				parse, errCh := parser.RRSParser(value)
				if errCh != nil {
					log.Fatal(err)
				}
				err := database.WriteData(parse)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(time.Now().Format(time.RFC1123), " - Database updated for RSS Feed:", value)

			}(value)
		}
		wg.Wait()
		time.Sleep(time.Minute * time.Duration(config.Request_period))
	}
}
