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

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return append(nonce, aesgcm.Seal(nil, nonce, plaintext, nil)...), nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	const fn = "internals.lib.crypto.decrypt"

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	if aesgcm.NonceSize() > len(ciphertext) {
		return nil, fmt.Errorf("too short ciphertext")
	}

	nonce, ciphertext := ciphertext[:aesgcm.NonceSize()], ciphertext[aesgcm.NonceSize():]

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
