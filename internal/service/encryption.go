package service

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const EncryptionServiceKeySize = 32

var EncryptionServiceKeySizeError = fmt.Errorf("invalid aes key size. expecting key to have length: %d", EncryptionServiceKeySize)

type EncryptionService interface {
	Encrypt(plainText []byte) []byte
	Decrypt(cipherText []byte) ([]byte, error)
}

type encryptionServiceImpl struct {
	gcm cipher.AEAD
}

// NewEncryptionService creates a new instance of EncryptionService. The `key` argument must be a byte slice
// with length of `EncryptionServiceKeySize`.
func NewEncryptionService(key []byte) (*encryptionServiceImpl, error) {
	if len(key) != EncryptionServiceKeySize {
		return nil, EncryptionServiceKeySizeError
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}

	return &encryptionServiceImpl{
		gcm: gcm,
	}, nil
}

func (s *encryptionServiceImpl) Encrypt(plainText []byte) []byte {
	// nonce does not need to be passed, as the AEAD returned by `NewGCMWithRandomNonce` already includes a random nonce.
	return s.gcm.Seal(nil, nil, plainText, nil)
}

func (s *encryptionServiceImpl) Decrypt(cipherText []byte) ([]byte, error) {
	// nonce does not need to be passed, as the AEAD returned by `NewGCMWithRandomNonce` will automatically
	// extract the nonce from the cipher text.
	return s.gcm.Open(nil, nil, cipherText, nil)
}
