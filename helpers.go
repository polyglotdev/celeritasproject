package celeritas

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_+"

// RandomString generates a random string of a given length using a cryptographically secure random number generator.
// It returns a string of random characters from the allowedChars constant.
func (c *Celeritas) RandomString(length int) string {
	randomChars, charSet := make([]rune, length), []rune(allowedChars)

	for i := range randomChars {
		prime, _ := rand.Prime(rand.Reader, len(charSet))
		randomInt, maxValue := prime.Uint64(), uint64(len(charSet))
		randomChars[i] = charSet[randomInt%maxValue]
	}
	return string(randomChars)
}

// CreateDirIfNotExist takes a path string
// and returns an error. It creates the directory if it does not exist.
func (c *Celeritas) CreateDirIfNotExist(path string) error {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, mode)
	}
	return nil
}

// CreateFileIfNotExists creates a new empty file at the specified path if one doesn't exist.
// If the file already exists, it does nothing and returns nil.
// The created file is automatically closed after creation.
// It returns an error if the file creation fails.
func (c *Celeritas) CreateFileIfNotExists(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				err = fmt.Errorf("error closing file: %w", err)
				// Handle or log the error appropriately
				fmt.Println(err)
			}
		}(file)
	}
	return nil
}

type Encryption struct {
	Key []byte
}

// Encrypt encrypts the given data using AES CFB mode encryption.
// It returns the encrypted data as a base64 encoded string.
// If an error occurs during encryption, it returns an empty string and the error.
func (e *Encryption) Encrypt(data string) (string, error) {
	plainText := []byte(data)
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts the given data using AES CFB mode decryption.
// It returns the decrypted data as a string.
// If an error occurs during decryption, it returns an empty string and the error.
func (e *Encryption) Decrypt(data string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
