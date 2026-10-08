package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Open connects to an existing MySQL database, applies schema migrations and
// seeds a new installation. The database itself must be created beforehand.
func Open(dsn string) (*sql.DB, error) {
	return open(dsn, true)
}

// OpenForImport prepares an empty MySQL database without demo records.
func OpenForImport(dsn string) (*sql.DB, error) {
	return open(dsn, false)
}

func open(dsn string, withSeed bool) (*sql.DB, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid MySQL DSN: %w", err)
	}
	config.ParseTime = false // Keep the API's existing SQL timestamp strings.
	if config.Params == nil {
		config.Params = make(map[string]string)
	}
	config.Params["charset"] = "utf8mb4"
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(3 * time.Minute)
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to MySQL: %w", err)
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	if withSeed {
		if err := seed(db); err != nil {
			db.Close()
			return nil, err
		}
	}
	return db, nil
}
