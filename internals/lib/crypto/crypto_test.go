package crypto

import (
	"crypto/rand"
	"io"
	"log"
	"testing"
)

func TestEncryptionDecryption(t *testing.T) {
	tests := []string{
		"Jane Doe",
		"CyberShiz",
		"mongol1999",
		"Accell_l",
	}

	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		t.Fatal("failed to generate encryption key")
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			ciphertext, err := Encrypt([]byte(tt), key)
			if err != nil {
				log.Fatalf("failed to encrypt: %v", err)
			}

			t.Logf("successfully encrypted %s -> %x", tt, ciphertext)

			plaintext, err := Decrypt(ciphertext, key)
			if err != nil {
				t.Fatalf("failed to decrypt: %v", err)
			}

			dec := string(plaintext)
			if dec != tt {
				t.Errorf("base and decrypted data are not equal: %s != %s", tt, dec)
			} else {
				t.Logf("successfully decrypted data: %s == %s", tt, dec)
			}
		})
	}
}
