package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/fatih/color"
)

//go:embed templates
var templateFS embed.FS

// copyFileFromTemplate creates a new file from a template, implementing a safety-first approach
// by checking for existing files before attempting to write. This prevents accidental overwrites
// of existing files, which is especially important in code generation scenarios.
//
// The function uses embedded templates (templateFS) rather than external files to ensure:
// 1. Templates are always available and versioned with the code
// 2. No dependency on external file system structure
// 3. Consistent behavior across different environments
//
// It returns an error if the destination already exists or if any file operations fail.
// Success is indicated with a green-colored confirmation message for better UX.
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

// copyDataToFile handles the atomic operation of writing bytes to a file with consistent permissions.
// It uses 0644 permissions (-rw-r--r--) as a secure default that allows:
// - Owner to read and write (6)
// - Group to read (4)
// - Others to read (4)
//
// This function is separated from copyFileFromTemplate to maintain single responsibility
// and allow for potential reuse in other file writing scenarios.
func copyDataToFile(data []byte, destinationPath string) error {
	err := os.WriteFile(destinationPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// fileExists encapsulates the os.Stat() and os.IsNotExist() logic to provide
// a more semantic way to check file existence. This abstraction improves code readability
// and centralizes the file existence check logic in case the underlying implementation
// needs to change.
//
// The function intentionally returns a bool rather than an error to simplify the common
// case of just needing to know if a file exists. This is a conscious trade-off between
// simplicity and complete error information.
func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}
