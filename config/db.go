package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func InitDB() {
	cstr, ok := os.LookupEnv("DB_URL")
	if !ok {
		logger.Fatal("Couldn't find DB_URL variable")
		os.Exit(1)
	}
	var err error
	db, err = sql.Open("postgres", cstr)
	if err != nil {
		logger.Fatalf("Couldn't open a connection to DB: %v", err)
		os.Exit(1)
	}
	db.SetMaxOpenConns(16)

	err = db.Ping()
	if err != nil {
		logger.Errorf("Couldn't ping to DB: %v", err)
		os.Exit(1)
	}
}

func GetDB() *sql.DB {
	return db
}
