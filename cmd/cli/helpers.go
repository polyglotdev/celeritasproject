package main

import (
	"os"

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
