package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" 
)

var db *sql.DB

func InitDB() {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error

	db, err = sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	const maxRetries = 5
	for i := range maxRetries {
		if err = db.Ping(); err == nil {
			fmt.Println("✅ Successfully connected to the database with pgx.")
			return
		}
		log.Printf("⏳ DB not ready (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	log.Fatalf("❌ Could not connect to database after %d attempts: %v", maxRetries, err)
}

func GetDB() *sql.DB {
	if db == nil {
		log.Fatalf("Database connection not initialized")
	}
	return db
}

func CloseDB(){
	if err := db.Close(); err != nil {
		log.Fatalf("Error Closing the Database :%v",err)
	}
}
