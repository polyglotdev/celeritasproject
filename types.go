package celeritas

import "database/sql"

// initPaths represents the paths and folder names to initialize the application
type initPaths struct {
	rootPath    string
	folderNames []string
}

// cookieConfig represents the configuration for the session cookie
type cookieConfig struct {
	name     string
	lifetime string
	persist  string
	secure   string
	domain   string
}

type databaseConfig struct {
	dsn      string
	database string
}

type Database struct {
	DataType string
	Pool     *sql.DB
}
