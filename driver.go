package celeritas

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func (c *Celeritas) OpenDB(dbType, dsn string) (*sql.DB, error) {
	// The if statement is used because we specify the driver type in
	// the .env file, and to avoid the issue of user confusion on postgres
	// vs pgx, because pgx is what we need to connect to postgres.
	if dbType == "postgres" || dbType == "postgresql" {
		dbType = "pgx"
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return nil, err
	}

	// We are pinging the database to ensure that the connection is
	// working and that the database is up and running.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
