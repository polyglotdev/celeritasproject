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

// databaseConfig represents the configuration for the database
type databaseConfig struct {
	dsn      string
	database string
}

// Database represents the database connection
type Database struct {
	DataType string
	Pool     *sql.DB
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}
