package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// doMake implements a code generation system following established patterns from popular
// frameworks like Rails and Laravel. It uses subcommands to generate different types of
// code artifacts (migrations, models, handlers, etc.) while enforcing consistent naming
// and structure across the application.
//
// The function employs a two-level command structure (make + subcommand) to provide
// a familiar interface for developers coming from other frameworks. This approach was
// chosen over a single-level structure to better organize related generation tasks
// and provide clearer command semantics.
//
// Each subcommand implements specific generation logic with careful consideration for:
// - File naming conventions (using snake_case for migrations, camelCase for handlers)
// - Timestamp-based ordering (for migrations)
// - Conflict prevention (checking for existing files)
// - Parameter naming collisions (avoiding conflicts with common variables)
//
// The error handling strategy focuses on early validation and clear error messages,
// including "did you mean?" suggestions for mistyped commands to improve developer
// experience.
//
// Generation templates are embedded in the binary rather than stored as separate files
// to ensure consistency across deployments and eliminate dependency on external files.
//
// Specific subcommand implementations:
//
// Migration:
//   - Uses microsecond precision timestamps to ensure unique ordering
//   - Supports multiple database types through templating
//   - Creates paired up/down migrations for reversibility
//
// Handler:
//   - Implements smart parameter naming to avoid conflicts with http.ResponseWriter (w)
//     and *http.Request (r) parameters
//   - Uses consistent naming conventions for better code organization
//
// Model:
//   - Automatically derives table names using proper pluralization
//   - Implements smart parameter naming to avoid conflicts with common model variables
//   - Uses consistent casing conventions (camelCase for types, snake_case for tables)
//
// Auth:
//   - Generates a complete authentication system with associated migrations
//   - Creates necessary middleware and helper functions
//
// Session:
//   - Generates database-backed session handling
//   - Creates necessary tables and configurations
//
// The function returns errors instead of handling them directly to allow the caller
// to implement custom error handling strategies, except for fatal errors that should
// terminate execution.
func doMake(arg2, arg3 string) error {
	validCommands := []string{"key", "migration", "model", "handler", "middleware"}
	if !contains(validCommands, arg2) {
		suggestion := findClosestMatch(arg2, validCommands)
		if suggestion != "" {
			return fmt.Errorf("invalid 'make' subcommand: %s\nDid you mean '%s'?\nValid subcommands are: key, migration, model, handler, middleware", arg2, suggestion)
		}
		return fmt.Errorf("invalid 'make' subcommand: %s\nValid subcommands are: key, migration, model, handler, middleware", arg2)
	}

	switch arg2 {
	case "key":
		// Check if .env exists and if it contains KEY=
		if _, err := os.Stat(".env"); err == nil {
			content, err := os.ReadFile(".env")
			if err != nil {
				return err
			}

			if strings.Contains(string(content), "KEY=") {
				// Ask for confirmation
				fmt.Print("KEY already exists in .env. Do you want to overwrite it? (y/n): ")
				var response string
				_, err = fmt.Scanln(&response)
				if err != nil {
					exitGracefully(err)
				}

				if strings.ToLower(response) != "y" {
					color.Yellow("Operation cancelled")
					return nil
				}
			}
		}

		// Generate the key
		rnd := cel.RandomString(32)

		// Read existing content
		var existingContent string
		if data, err := os.ReadFile(".env"); err == nil {
			existingContent = string(data)
		}

		// Remove existing KEY line if it exists
		lines := strings.Split(existingContent, "\n")
		var newLines []string
		for _, line := range lines {
			if !strings.HasPrefix(line, "KEY=") {
				newLines = append(newLines, line)
			}
		}

		// Add new KEY
		newLines = append(newLines, fmt.Sprintf("KEY=%s", rnd))

		// Write back to file
		err := os.WriteFile(".env", []byte(strings.Join(newLines, "\n")), 0644)
		if err != nil {
			return err
		}

		color.Yellow("Encryption key: " + rnd)
		color.Green("Key has been added to .env file")

	case "migration":
		dbType := cel.DB.DataType
		if arg3 == "" {
			return errors.New("you must specify a migration name")
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), strcase.ToSnake(arg3))

		upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			exitGracefully(err)
		}

		err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			exitGracefully(err)
		}

	case "handler":
		if arg3 == "" {
			exitGracefully(errors.New("you must specify a handler name"))
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(fmt.Errorf("%s already exists", fileName))
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		handler := string(data)
		handlerName := strcase.ToCamel(arg3)

		// Get first letter for param name
		firstLetter := strings.ToLower(string(handlerName[0]))

		// Check if firstLetter conflicts with 'w' or 'r'
		paramName := firstLetter
		if firstLetter == "w" || firstLetter == "r" {
			paramName = "h" // Use 'h' as fallback
		}

		handler = strings.ReplaceAll(handler, "$HANDLERNAME", handlerName)
		handler = strings.ReplaceAll(handler, "$FIRSTLETTER", paramName)

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			exitGracefully(err)
		}

	case "model":
		if arg3 == "" {
			exitGracefully(errors.New("you must specify a model name"))
		}

		fileName := cel.RootPath + "/data/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			exitGracefully(fmt.Errorf("%s already exists", fileName))
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			exitGracefully(err)
		}

		model := string(data)
		modelName := strcase.ToCamel(arg3)

		// Get first letter for param name
		firstLetter := strings.ToLower(string(modelName[0]))

		// Check if firstLetter conflicts with 't', 'm', or 'c' (table, model, collection)
		paramName := firstLetter
		if firstLetter == "t" || firstLetter == "m" || firstLetter == "c" {
			paramName = "d" // Use 'd' as fallback (for 'data')
		}

		// Initialize pl
		pl := pluralize.NewClient()

		// Get pluralized, snake_cased table name
		tableName := strcase.ToSnake(pl.Plural(modelName))

		model = strings.ReplaceAll(model, "$MODELNAME$", modelName)
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)
		model = strings.ReplaceAll(model, "$FIRSTLETTER$", paramName)

		err = os.WriteFile(fileName, []byte(model), 0644)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}

	case "session":
		err := doSessionTable()
		if err != nil {
			exitGracefully(err)
		}
	}
	return nil
}
