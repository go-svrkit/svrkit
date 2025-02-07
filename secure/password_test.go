// Copyright © Johnnie Chen ( qi7chen@github ). All rights reserved.
// See accompanying LICENSE file

package secure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestGeneratePasswordHash(t *testing.T) {
	var methods = []string{"plain", "default"}
	for _, method := range methods {
		for i := 0; i < 20; i++ {
			var password = randString(12)
			var hashText = GeneratePasswordHash(password, method)
			var ok = VerifyPasswordHash(hashText, password)
			if !ok {
				t.Fatalf("password mismatch: %s, %s", password, hashText)
			}
		}
	}
}

func BenchmarkGeneratePasswordHash(b *testing.B) {
	b.StopTimer()
	var password = randString(12)
	b.StartTimer()
	var hashText = GeneratePasswordHash(password, "default")
	var ok = VerifyPasswordHash(hashText, password)
	if !ok {
		b.Fatalf("password mismatch: %s, %s", password, hashText)
	}
}

func TestVerifyEncryptSignature(t *testing.T) {
	var method = "aes-192-cfb"
	aesCrypt, err := CreateAESCryptor(method)
	assert.Nil(t, err)
	prikey, err := LoadRSAPrivateKey(RSATestPrivateKey)
	assert.Nil(t, err)
	signature, err := SignEncryptSignature(method, aesCrypt, prikey)
	assert.Nil(t, err)
	assert.True(t, len(signature) > 0)

	pubkey, err := LoadRSAPublicKey(RSATestPublicKey)
	assert.Nil(t, err)
	err = VerifyEncryptSignature(method, signature, aesCrypt, pubkey)
	assert.Nil(t, err)
}
