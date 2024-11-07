package celeritas

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCeleritas_RandomString(t *testing.T) {
	c := &Celeritas{}
	tests := []struct {
		name   string
		length int
		want   int // testing length of output
	}{
		{
			name:   "zero length",
			length: 0,
			want:   0,
		},
		{
			name:   "length of 10",
			length: 10,
			want:   10,
		},
		{
			name:   "length of 20",
			length: 20,
			want:   20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			got := c.RandomString(tt.length)

			// Check length
			if len(got) != tt.want {
				ts.Errorf("RandomString() length = %v, want %v", len(got), tt.want)
			}

			// Check characters are valid
			if tt.length > 0 {
				for _, char := range got {
					if !strings.ContainsRune(allowedChars, char) {
						ts.Errorf("RandomString() contains invalid character: %c", char)
					}
				}
			}
		})
	}
}

func TestCeleritas_CreateDirIfNotExist(t *testing.T) {
	c := &Celeritas{}
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid new directory",
			path:    filepath.Join(tmpDir, "testdir"),
			wantErr: false,
		},
		{
			name:    "existing directory",
			path:    tmpDir,
			wantErr: false,
		},
		{
			name:    "invalid path",
			path:    "/path/that/cannot/be/created/due/to/permissions",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			err := c.CreateDirIfNotExist(tt.path)
			if (err != nil) != tt.wantErr {
				ts.Errorf("CreateDirIfNotExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if directory exists
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					ts.Errorf("CreateDirIfNotExist() directory was not created at %v", tt.path)
				}
			}
		})
	}
}

func TestCeleritas_CreateFileIfNotExists(t *testing.T) {
	c := &Celeritas{}
	tmpDir := t.TempDir()

	// Create an existing file
	existingFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "new file",
			path:    filepath.Join(tmpDir, "newfile.txt"),
			wantErr: false,
		},
		{
			name:    "existing file",
			path:    existingFile,
			wantErr: false,
		},
		{
			name:    "invalid path",
			path:    "/path/that/cannot/be/created/file.txt",
			wantErr: true,
		},
		{
			name:    "directory path",
			path:    tmpDir,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(ts *testing.T) {
			err := c.CreateFileIfNotExists(tt.path)
			if (err != nil) != tt.wantErr {
				ts.Errorf("CreateFileIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if file exists
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					ts.Errorf("CreateFileIfNotExists() file was not created at %v", tt.path)
				}
			}
		})
	}
}
