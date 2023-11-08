// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package secure

import (
	"testing"
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
	var hashText = GeneratePasswordHash(password, "")
	var ok = VerifyPasswordHash(hashText, password)
	if !ok {
		b.Fatalf("password mismatch: %s, %s", password, hashText)
	}
}
