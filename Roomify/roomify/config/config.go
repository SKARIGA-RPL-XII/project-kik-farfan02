package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	JWTKey string
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	return &Config{
		JWTKey: os.Getenv("JWT_KEY"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPass: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
	}
}

func ConnectDB(cfg *Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to open database", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	DB = db
}