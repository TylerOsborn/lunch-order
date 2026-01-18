package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

// GetEncryptionKey retrieves the encryption key from the environment
func GetEncryptionKey() ([]byte, error) {
	keyHex := os.Getenv("DATA_ENCRYPTION_KEY")
	if keyHex == "" {
		return nil, errors.New("DATA_ENCRYPTION_KEY environment variable is not set")
	}
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, errors.New("DATA_ENCRYPTION_KEY must be a valid hex string")
	}
	if len(key) != 32 {
		return nil, errors.New("DATA_ENCRYPTION_KEY must be 32 bytes (64 hex characters) for AES-256")
	}
	return key, nil
}

// Encrypt encrypts plaintext using AES-GCM
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func Decrypt(ciphertextBase64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Hash computes the HMAC-SHA256 of the input
func Hash(input string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
