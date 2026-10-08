package repositories

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
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
	if err := config.Apply(mysql.Charset("utf8mb4", "")); err != nil {
		return nil, err
	}
	if caPath := os.Getenv("MYSQL_CA_CERT_PATH"); caPath != "" {
		pem, err := os.ReadFile(caPath)
		if err != nil {
			return nil, fmt.Errorf("read MySQL CA certificate: %w", err)
		}
		roots := x509.NewCertPool()
		if !roots.AppendCertsFromPEM(pem) {
			return nil, fmt.Errorf("MySQL CA certificate contains no valid certificates")
		}
		config.TLS = &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
	}
	connector, err := mysql.NewConnector(config)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(connector)
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
