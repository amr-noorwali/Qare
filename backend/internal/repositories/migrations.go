package repositories

import (
	"context"
	"database/sql"
	"fmt"
)

// Append a new migration for later schema changes. MySQL DDL commits itself,
// so each statement must be safe to retry after an interrupted deployment.
var migrations = [][]string{{
	`CREATE TABLE IF NOT EXISTS users (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	`CREATE TABLE IF NOT EXISTS sessions (
		token CHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL PRIMARY KEY,
		user_id BIGINT NOT NULL,
		expires_at VARCHAR(32) NOT NULL,
		INDEX idx_sessions_user (user_id),
		CONSTRAINT fk_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	`CREATE TABLE IF NOT EXISTS books (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(500) NOT NULL,
		author VARCHAR(255) NOT NULL,
		description TEXT NOT NULL,
		cover_url TEXT NOT NULL,
		genre VARCHAR(255) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	`CREATE TABLE IF NOT EXISTS reviews (
		id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		book_id BIGINT NOT NULL,
		user_id BIGINT NOT NULL,
		rating TINYINT NOT NULL,
		body TEXT NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY uq_reviews_book_user (book_id, user_id),
		INDEX idx_reviews_book_created (book_id, created_at DESC),
		INDEX idx_reviews_user (user_id),
		CONSTRAINT chk_reviews_rating CHECK (rating BETWEEN 1 AND 5),
		CONSTRAINT fk_reviews_book FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
		CONSTRAINT fk_reviews_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
	`CREATE TABLE IF NOT EXISTS library (
		user_id BIGINT NOT NULL,
		book_id BIGINT NOT NULL,
		status VARCHAR(7) NOT NULL,
		PRIMARY KEY (user_id, book_id),
		INDEX idx_library_book (book_id),
		CONSTRAINT chk_library_status CHECK (status IN ('want', 'reading', 'read')),
		CONSTRAINT fk_library_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		CONSTRAINT fk_library_book FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`,
}}

func migrate(db *sql.DB) error {
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	var locked int
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK('qare_schema_migration', 30)").Scan(&locked); err != nil || locked != 1 {
		return fmt.Errorf("could not acquire schema migration lock: %v", err)
	}
	defer func() {
		var released int
		_ = conn.QueryRowContext(ctx, "SELECT RELEASE_LOCK('qare_schema_migration')").Scan(&released)
	}()
	if _, err := conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INT NOT NULL PRIMARY KEY
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`); err != nil {
		return err
	}
	var version int
	if err := conn.QueryRowContext(ctx, "SELECT COALESCE(MAX(version),0) FROM schema_migrations").Scan(&version); err != nil {
		return err
	}
	if version > len(migrations) {
		return fmt.Errorf("database schema version %d is newer than supported version %d", version, len(migrations))
	}
	for i := version; i < len(migrations); i++ {
		for _, statement := range migrations[i] {
			if _, err := conn.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("migration %d: %w", i+1, err)
			}
		}
		if _, err := conn.ExecContext(ctx, "INSERT INTO schema_migrations(version) VALUES(?)", i+1); err != nil {
			return err
		}
	}
	return nil
}
