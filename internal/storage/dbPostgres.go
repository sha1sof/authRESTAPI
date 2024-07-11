package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sha1sof/authRESTAPI/internal/config"
)

type DBPostgres struct {
	db *sql.DB
}

// New initializes and returns a new instance of DB Postgres,
// which represents a connection to a PostgreSQL database.
func New(cfg *config.Config) (*DBPostgres, error) {
	const op = "storage.dbPostgres.New"

	pathToPostgre := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode)
	db, err := sql.Open("postgres", pathToPostgre)
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic("Failed to ping database: " + err.Error())
	}

	// TODO: Most likely, you will have to redo the tables.
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);
	CREATE TABLE IF NOT EXISTS tokens (
		id SERIAL PRIMARY KEY,
		token VARCHAR(255) NOT NULL,
		timeout INTEGER NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_email ON users (email);
	CREATE INDEX IF NOT EXISTS idx_token ON tokens (token);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &DBPostgres{db: db}, nil
}
