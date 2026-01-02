package model

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	dbHost      = "db"
	dbPort      = 5432
	dbUser      = "postgres"
	dbName      = "skvdmt_ru"
	DB_PASSWORD = "DB_PASSWORD"
)

var DB *sql.DB

// PostgressConnect соединение с базой данных postgress и
// установка дескриптора соединения в глобальную пременную DB
func PostgressConnect() error {
	dbPass, ok := os.LookupEnv(DB_PASSWORD)
	if !ok {
		return fmt.Errorf("env %s not set", DB_PASSWORD)
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	))
	if err != nil {
		return err
	}
	DB = db
	return nil
}
