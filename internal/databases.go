package internal

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Game struct {
	Id                int      
	Game_name         string
	Author_corp       string
	Images_dir        string   
	Images_url        []string 
	Game_post_content string
	Date_realease     string
	Waitings          int
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

// Add Game to db
func AddGame(db *sql.DB, game *Game) error {
	_, err := db.Exec(`INSERT INTO games (game_name, author_corp, images_dir, game_post_content, date_realease) 
	VALUES ($1, $2, $3, $4, $5)`, game.Game_name, game.Author_corp, game.Images_dir, game.Game_post_content, game.Date_realease)
	if err != nil {
		return fmt.Errorf("db INSERT problem: %s", err)
	}
	err = os.MkdirAll(game.Images_dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("images_dir error: %s", err)
	}
	return nil
}

// Gives a specific Game by id
func GetSomeGame(db *sql.DB, game_id int) (*Game, error) {
	rows, err := db.Query(`SELECT g.*, COUNT(w.user_id) AS waitings_count FROM Games AS g 
	LEFT JOIN Waitings AS w ON g.id = w.game_id WHERE id=$1 GROUP BY g.id;`, game_id)
	if err != nil {
		return nil, fmt.Errorf("can't get Game with id: %d -- %s", game_id, err)
	}
	defer rows.Close()
	game := &Game{}
	for rows.Next() {
		err = rows.Scan(&game.Id, &game.Game_name, &game.Author_corp, &game.Images_dir, &game.Game_post_content, &game.Date_realease, &game.Waitings)
		if err != nil {
			return nil, fmt.Errorf("error in Scan for params, id: %d -- %s", game_id, err)
		}
		images, err := os.ReadDir(fmt.Sprintf("./static/images/Games/%s", game.Images_dir))
		if err != nil {
			return nil, fmt.Errorf("error in Open image dir params %s", err)
		}
		for _, img := range images {
			game.Images_url = append(game.Images_url, img.Name())
		}
	}
	return game, nil
}

// Return lim Games with the most waitings
func GetGameWMW(db *sql.DB, lim int) ([]Game, error) {
	var games []Game
	rows, err := db.Query(`SELECT g.*, COUNT(w.user_id) AS waitings_count FROM Games AS g 
	LEFT JOIN Waitings AS w ON g.id = w.game_id GROUP BY g.id LIMIT $1;`, lim)
	if err != nil {
		return nil, fmt.Errorf("can't get %d Posts, Server error: %s", lim, err)
	}
	defer rows.Close()
	for rows.Next() {
		game := Game{}
		err = rows.Scan(&game.Id, &game.Game_name, &game.Author_corp, &game.Images_dir, &game.Game_post_content, &game.Date_realease, &game.Waitings)
		if err != nil {
			return nil, fmt.Errorf("error in Scan params %s", err)
		}
		images, err := os.ReadDir(fmt.Sprintf("./static/images/Games/%s", game.Images_dir))
		if err != nil {
			return nil, fmt.Errorf("error in Open image dir params %s", err)
		}
		for _, img := range images {
			game.Images_url = append(game.Images_url, img.Name())
		}
		games = append(games, game)
	}
	rows.Close()
	return games, nil
}

// Delete game by id
func DeleteGame(db *sql.DB, game_id int) error {
	_, err := db.Exec(`DELETE FROM games WHERE id = $1`, game_id)
	return err
}

// Incoming game has all edits
func EditGame(db *sql.DB, game *Game) error {
	_, err := db.Exec(`UPDATE games SET game_name=$1, author_corp=$2, images_dir=$3, game_post_content=$4, date_realease=$5 
	WHERE id=$6;`, game.Game_name, game.Author_corp, game.Images_dir, game.Game_post_content, game.Date_realease, game.Id)
	return err
}

// Add +1 waiter to game
func AddWaiter(db *sql.DB, game_id int, user_id int) error {
	_, err := db.Exec(`INSERT INTO waitings (user_id, game_id) VALUES ($1, $2)`, user_id, game_id)
	return err
}

// Take away waiter from game
func DeleteWaiter(db *sql.DB, game_id int, user_id int) error {
	_, err := db.Exec(`DELETE FROM waitings WHERE user_id=$1 AND game_id=$2`, user_id, game_id)
	return err
}

// Add User to db
func AddUser(db *sql.DB, usr *User) error {
	usr.Pass_hash = base64.RawStdEncoding.EncodeToString([]byte(usr.Password))
	_, err := db.Exec(`INSERT INTO users (username, email, pass_hash) VALUES ($1, $2, $3)`, usr.Username, usr.Email, usr.Pass_hash)
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
	rows, err := db.Query(`SELECT pass_hash FROM users WHERE username=$1`, usr.Username)
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
