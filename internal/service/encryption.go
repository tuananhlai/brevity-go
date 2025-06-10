package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// EncryptionConfig holds configuration for the encryption service.
type EncryptionConfig struct {
	EncryptionKey string // Base64 encoded or raw 32-byte key for AES-256
}

type EncryptionService interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

type encryptionServiceImpl struct {
	key []byte
}

// NewService creates a new encryption service.
// It takes a raw 32-byte key (for AES-256).
func NewService(cfg EncryptionConfig) (EncryptionService, error) {
	if cfg.EncryptionKey == "" {
		return nil, errors.New("encryption key cannot be empty")
	}

	// Assuming the key is provided as a raw 32-byte string for simplicity
	// In a real app, it might be base64 encoded, or derived from a passphrase.
	keyBytes := []byte(cfg.EncryptionKey)

	if len(keyBytes) != 32 { // AES-256 requires a 32-byte key
		return nil, fmt.Errorf("encryption key must be 32 bytes long for AES-256, got %d bytes", len(keyBytes))
	}

	return &encryptionServiceImpl{
		key: keyBytes,
	}, nil
}

// Encrypt encrypts plaintext using AES-256 GCM.
// The nonce is prepended to the ciphertext.
func (s *encryptionServiceImpl) Encrypt(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Never use more than 2^32 random nonces with a given key because of the risk of a repeat.
	// A 12-byte nonce is standard for GCM.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Seal encrypts and authenticates the plaintext, then appends the authentication tag.
	// The `nil` first argument means Seal will allocate a new slice for the result.
	// The `nonce` is prepended to the ciphertext for easy storage and retrieval during decryption.
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)

	return cipherText, nil
}

// Decrypt decrypts ciphertext using AES-256 GCM.
// It expects the nonce to be prepended to the ciphertext.
func (s *encryptionServiceImpl) Decrypt(cipherTextWithNonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherTextWithNonce) < nonceSize {
		return nil, errors.New("ciphertext is too short to contain nonce")
	}

	nonce, cipherText := cipherTextWithNonce[:nonceSize], cipherTextWithNonce[nonceSize:]

	// Open decrypts and authenticates the ciphertext.
	// If the data has been tampered with, or the key/nonce is incorrect, it will return an error.
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		// This error typically indicates tampering or incorrect key/nonce.
		return nil, fmt.Errorf("failed to decrypt or authenticate data: %w", err)
	}

	return plainText, nil
}
