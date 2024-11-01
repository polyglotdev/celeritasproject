package celeritas

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/polyglotdev/celeritasproject/render"
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
		testCase := tt // Better name that describes the variable's purpose
		t.Run(testCase.name, func(ts *testing.T) {
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
					got    any
					want   any
					errMsg string
				}{
					{"AppName", c.AppName, "celeritas", "AppName not set correctly"},
					{"Debug", c.Debug, true, "Debug not set correctly"},
					{"Version", c.Version, Version, "Version not set correctly"},
					{"RootPath", c.RootPath, tt.rootPath, "RootPath not set correctly"},
					{"config.port", c.config.port, "8080", "Port not set correctly"},
					{"config.renderer", c.config.renderer, "go", "Renderer not set correctly"},
					// Add new test for Renderer initialization
					{"Renderer", c.Render, &render.Render{
						Renderer: c.config.renderer,
						RootPath: c.RootPath,
						Port:     c.config.port,
					}, "Renderer not initialized correctly"},
				}

				for _, tst := range appSettingsTests {
					testSetting := tst // Better name that describes the variable's purpose
					if testSetting.name == "Renderer" {
						got := testSetting.got.(*render.Render)
						want := testSetting.want.(*render.Render)
						if got.Renderer != want.Renderer ||
							got.RootPath != want.RootPath ||
							got.Port != want.Port ||
							got.Secure != want.Secure ||
							got.ServerName != want.ServerName {
							ts.Errorf("%s:\ngot:  %+v\nwant: %+v",
								testSetting.errMsg,
								got,
								want)
						}
					} else if testSetting.got != testSetting.want {
						ts.Errorf("%s: got %v, want %v",
							testSetting.errMsg,
							testSetting.got,
							testSetting.want)
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

func TestCeleritas_ListenAndServe(t *testing.T) {
	var logBuffer bytes.Buffer

	tests := []struct {
		name          string
		setup         func(*Celeritas)
		wantPort      string
		wantLogMsg    string
		wantErr       bool
		checkTimeouts bool
	}{
		{
			name: "valid port configuration",
			setup: func(c *Celeritas) {
				c.AppName = "test_app"
				c.config.port = "0"
				c.InfoLog = log.New(&logBuffer, "INFO\t", log.Ldate|log.Ltime)
				c.ErrorLog = log.New(io.Discard, "", 0)
				c.Routes = chi.NewRouter()
			},
			wantPort:      "0",
			wantLogMsg:    "Starting test_app on port 0",
			wantErr:       false,
			checkTimeouts: true,
		},
		{
			name: "missing port configuration",
			setup: func(c *Celeritas) {
				c.AppName = "test_app"
				c.config.port = ""
				c.InfoLog = log.New(&logBuffer, "INFO\t", log.Ldate|log.Ltime)
				c.ErrorLog = log.New(io.Discard, "", 0)
				c.Routes = chi.NewRouter()
			},
			wantPort:   "",
			wantLogMsg: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(ts *testing.T) {
			logBuffer.Reset()
			c := &Celeritas{}
			testCase.setup(c)

			// Test the actual ListenAndServe method
			errChan := make(chan error, 1)
			go func() {
				errChan <- c.ListenAndServe()
			}()

			// Give the server a moment to start
			time.Sleep(100 * time.Millisecond)

			// Check for errors
			select {
			case err := <-errChan:
				if !testCase.wantErr && err != nil {
					ts.Errorf("ListenAndServe() unexpected error: %v", err)
				}
				if testCase.wantErr && err == nil {
					ts.Error("ListenAndServe() expected error but got none")
				}
			default:
				if testCase.wantErr {
					ts.Error("ListenAndServe() expected error but got none")
				}
			}

			// Verify log message
			if testCase.wantLogMsg != "" {
				logMsg := strings.TrimSpace(logBuffer.String())
				if !strings.Contains(logMsg, testCase.wantLogMsg) {
					ts.Errorf("Log message incorrect\nwant: %q\ngot: %q",
						testCase.wantLogMsg,
						logMsg)
				}
			}

			// Verify port configuration
			if c.config.port != testCase.wantPort {
				ts.Errorf("Port configuration incorrect\nwant: %q\ngot: %q",
					testCase.wantPort,
					c.config.port)
			}
		})
	}
}
