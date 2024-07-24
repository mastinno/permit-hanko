package aes_gcm

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestNewEncryptionKey(t *testing.T) {
	key1 := newEncryptionKey()
	key2 := newEncryptionKey()
	assert.Equal(t, len(key1), 32)
	assert.Equal(t, len(key2), 32)
	assert.NotEqualf(t, key1, key2, "two separate constructed keys should not be equal.")
}

func TestNewAESGCM(t *testing.T) {
	for k, c := range []struct {
		keys  []string
		check func(aesgcm *AESGCM, err error)
	}{
		{
			keys: []string{},
			check: func(aesgcm *AESGCM, err error) {
				assert.Error(t, err, "empty key list should get rejected.")
				assert.Nil(t, aesgcm)
			},
		},
		{
			keys: []string{"too-short"},
			check: func(aesgcm *AESGCM, err error) {
				assert.Error(t, err, "too short key should get rejected.")
				assert.Nil(t, aesgcm)
			},
		},
		{
			keys: []string{string(newEncryptionKey()[:]), "too-short"},
			check: func(aesgcm *AESGCM, err error) {
				assert.Error(t, err, "too short key in any position should get rejected.")
				assert.Nil(t, aesgcm)
			},
		},
		{
			keys: []string{string(newEncryptionKey()[:])},
			check: func(aesgcm *AESGCM, err error) {
				assert.NoError(t, err, "Generated Key should be accepted")
				assert.NotNil(t, aesgcm)
			},
		},
		{
			keys: []string{string(newEncryptionKey()[:]), string(newEncryptionKey()[:])},
			check: func(aesgcm *AESGCM, err error) {
				assert.NoError(t, err, "two generated keys should be accepted")
				assert.NotNil(t, aesgcm)
			},
		},
	} {
		t.Run(fmt.Sprintf("case=%d", k), func(t *testing.T) {
			c.check(NewAESGCM(c.keys))
		})
	}
}

func TestAESGCM_EncryptDecrypt(t *testing.T) {
	// Encrypt
	plaintext := "testTesttestTestTestTEST"
	aesgcm, err := NewAESGCM([]string{string(newEncryptionKey()[:])})
	assert.NoError(t, err)
	assert.NotNil(t, aesgcm)
	ciphertext, err := aesgcm.Encrypt([]byte(plaintext))
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)
	//Decrypt
	plainAgain, err := aesgcm.Decrypt(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, string(plainAgain), plaintext)
}

func TestAESGCM_Decrypt(t *testing.T) {
	plaintext := "testTesttestTestTestTEST"
	key := []string{"superSecureAndRandomlyCreatedKey"}
	ciphertext := "1SiZZlBRnNM8P9xzppDI4n8YYVMG4Hf9UKNANqmxXtsxpHPt3fojlxRmwwOhkfQSCPooOA=="
	aesgcm, err := NewAESGCM(key)
	assert.NoError(t, err)
	assert.NotNil(t, aesgcm)
	decrypted, err := aesgcm.Decrypt(ciphertext)
	assert.NoError(t, err)
	assert.Equal(t, []byte(plaintext), decrypted)
}

func TestAESGCM_SomeoneModifiedTheCiphertext(t *testing.T) {
	// Encrypt
	plaintext := "testTesttestTestTestTEST"
	aesgcm, err := NewAESGCM([]string{string(newEncryptionKey()[:])})
	assert.NoError(t, err)
	assert.NotNil(t, aesgcm)
	ciphertext, err := aesgcm.Encrypt([]byte(plaintext))
	assert.NoError(t, err)
	assert.NotEmpty(t, ciphertext)

	// Modify cipher
	cipher := []rune(ciphertext)
	cipher[35] = cipher[35] + 1

	//Try to decrypt
	plainAgain, err := aesgcm.Decrypt(string(cipher))
	assert.Error(t, err)
	assert.NotEqual(t, string(plainAgain), plaintext)
}

// newEncryptionKey generates a random 256-bit key. It panics if the source of randomness fails.
func newEncryptionKey() *[32]byte {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		panic(err)
	}
	return &key
}
