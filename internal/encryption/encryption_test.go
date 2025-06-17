package encryption_test

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tuananhlai/brevity-go/internal/encryption"
)

func TestNewEncryptionService_InvalidKey(t *testing.T) {
	key := make([]byte, 1)
	_, err := encryption.New(key)
	require.ErrorIs(t, err, encryption.ErrInvalidKeySize)
}

func TestEncryptDecrypt_Success(t *testing.T) {
	key := make([]byte, encryption.KeySize)

	s, err := encryption.New(key)
	require.NoError(t, err)

	plainText := "Hello, World!"

	encryptedText := s.Encrypt([]byte(plainText))
	decryptedText, err := s.Decrypt(encryptedText)
	require.NoError(t, err)

	require.Equal(t, plainText, string(decryptedText))
}

func TestEncryptDecrypt_InvalidKey(t *testing.T) {
	key1 := make([]byte, encryption.KeySize)
	key2 := make([]byte, encryption.KeySize)
	_, err := rand.Read(key2)
	require.NoError(t, err)

	s1, err := encryption.New(key1)
	require.NoError(t, err)

	s2, err := encryption.New(key2)
	require.NoError(t, err)

	plainText := "Hello, World!"
	encryptedText := s1.Encrypt([]byte(plainText))

	_, err = s2.Decrypt(encryptedText)
	require.Error(t, err)
}

// TestEncrypt_DifferentCipherTextEveryTime tests that the encrypted text is different every time `Encrypt` is called.
// This is to ensure that the encrypted text is not predictable.
func TestEncrypt_DifferentCipherTextEveryTime(t *testing.T) {
	key := make([]byte, encryption.KeySize)
	s, err := encryption.New(key)
	require.NoError(t, err)

	plainText := "Hello, World!"
	encryptedText1 := s.Encrypt([]byte(plainText))
	encryptedText2 := s.Encrypt([]byte(plainText))

	assert.NotEqual(t, encryptedText1, encryptedText2)
}
