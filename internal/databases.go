package internal

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Post struct {
	Id           int `json:"-"`
	AuthorName   string
	Title        string
	Post_content string
	Image_url    sql.NullString
	Waitings     int
	Created_at   time.Time
	Author_id    int `json:"-"`
}

type User struct {
	Id        int
	Pass_hash string
	Username  string
	Email     string
	Password  string
}

// Conn to db and return it
func GetDb() (*sql.DB, error) {
	err := godotenv.Load("keys.env")
	if err != nil {
		return nil, fmt.Errorf("can`t load data from .env")
	}
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("can`t open database, invalid dsn params")
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ping error: %s", err.Error())
	}
	return db, nil
}

// Add Post to db
func AddPost(db *sql.DB, post *Post) error {
	var err error
	if post.Image_url.Valid {
		_, err = db.Exec(`INSERT INTO "Posts" (title, post_content, image_url, author_id) VALUES ($1, $2, $3, $4)`,
			post.Title, post.Post_content, post.Image_url, post.Author_id)
	} else {
		_, err = db.Exec(`INSERT INTO "Posts" (title, post_content, author_id) VALUES ($1, $2, $3)`,
			post.Title, post.Post_content, post.Author_id)
	}
	if err != nil {
		return fmt.Errorf("db INSERT problem: %s", err)
	}
	return nil
}

// Gives a specific Post by id
func GetPost(db *sql.DB, post_id int) (*Post, error) {
	rows, err := db.Query(`SELECT U.username AS authorname, P.* 
	FROM "Posts" AS P JOIN "Users" AS U ON P.author_id = U.id WHERE P.id=$1`, post_id)
	if err != nil {
		return nil, fmt.Errorf("can't get Post with id: %d -- %s", post_id, err)
	}
	defer rows.Close()
	post := &Post{}
	for rows.Next() {
		err = rows.Scan(&post.AuthorName, &post.Id, &post.Title, &post.Post_content, &post.Image_url, &post.Waitings, &post.Created_at, &post.Author_id)
		if err != nil {
			return nil, fmt.Errorf("error in Scan for params, id: %d -- %s", post_id, err)
		}
	}
	return post, nil
}

// Return lim Posts with the most waitings
func GetPostsWMV(db *sql.DB, lim int) ([]Post, error) {
	var posts []Post
	rows, err := db.Query(`SELECT U.username AS authorname, P.* 
	FROM "Posts" AS P JOIN "Users" AS U ON P.author_id = U.id ORDER BY P.waitings DESC LIMIT $1`, lim)
	if err != nil {
		return nil, fmt.Errorf("can't get %d Posts, Server error: %s", lim, err)
	}
	defer rows.Close()
	post := Post{}
	for rows.Next() {
		err = rows.Scan(&post.AuthorName, &post.Id, &post.Title, &post.Post_content, &post.Image_url, &post.Waitings, &post.Created_at, &post.Author_id)
		posts = append(posts, post)
		if err != nil {
			rows.Close()
			return nil, fmt.Errorf("error in Scan params %s", err)
		}
	}
	rows.Close()
	return posts, nil
}

// Delete post by id
func DeletePost(db *sql.DB, post_id int) error {
	_, err := db.Exec(`DELETE FROM "Posts" WHERE id = $1`, post_id)
	return err
}

// Incoming post has all edits
func EditPost(db *sql.DB, post *Post) error {
	req := `UPDATE "Posts" SET `
	if post.Image_url.Valid {
		req += "image_url=$1, title=$2, post_content=$3 WHERE id=$4"
		_, err := db.Exec(req, post.Image_url, post.Title, post.Post_content, post.Id)
		return err
	}
	req += "title=$1, post_content=$2 WHERE id=$3"
	_, err := db.Exec(req, post.Title, post.Post_content, post.Id)
	return err
}

// Add +1 waiter to post
func AddWaiter(db *sql.DB, post_id int) error {
	_, err := db.Exec(`UPDATE "Posts" SET waitings=waitings+1 WHERE id=$1`, post_id)
	return err
}

// Add to User db
func AddUser(db *sql.DB, usr *User) error {
	usr.Pass_hash = base64.RawStdEncoding.EncodeToString([]byte(usr.Password))
	_, err := db.Exec(`INSERT INTO "Users" (username, email, pass_hash) VALUES ($1, $2, $3)`, usr.Username, usr.Email, usr.Pass_hash)
	if err != nil {
		return fmt.Errorf("that username or email is busy: %s", err)
	}
	return nil
}

// Serach User by username and password:
// return -1 if user dosn't exist or server error =),
// return 0  if password incorrect,
// return 1  if user and pass correct.
func SearchUser(db *sql.DB, usr *User) int8 {
	usr.Pass_hash = base64.RawStdEncoding.EncodeToString([]byte(usr.Password))
	rows, err := db.Query(`SELECT pass_hash FROM "Users" WHERE username=$1`, usr.Username)
	if err != nil {
		return -1
	}
	defer rows.Close()
	var db_pass string
	for rows.Next() {
		err = rows.Scan(&db_pass)
		if err != nil {
			return -1
		}
	}
	if db_pass == "" {
		return -1
	}
	if strings.Compare(db_pass, usr.Pass_hash) != 0 {
		return 0
	}
	return 1
}
