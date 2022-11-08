package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
)

// ErrCh channel for errors across entire application
var ErrCh = make(chan error, 10)

// DB structure for postgres pool
type DB struct {
	m  sync.Mutex
	db *pgxpool.Pool
}

// Post from RSS feed.
type Post struct {
	ID      int
	Title   string
	Content string
	PubTime int64
	Link    string
}

// ConnectToPostgres connects to the database via connection string
func ConnectToPostgres(connString string) (*DB, error) {
	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	database := DB{db: db}
	return &database, nil
}

// WriteData to the database, ignoring links that are already in the database
func (d *DB) WriteData(inc []Post) chan error {
	ctx := context.Background()
	d.m.Lock()
	defer d.m.Unlock()

	tx, err := d.db.Begin(ctx)
	if err != nil {
		ErrCh <- err
		log.Println(err)
		return ErrCh
	}
	defer tx.Rollback(ctx)

	batch := new(pgx.Batch)
	for _, post := range inc {
		batch.Queue(`
		INSERT INTO posts(title, link, content, pubdate)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (link) DO NOTHING;
	`,
			post.Title,
			post.Link,
			post.Content,
			post.PubTime)
	}
	result := tx.SendBatch(ctx, batch)
	err = result.Close()
	if err != nil {
		ErrCh <- err
		log.Println(err)
		return ErrCh
	}
	err = tx.Commit(ctx)
	if err != nil {
		ErrCh <- err
		log.Println(err)
		return ErrCh
	}
	return nil
}

// ReadPosts return channel array of posts stored in the database
func (d *DB) ReadPosts(amountToShow int) ([]Post, error) {
	d.m.Lock()
	defer d.m.Unlock()

	var Posts []Post
	//SQL request. Getting all values from DB
	rows, err := d.db.Query(context.Background(), `
	SELECT id, title, link, content, pubdate FROM posts
	ORDER BY pubdate desc
	LIMIT $1;
	`, &amountToShow)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Link,
			&p.Content,
			&p.PubTime,
		)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		Posts = append(Posts, p)
	}
	return Posts, rows.Err()
}

// DeletePost deletes post by its link
func (d *DB) DeletePost(post Post) error {
	d.m.Lock()
	defer d.m.Unlock()
	_, err := d.db.Exec(context.Background(), `
	DELETE FROM posts
	WHERE posts.link = $1;
	`, &post.Link)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
