package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// setup initializes the application environment by loading environment variables
// and setting up the root path. It loads variables from a .env file and sets
// the database type from the DATABASE_TYPE environment variable. If any operation
// fails, the program exits gracefully with an error message.
func setup() {
	err := godotenv.Load()
	if err != nil {
		exitGracefully(err)
	}

	path, err := os.Getwd()
	if err != nil {
		exitGracefully(err)
	}

	cel.RootPath = path
	cel.DB.DataType = os.Getenv("DATABASE_TYPE")
}

func getDSN() string {
	dbType := cel.DB.DataType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}
	return "mysql://" + cel.BuildDSN()
}

func showHelp() {
	color.Yellow(`Available commands:
	help           				   - show the help commands
	version        				   - print application version
	migrate up     				   - runs all up migrations that have not been run previously
	migrate down   				   - reverses the most recent migration
	migrate reset  				   - reverses all migrations and runs them again
	make migration <name>   		   - creates a new migration file
	make auth                 	           - creates and runs migrations for authentication tables, and creates models and middleware
	make handler <name>   				   - creates a new handler file
	`)
}
