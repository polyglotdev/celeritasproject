package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

// setup initializes the application's core configuration by loading environment
// variables and establishing the root path. This approach was chosen over
// hard-coded configuration to support flexible deployment environments and
// follow twelve-factor app principles for configuration management.
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

// getDSN constructs database connection strings with support for multiple database
// types. The function handles PostgreSQL's unique connection string format separately
// from other databases to accommodate its specific security and connection requirements.
// This separation enables clean handling of optional password authentication while
// maintaining consistent connection patterns across different database backends.
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
				os.Getenv("DATABASE_SSL_MODE"),
			)
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"),
			)
		}
		return dsn
	}
	return "mysql://" + cel.BuildDSN()
}

// contains provides a type-safe way to check for string membership in a slice.
// While maps could offer O(1) lookup, this slice-based approach was chosen for
// small collections where the memory overhead of a map would be unnecessary.
func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// findClosestMatch implements a "did you mean?" feature using [Levenshtein distance](https://en.wikipedia.org/wiki/Levenshtein_distance)
// to enhance UX when users make typographical errors in commands. The threshold of 3
// was chosen as it represents a good balance between catching common typos while
// avoiding false positives. This approach is particularly valuable in CLI applications
// where user input errors are common and feedback is immediate.
func findClosestMatch(input string, commands []string) string {
	minDistance := 1000
	var closest string

	for _, cmd := range commands {
		distance := levenshteinDistance(input, cmd)
		if distance < minDistance {
			minDistance = distance
			closest = cmd
		}
	}

	// Only suggest if the distance is reasonable (e.g., less than 3)
	if minDistance <= 3 {
		return closest
	}
	return ""
}

// levenshteinDistance implements the Levenshtein algorithm using a matrix-based
// approach rather than recursion. While this uses more memory (O(mn)), it provides
// better performance through dynamic programming, avoiding the exponential time
// complexity of recursive implementations. This tradeoff favors CPU efficiency
// over memory usage, which is appropriate for the small string lengths typical
// in command names.
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 0; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min(
					matrix[i-1][j]+1,   // deletion
					matrix[i][j-1]+1,   // insertion
					matrix[i-1][j-1]+1, // substitution
				)
			}
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min is a variadic function that finds the minimum value among integers.
// The variadic approach was chosen over multiple parameters or slices to
// provide a more flexible API that can handle any number of comparisons
// while maintaining readable call sites. This is particularly useful in
// the Levenshtein distance calculation where we need to find the minimum
// of exactly three values.
func min(numbers ...int) int {
	result := numbers[0]
	for _, num := range numbers[1:] {
		if num < result {
			result = num
		}
	}
	return result
}

// showHelp uses color-coded output to improve command visibility and user
// experience in terminal environments. The choice of cyan for command help
// provides good visibility across different terminal color schemes while
// maintaining readability. The structured format helps users quickly scan
// available commands and their purposes.
func showHelp() {
	color.Cyan(`Available commands:
	help                  - show the help commands
	version               - print application version
	migrate up            - run all up migrations
	migrate down          - reverses most recent migration
	migrate reset         - drop all tables and re-run all migrations
	make auth             - create and runs auth migrations, models, and middleware
	make session          - create a table in the database as a session store
	make migration <name> - create a new migration files for up and down migrations
	make model <name>     - create a new model file
	make handler <name>   - create a new handler file
	make key              - generate a new encryption key
	`)
}
