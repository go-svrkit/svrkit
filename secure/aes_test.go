// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package secure

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func testSimpleEncryption(t *testing.T, method string, cryptor AESCryptor) {
	var size = 32 + rand.Intn(2048)
	var plainData = []byte(randString(size))
	encrypted, err := cryptor.Encrypt(plainData)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	decrypted, err := cryptor.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if !bytes.Equal(decrypted, plainData) {
		t.Fatalf("decrypt data mismatch: %s", method)
	}
}

func testEncryption(t *testing.T, method string, cryptor AESCryptor) {
	var plainTextList = make([][]byte, 100)
	var encryptedList = make([][]byte, 100)
	for i := 0; i < 100; i++ {
		var size = 32 + rand.Intn(2048)
		var data = []byte(randString(size))
		plainTextList[i] = data
		encrypted, err := cryptor.Encrypt(data)
		if err != nil {
			t.Fatalf("encrypt failed: %v", err)
		}
		encryptedList[i] = encrypted
	}
	for i := 0; i < 100; i++ {
		decrypted, err := cryptor.Decrypt(encryptedList[i])
		if err != nil {
			t.Fatalf("decrypt failed: %v", err)
		}
		if !bytes.Equal(decrypted, plainTextList[i]) {
			t.Fatalf("decrypt data mismatch: %s", method)
		}
	}
}

func TestAESEncryption(t *testing.T) {
	var modes = []string{"cbc", "cfb", "ctr", "gcm", "ofb"}
	var sizes = []int{128, 192, 256}
	for _, mode := range modes {
		for _, size := range sizes {
			var method = fmt.Sprintf("aes-%d-%s", size, mode)
			cryptor, err := CreateAESCryptor(method)
			if err != nil {
				t.Fatalf("create cryptor failed: %v", err)
			}
			testSimpleEncryption(t, method, cryptor)
			testEncryption(t, method, cryptor)
		}
	}
}

func benchAESCryptor(b *testing.B, cryptor AESCryptor) {
	var size = 64 + rand.Intn(2048)
	var plainData = []byte(randString(size))
	for i := 0; i < b.N; i++ {
		encrypted, err := cryptor.Encrypt(plainData)
		if err != nil {
			b.Fatalf("encrypt failed: %v", err)
		}
		_, err = cryptor.Decrypt(encrypted)
		if err != nil {
			b.Fatalf("decrypt failed: %v", err)
		}
	}
}

func BenchmarkAESCtr(b *testing.B) {
	cryptor, err := CreateAESCryptor("aes-192-ctr")
	if err != nil {
		b.Fatalf("create cryptor failed: %v", err)
	}
	benchAESCryptor(b, cryptor)
}

func BenchmarkAESGCM(b *testing.B) {
	cryptor, err := CreateAESCryptor("aes-192-gcm")
	if err != nil {
		b.Fatalf("create cryptor failed: %v", err)
	}
	benchAESCryptor(b, cryptor)
}
