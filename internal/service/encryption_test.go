package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tuananhlai/brevity-go/internal/service"
)

func TestNewEncryptionService_InvalidKey(t *testing.T) {
	key := make([]byte, 1)
	_, err := service.NewEncryptionService(key)
	assert.ErrorIs(t, err, service.ErrEncryptionServiceInvalidKeySize)
}

func TestEncryptDecrypt_Success(t *testing.T) {
	key := make([]byte, service.EncryptionServiceKeySize)

	s, err := service.NewEncryptionService(key)
	require.NoError(t, err)

	plainText := "Hello, World!"

	encryptedText := s.Encrypt([]byte(plainText))
	decryptedText, err := s.Decrypt(encryptedText)
	require.NoError(t, err)

	assert.Equal(t, plainText, string(decryptedText))
}
