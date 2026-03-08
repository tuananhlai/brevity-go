package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

const KeySize = 32

var ErrInvalidKeySize = fmt.Errorf("invalid aes key size. expecting key to have length: %d", KeySize)

// Cipher is a wrapper around the AES-GCM cipher, providing a simple interface for encrypting and decrypting data.
type Cipher struct {
	gcm cipher.AEAD
}

// New creates a new instance of EncryptionService. The `key` argument must be a byte slice
// with length of `KeySize`.
func New(key []byte) (*Cipher, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}

	return &Cipher{
		gcm: gcm,
	}, nil
}

func (s *Cipher) Encrypt(plainText []byte) []byte {
	// nonce does not need to be passed, as the AEAD returned by `NewGCMWithRandomNonce` already includes a random nonce.
	return s.gcm.Seal(nil, nil, plainText, nil)
}

func (s *Cipher) Decrypt(cipherText []byte) ([]byte, error) {
	// nonce does not need to be passed, as the AEAD returned by `NewGCMWithRandomNonce` will automatically
	// extract the nonce from the cipher text.
	return s.gcm.Open(nil, nil, cipherText, nil)
}
