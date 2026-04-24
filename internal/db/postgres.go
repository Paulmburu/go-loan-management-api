package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// NewPostgresConnection opens and verifies a Postgres database connection.
func NewPostgresConnection(databaseURL string) (*sql.DB, error) {
	// Open creates a database handle, but does not fully verify the connection yet.
	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Open creates a database handle, but does not fully verify the connection yet.
	//
	// We declare 'err' and check it in one line.
	// This keeps 'err' scoped ONLY to this error-handling block.
	//
	// err := db.Ping() // 'err' is now available for the REST of the function
	// if err != nil {
	//     return err
	// }
	// 'err' still exists here, which might be confusing
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}
