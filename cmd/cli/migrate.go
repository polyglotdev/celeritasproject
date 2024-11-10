package main

import (
	"fmt"
)

// doMigrate handles database migration operations with input validation and helpful error messages.
// It supports different migration strategies through subcommands while maintaining atomic operations
// and proper error handling.
//
// The function implements three key migration strategies:
//   - "up": Applies all pending migrations
//   - "down": Rolls back either the most recent migration or all migrations
//   - "reset": Performs a complete database reset by rolling back all migrations
//     and then re-applying them, useful for development and testing environments
//
// The function includes input validation that:
//   - Validates subcommands against a known list of valid options
//   - Provides helpful suggestions for mistyped commands using string similarity
//   - Returns clear error messages showing valid subcommand options
//
// Parameters:
//   - arg2: The migration subcommand ("up", "down", or "reset")
//   - arg3: Additional parameter for the "down" command ("all" to reverse all migrations)
//
// Returns:
//   - error: Returns nil on successful migration, otherwise returns either:
//   - A validation error with suggestions for invalid commands
//   - The first encountered error during the migration process
//
// The function uses the application's DSN (Data Source Name) configuration
// to maintain consistent database connectivity across all migration operations.
// Each operation is designed to be atomic - either completely succeeding or
// failing with appropriate error information.
func doMigrate(arg2, arg3 string) error {
	validSubcommands := []string{"up", "down", "reset"}
	if !contains(validSubcommands, arg2) {
		suggestion := findClosestMatch(arg2, validSubcommands)
		if suggestion != "" {
			return fmt.Errorf("invalid migrate subcommand: %s\nDid you mean '%s'?", arg2, suggestion)
		}
		return fmt.Errorf(
			"invalid migrate subcommand: %s\nValid subcommands are: up, down, reset",
			arg2,
		)
	}

	dsn := getDSN()

	switch arg2 {
	case "up":
		err := cel.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 == "all" {
			err := cel.MigrateDownAll(dsn)
			if err != nil {
				return err
			}
		} else {
			err := cel.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := cel.MigrateDownAll(dsn)
		if err != nil {
			return err
		}
		err = cel.MigrateUp(dsn)
		if err != nil {
			return err
		}
	}
	return nil
}
