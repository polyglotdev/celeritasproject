package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fatih/color"
	celeritas "github.com/polyglotdev/celeritasproject"
)

const version = "1.0.0"

var cel celeritas.Celeritas

func main() {
	var message string
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		exitGracefully(err)
	}

	setup()

	validCommands := []string{"help", "version", "make", "migrate", "auth", "model"}
	if !contains(validCommands, arg1) {
		suggestion := findClosestMatch(arg1, validCommands)
		if suggestion != "" {
			exitGracefully(fmt.Errorf("invalid command: %s\nDid you mean '%s'?", arg1, suggestion))
		}
		exitGracefully(fmt.Errorf("invalid command: %s\nRun 'celeritas help' for usage", arg1))
	}

	switch arg1 {
	case "help":
		showHelp()

	case "version":
		color.Yellow("Application version: " + version)

	case "make":
		if arg2 == "" {
			exitGracefully(errors.New("make requires a subcommand: (migration|model|handler|middleware)"))
		}
		err = doMake(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}

	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = doMigrate(arg2, arg3)
		if err != nil {
			exitGracefully(err)
		}
		message = "Migration " + arg2 + " completed"

	default:
		showHelp()
	}
	exitGracefully(nil, message)
}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string

	if len(os.Args) > 1 {
		arg1 = os.Args[1]

		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}

		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		showHelp()
		return "", "", "", errors.New("command required")
	}

	return arg1, arg2, arg3, nil
}

func exitGracefully(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
		os.Exit(1)
	}

	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("process completed successfully")
	}

	os.Exit(0)
}
