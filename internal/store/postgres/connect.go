package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DatabaseParameters struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	SSLMode  string
}

func Connect(dbParams *DatabaseParameters) (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=%s",
			dbParams.Host,
			dbParams.Port,
			dbParams.Username,
			dbParams.Password,
			dbParams.Database,
			dbParams.SSLMode),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
