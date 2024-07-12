package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sha1sof/authRESTAPI/internal/config"
	"github.com/sha1sof/authRESTAPI/internal/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

// DBPostgres is a framework for managing
// and connecting to the PostgreSQL database and related configurations.
type DBPostgres struct {
	db     *sql.DB
	logger *slog.Logger
	Cost   int
	Secret string
	Time   time.Duration
}

// New initializes and returns a new instance of DB Postgres,
// which represents a connection to a PostgreSQL database.
func New(cfg *config.Config, log *slog.Logger, cost int, secret string, timeD time.Duration) (*DBPostgres, error) {
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
		log.Info("Failed to connect to database: " + err.Error())
		panic("Failed to connect to database: " + err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Info("Failed to ping database: " + err.Error())
		panic("Failed to ping database: " + err.Error())
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);
	CREATE INDEX IF NOT EXISTS idx_email ON users (email);
	`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Info("Failed to create table: " + err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Successfully created table: " + op)

	return &DBPostgres{
		db:     db,
		logger: log,
		Cost:   cost,
		Secret: secret,
		Time:   timeD,
	}, nil
}

// RegisterUser is a method of the DB Postgres struct
// that handles the registration of a new user.
func (s *DBPostgres) RegisterUser(login, pass string) (bool, error) {
	const op = "storage.dbPostgres.RegisterUser"
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := s.db.QueryRow(query, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		password := []byte(pass)
		hash, err := bcrypt.GenerateFromPassword(password, s.Cost)
		if err != nil {
			s.logger.Info("Failed to hash password: "+op+": ", err)
			return false, fmt.Errorf("%s: %w", op, err)
		}

		query = "INSERT INTO users (email, password) VALUES ($1, $2)"
		_, err = s.db.Exec(query, login, string(hash))
		if err != nil {
			s.logger.Info(op+": ", err)
			return false, fmt.Errorf("%s: %w", op, err)
		}

		s.logger.Info("Created a user with: " + login + ".")
		return true, nil
	} else {
		s.logger.Info("There is already a user with such an email: " + login + ".")
		return false, nil
	}
}

// LoginUser is a method of the DBPostgres struct
// that handles the login process for a user.
func (s *DBPostgres) LoginUser(login, pass string) (string, bool, error) {
	const op = "storage.dbPostgres.LoginUser"
	var exists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err := s.db.QueryRow(query, login).Scan(&exists)
	if err != nil {
		s.logger.Info(op+": Error checking if user exists: ", err)
		return "", false, fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		s.logger.Info(op + ": User exists: " + login)
		query = "SELECT email, password FROM users WHERE email = $1"
		var email, passHash string
		err = s.db.QueryRow(query, login).Scan(&email, &passHash)
		if err != nil {
			s.logger.Info(op+": Error retrieving user: ", err)
			return "", false, fmt.Errorf("%s: %w", op, err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass))
		if err != nil {
			s.logger.Info(op + ": The hash does not match the password: " + login)
			return "", false, fmt.Errorf("%s: %w", op, err)
		}
		token, err := jwt.NewToken(email, s.Time, s.Secret)
		if err != nil {
			s.logger.Info(op+": Error creating token: ", err)
			return "", false, fmt.Errorf("%s: %w", op, err)
		}

		return token, true, nil

	} else {
		s.logger.Info(op + ": The user is not registered: " + login + ".")
		return "", false, nil
	}

	return "", exists, nil
}
