package internal

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDb() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("Can`t load data of user")
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
		panic(err)
	}
	return db, nil
}

func AddPost() {

}
