// todoapp/config/config.go

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetDBConnectionString() string {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", user, password, dbName, host, port, dbSSLMode)

	return connStr
}
