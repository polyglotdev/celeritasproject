package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/fatih/color"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath, destinationPath string) error {
	if fileExists(destinationPath) {
		return fmt.Errorf("%s already exists", destinationPath)
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(data, destinationPath)
	if err != nil {
		exitGracefully(err)
	}

	color.Green("created template at %s", destinationPath)
	return nil
}

func copyDataToFile(data []byte, destinationPath string) error {
	err := os.WriteFile(destinationPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}
