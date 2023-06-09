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
	} else {
		logger := logger.NewLogger("DB_INFO: ")
		logger.Info("¡Database Connected!")
	}
	// Check if the table exists in the database
	rows, err := db.Query("SHOW TABLES LIKE 'stori_transactions'")
	if err != nil {
		logger := logger.NewLogger("DB_ERROR: ")
		logger.Error(fmt.Sprintf("Error checking if the table stori_transactions exists: %s", err.Error()))
	}

	tableExists := false
	if rows.Next() {
		tableExists = true
	}

	// Create the table if it does not exist
	if !tableExists {
		err = createTable(db)
		if err != nil {
			logger := logger.NewLogger("DB_ERROR: ")
			logger.Error(fmt.Sprintf("Error creating the table stori_transactions: %s", err.Error()))
		} else {
			logger := logger.NewLogger("DB_INFO: ")
			logger.Info("The stori_transactions table has been created")
		}
	}

	return &Database{db}, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS stori_transactions (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		file VARCHAR(100),
		idFile INT,
		transaction FLOAT,
		date VARCHAR(100)
	)`)
	if err != nil {
		return err
	}
	return nil
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
		} else {
			logger := logger.NewLogger("DB_INFO: ")
			logger.Info("Value inserted on stori_transactions table")
		}
	}

	return nil
}
