package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	Conn *sql.DB
}

type DBParams struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func New(conn *sql.DB) *Database {
	return &Database{
		Conn: conn,
	}
}

func Connect(params DBParams) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		params.Host, params.Port, params.Username, params.Password, params.Database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established")

	return conn, nil
}
