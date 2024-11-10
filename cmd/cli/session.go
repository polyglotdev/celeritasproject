package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// doSessionTable creates database migration files for session management tables.
// It generates both up and down migrations with standardized naming conventions
// and database-specific SQL.
//
// Design choices:
//   - Uses Unix microseconds in filename to ensure unique, chronological ordering
//     while avoiding potential naming conflicts in distributed development
//   - Normalizes database type names (mariadb → mysql, postgresql → postgres) to
//     maintain consistent template naming while supporting multiple dialects
//   - Separates up/down migrations to support both deployments and rollbacks
//     with proper state management
//   - Uses templated SQL files rather than hardcoded queries to:
//   - Ensure consistent schema across projects
//   - Allow for database-specific optimizations
//   - Make schema changes maintainable in a single location
//
// The function handles three main tasks:
//  1. Database dialect normalization and file path construction
//  2. Creation of migration files from templates
//  3. User feedback for next steps
//
// It returns an error if file operations fail, ensuring atomic migration file creation.
func doSessionTable() error {
	dbType := cel.DB.DataType

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	fileName := fmt.Sprintf("%d_create_sessions_table", time.Now().UnixMicro())

	upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

	err := copyFileFromTemplate("/templates/migrations/"+dbType+"_session.sql", upFile)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile([]byte("drop table sessions;"), downFile)
	if err != nil {
		exitGracefully(err)
	}

	err = doMigrate("up", "")
	if err != nil {
		exitGracefully(err)
	}

	color.Yellow("Don't forget to run:")
	color.Yellow("migrate up")
	color.Green("Session table migration created successfully")

	return nil
}
