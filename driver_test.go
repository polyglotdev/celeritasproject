package celeritas

import (
	"os"
	"testing"
)

func TestCeleritas_OpenDB(t *testing.T) {
	tests := []struct {
		name    string
		dbType  string
		dsn     string
		wantErr bool
	}{
		{
			name:    "postgres type",
			dbType:  "postgres",
			dsn:     os.Getenv("POSTGRES_TEST_DSN"),
			wantErr: false,
		},
		{
			name:    "postgresql type",
			dbType:  "postgresql",
			dsn:     os.Getenv("POSTGRES_TEST_DSN"),
			wantErr: false,
		},
		{
			name:    "invalid database type",
			dbType:  "invalid",
			dsn:     "invalid",
			wantErr: true,
		},
		{
			name:    "invalid dsn",
			dbType:  "pgx",
			dsn:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			c := &Celeritas{}

			// Skip postgres tests if no DSN provided
			if (tt.dbType == "postgres" || tt.dbType == "postgresql") && tt.dsn == "" {
				ts.Skip("Skipping postgres test - no DSN provided")
			}

			db, err := c.OpenDB(tt.dbType, tt.dsn)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				ts.Errorf("OpenDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we got a db connection, verify it works and clean up
			if db != nil {
				err = db.Ping()
				if err != nil {
					ts.Errorf("OpenDB() returned unpingable database: %v", err)
				}
				defer func() {
					err := db.Close()
					if err != nil {
						ts.Errorf("Error closing database: %v", err)
					}
				}()
			}
		})
	}
}
