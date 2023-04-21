package database

import (
	"database/sql"
	"fmt"

	"/Users/ernestoguevara/Desktop/StoriChallenge/cmd/app/config"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDB() (*Database, error) {
	// Load the configuration values
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Use sql.Open to initialize a new sql.DB object
	dataBaseString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.HOST, cfg.Database.Port, cfg.Database.DBName)
	db, err := sql.Open("mysql", dataBaseString)
	if err != nil {
		return nil, err
	}

	// Call db.Ping() to check the connection
	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	fmt.Println("Connected!")
	return &Database{db}, nil
}
