package database

import (
	"database/sql"
	"fmt"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/config"
	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	db *sql.DB
}

func NewDB() (*Database, error) {
	// Load the configuration values
	cfg, err := config.LoadConfig()
	if err != nil {
		logger := logger.NewLogger("CONFIG_ERROR: ")
		logger.Error(fmt.Sprintf("Error loading configuration values: %v", err))

	}

	// Use sql.Open to initialize a new sql.DB object
	dataBaseString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.Database.User, cfg.Database.Password, cfg.Database.HOST, cfg.Database.Port, cfg.Database.DBName)
	db, err := sql.Open("mysql", dataBaseString)
	if err != nil {
		logger := logger.NewLogger("DB_ERROR: ")
		logger.Error(fmt.Sprintf("Error initialazing database: %s", err.Error()))

	}

	// Call db.Ping() to check the connection
	pingErr := db.Ping()
	if pingErr != nil {
		logger := logger.NewLogger("DB_ERROR: ")
		logger.Error(fmt.Sprintf("Error to ping database: %s", pingErr.Error()))
	}
	logger := logger.NewLogger("DB_INFO: ")
	logger.Info("¡Database Connected!")
	return &Database{db}, nil
}

// InsertTransaction inserts a transaction into the database
func (d *Database) InsertTransaction(file string, idFile int, transaction float64, date string) error {
	// Check if the row already exists in the table
	var count int
	err := d.db.QueryRow("SELECT COUNT(*) FROM stori_transactions WHERE file = ? AND idFile = ?", file, idFile).Scan(&count)
	if err != nil {
		logger := logger.NewLogger("DB_ERROR: ")
		logger.Error(fmt.Sprintf("Error checking if the row already exists in the table stori_transactions: %s", err.Error()))
	}

	// Insert the row if it doesn't already exist in the table
	if count == 0 {
		_, err = d.db.Exec("INSERT INTO stori_transactions(file,idFile,transaction,date) VALUES (?,?, ?, ?)", file, idFile, transaction, date)
		if err != nil {
			logger := logger.NewLogger("DB_ERROR: ")
			logger.Error(fmt.Sprintf("Error inserting the row in the table stori_transactions: %s", err.Error()))
		}
	}

	return nil
}
