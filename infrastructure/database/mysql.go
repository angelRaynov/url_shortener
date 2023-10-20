package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"url_shortener/infrastructure/config"
)

type Logger interface {
	Fatalw(msg string, keysAndValues ...interface{})
}

func Init(cfg *config.Application, l Logger) *sql.DB {
	// Create a connection string
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Open a connection to the database
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		l.Fatalw("Error opening database connection:", "error", err)
	}

	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		l.Fatalw("Error connecting to the database:", "error", err)
	}

	return db
}
