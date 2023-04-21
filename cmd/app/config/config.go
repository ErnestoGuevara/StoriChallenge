package config

import (
	"fmt"
	"os"

	"github.com/ErnestoGuevara/StoriChallenge/cmd/app/logger"
	"github.com/joho/godotenv"
)

type config struct {
	Database struct {
		User     string
		Password string
		HOST     string
		Port     string
		DBName   string
	}
	SendGrid struct {
		Api      string
		Template string
	}
}

func LoadConfig() (*config, error) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		logger := logger.NewLogger("FILE_ERROR: ")
		logger.Error(fmt.Sprintf("Error loading .env file: %v", err))
	}

	// Create a new config instance
	cfg := &config{}

	// Read the configuration values from the environment variables
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.HOST = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.DBName = os.Getenv("DB_NAME")

	cfg.SendGrid.Api = os.Getenv("SENDGRID_API_KEY")
	cfg.SendGrid.Template = os.Getenv("SENDGRID_TEMPLATE_ID")

	return cfg, nil
}
