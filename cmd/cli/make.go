package main

import (
	"errors"
	"fmt"
	"time"
)

func doMake(arg2, arg3 string) error {
	validSubcommands := []string{"migration", "model", "handler", "middleware", "auth"}
	if !contains(validSubcommands, arg2) {
		suggestion := findClosestMatch(arg2, validSubcommands)
		if suggestion != "" {
			return fmt.Errorf("invalid make subcommand: %s\nDid you mean '%s'?", arg2, suggestion)
		}
		return fmt.Errorf("invalid make subcommand: %s\nValid subcommands are: migration, model, handler, middleware, auth", arg2)
	}

	switch arg2 {
	case "migration":
		dbType := cel.DB.DataType
		if arg3 == "" {
			return errors.New("you must specify a migration name")
		}

		fileName := fmt.Sprintf("%d_%s.go", time.Now().UnixMicro(), arg3)

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

	case "model":
		fileName := cel.RootPath + "/models/" + arg3 + ".go"
		err := copyFileFromTemplate("templates/models/model.go", fileName)
		if err != nil {
			exitGracefully(err)
		}

	case "auth":
		err := doAuth()
		if err != nil {
			exitGracefully(err)
		}
	}
	return nil
}
