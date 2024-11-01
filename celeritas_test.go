package celeritas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCeleritas_New(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name     string
		rootPath string
		wantErr  bool
	}{
		{
			name:     "valid path",
			rootPath: tempDir,
			wantErr:  false,
		},
		{
			name:     "invalid path",
			rootPath: "/path/that/cannot/be/created",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			c := &Celeritas{}
			err := c.New(tt.rootPath)

			// Check if error matches expected
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error, verify all directories were created
			if !tt.wantErr {
				expectedDirs := []string{
					"handlers",
					"migrations",
					"views",
					"data",
					"public",
					"tmp",
					"logs",
					"middleware",
				}

				for _, dir := range expectedDirs {
					path := filepath.Join(tt.rootPath, dir)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						t.Errorf("Directory not created: %s", path)
					}
				}

				// Verify struct fields were set correctly
				if c.AppName != "celeritas" {
					t.Errorf("AppName not set correctly, got %s", c.AppName)
				}
				if !c.Debug {
					t.Error("Debug not set correctly")
				}
				if c.Version != Version {
					t.Errorf("Version not set correctly, got %s, want %s", c.Version, Version)
				}
			}
		})
	}
}
