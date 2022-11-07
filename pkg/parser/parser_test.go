package parser

import (
	"NewsFeed/pkg/db"
	"testing"
)

var (
	database = db.DB{}
	post     = db.Post{}
)

func TestRRSParser(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name string
		args args
		val  chan db.Post
		err  chan error
	}{
		{
			name: "Invalid input",
			args: args{address: "invalid_input"},
			val:  nil,
		},
		{
			name: "Invalid RSS link",
			args: args{address: "http://google.com"},
			val:  nil,
		},
		{
			name: "Correct RSS link. time.RFC1123 format",
			args: args{address: "https://habr.com/ru/rss/hub/go/all/?fl=ru"},
			err:  nil,
		},
		{
			name: "Correct RSS link. time.RFC1123Z format",
			args: args{address: "http://static.feed.rbc.ru/rbc/logical/footer/news.rss"},
			err:  nil,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, goterr := RRSParser(tt.args.address)
			switch {
			case i == 0 && got != nil: // invalid input, expecting error
				t.Errorf("Got post value - %v, wanted %v", got, nil)
			case i == 1 && got != nil: // incorrect link, expecting error
				t.Errorf("Got post value - %v, wanted %v", got, nil)
			case i == 2 && goterr != nil: // correct link time.RFC1123 format, expecting err = nil
				t.Errorf("Got error - %v, wanted %v", goterr, nil)
			case i == 3 && goterr != nil: // correct link time.RFC1123Z format, expecting err = nil
				t.Errorf("Got error - %v, wanted %v", goterr, nil)
			default:
			}
		})
	}
}
