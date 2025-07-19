package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(plaintext, key []byte) ([]byte, error) {
	const fn = "internals.lib.crypto.encrypt"

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return append(nonce, gcm.Seal(nil, nonce, plaintext, nil)...), nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	const fn = "internals.lib.crypto.decrypt"

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if gcm.NonceSize() > len(ciphertext) {
		return nil, fmt.Errorf("too short ciphertext")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
