package parser

import (
	"NewsFeed/pkg/db"
	"encoding/xml"
	"log"
	"net/http"
	"time"
)

var (
	amazonShittyLayout = "Mon, 2 Jan 2006 15:04:05 -0700"
)

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Cont  string `xml:"description"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Content string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

type Rss struct {
	Channel Channel `xml:"channel"`
}

// RRSParser takes RSS link and assign data to structure fields
// Then returns channel of Post structure and if any, error
func RRSParser(address string) ([]db.Post, chan error) {

	response, err := http.Get(address)
	if err != nil {
		db.ErrCh <- err
		log.Printf("Could not get RSS feed: %v\n", err)
		return nil, db.ErrCh
	}
	rss := Rss{}
	defer response.Body.Close()

	decoder := xml.NewDecoder(response.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		log.Printf("Error Decode: %v\n", err)
		db.ErrCh <- err
		return nil, db.ErrCh
	}
	var data Item
	posts := make([]db.Post, 0)
	for _, item := range rss.Channel.Items {
		data = item
		Post := db.Post{}
		//assigning values parsed from XML document to structure db.Post
		Post.Title = data.Title
		Post.Content = data.Content
		Post.Link = data.Link
		unixTime, err := time.Parse(time.RFC1123, data.PubDate) //Most used time standard
		if err != nil {
			unixTime, err = time.Parse(time.RFC1123Z, data.PubDate) //Same as above, but with numeric time zone
			if err != nil {
				unixTime, err = time.Parse(amazonShittyLayout, data.PubDate) //Shitty layout used by amazon and other weirdos
				if err != nil {
					log.Fatalf("Could not parse time: %v", err)
				}
			}
		}
		Post.PubTime = unixTime.Unix()
		posts = append(posts, Post)
	}
	//fmt.Println("amount of items in this rss feed:", len(posts))
	return posts, nil
}
