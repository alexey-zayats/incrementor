package db

import (
	"database/sql"
	"fmt"

	// Register postgres libpq
	_ "github.com/lib/pq"

	"github.com/pkg/errors"
	"incrementor/internal/config"
)

// Connect open connection to database
func Connect(config *config.Config) (*sql.DB, error) {
	// postgres://username:password@hostname:port/database?sslmode=disable
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		config.SQL.Driver,
		config.SQL.Username, config.SQL.Password,
		config.SQL.Hostname, config.SQL.Port,
		config.SQL.Database, config.SQL.SslMode)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error connecting to PostgreSQL with dsn '%s'", dsn))
	}

	return db, nil
}
