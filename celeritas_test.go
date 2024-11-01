package celeritas

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCeleritas_New(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Set up environment variables for testing
	err := os.Setenv("DEBUG", "true")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("PORT", "8080")
	if err != nil {
		t.Fatal(err)
	}
	err = os.Setenv("RENDERER", "go")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		rootPath string
		envVars  map[string]string
		wantErr  bool
	}{
		{
			name:     "valid path with env vars",
			rootPath: tempDir,
			envVars: map[string]string{
				"DEBUG":    "true",
				"PORT":     "8080",
				"RENDERER": "go",
			},
			wantErr: false,
		},
		{
			name:     "invalid path",
			rootPath: "/path/that/cannot/be/created",
			envVars:  nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			// Set up environment variables for this test
			if tt.envVars != nil {
				for k, v := range tt.envVars {
					err := os.Setenv(k, v)
					if err != nil {
						ts.Fatal(err)
					}
				}
			}

			c := &Celeritas{}
			err := c.New(tt.rootPath)

			// Check if error matches expected
			if (err != nil) != tt.wantErr {
				ts.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If no error, verify the setup
			if !tt.wantErr {
				// Test directory creation
				expectedDirs := []string{
					"handlers", "migrations", "views", "data",
					"public", "tmp", "logs", "middleware",
				}

				for _, dir := range expectedDirs {
					path := filepath.Join(tt.rootPath, dir)
					if _, err := os.Stat(path); os.IsNotExist(err) {
						ts.Errorf("Directory not created: %s", path)
					}
				}

				// Test struct field initialization
				appSettingsTests := []struct {
					name   string
					got    interface{}
					want   interface{}
					errMsg string
				}{
					{"AppName", c.AppName, "celeritas", "AppName not set correctly"},
					{"Debug", c.Debug, true, "Debug not set correctly"},
					{"Version", c.Version, Version, "Version not set correctly"},
					{"RootPath", c.RootPath, tt.rootPath, "RootPath not set correctly"},
					{"config.port", c.config.port, "8080", "Port not set correctly"},
					{"config.renderer", c.config.renderer, "go", "Renderer not set correctly"},
				}

				for _, tst := range appSettingsTests {
					if tst.got != tst.want {
						ts.Errorf("%s: got %v, want %v", tst.errMsg, tst.got, tst.want)
					}
				}

				// Test logger initialization
				if c.InfoLog == nil {
					ts.Error("InfoLog not initialized")
				}
				if c.ErrorLog == nil {
					ts.Error("ErrorLog not initialized")
				}
			}

			// Clean up environment variables
			if tt.envVars != nil {
				for k := range tt.envVars {
					if err := os.Unsetenv(k); err != nil {
						ts.Fatal(err)
					}
				}
			}
		})
	}
}
