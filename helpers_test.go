package celeritas

import (
	"crypto/rand"
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

func TestCreateFileIfNotExists(t *testing.T) {
	c := &Celeritas{}

	// Create a temporary directory for testing
	tmpDir := os.TempDir()
	testPath := filepath.Join(tmpDir, "test_file.txt")

	// Clean up after the test
	defer func() {
		_ = os.Remove(testPath)
	}()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "create new file",
			path:    testPath,
			wantErr: false,
		},
		{
			name:    "file already exists",
			path:    testPath,
			wantErr: false,
		},
		{
			name:    "invalid path",
			path:    "/nonexistent/directory/file.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.CreateFileIfNotExists(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileIfNotExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Check if file exists
				_, err := os.Stat(tt.path)
				if os.IsNotExist(err) {
					t.Errorf("CreateFileIfNotExists() file was not created at %s", tt.path)
				}
			}
		})
	}
}

func TestEncryptionMethods(t *testing.T) {
	// Create encryption instance with a test key (must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256)
	key := make([]byte, 32) // Using AES-256
	_, err := rand.Read(key)
	if err != nil {
		t.Fatalf("Failed to generate random key: %v", err)
	}

	enc := &Encryption{Key: key}

	tests := []struct {
		name      string
		inputData string
		wantErr   bool
	}{
		{
			name:      "encrypt and decrypt empty string",
			inputData: "",
			wantErr:   false,
		},
		{
			name:      "encrypt and decrypt normal string",
			inputData: "Hello, World!",
			wantErr:   false,
		},
		{
			name:      "encrypt and decrypt long string",
			inputData: "This is a longer string that we'll use to test encryption and decryption with multiple blocks",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encryption
			encrypted, err := enc.Encrypt(tt.inputData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Test decryption
				decrypted, err := enc.Decrypt(encrypted)
				if err != nil {
					t.Errorf("Decrypt() error = %v", err)
					return
				}

				// Compare original and decrypted data
				if decrypted != tt.inputData {
					t.Errorf("Decrypt() = %v, want %v", decrypted, tt.inputData)
				}
			}
		})
	}

	t.Run("invalid decrypt input", func(t *testing.T) {
		_, err := enc.Decrypt("invalid base64 data")
		if err == nil {
			t.Error("Decrypt() expected error for invalid input")
		}
	})
}
