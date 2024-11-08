package main

import (
	"embed"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(templatePath, targetFile string) error {
	// check if the file exists in the templateFS
	_, err := templateFS.Open(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		exitGracefully(err)
	}

	err = copyDataToFile(targetFile, data)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func copyDataToFile(targetFile string, data []byte) error {
	err := os.WriteFile(targetFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
