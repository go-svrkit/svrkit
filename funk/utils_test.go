// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package funk

import (
	"fmt"
	"testing"
)

func TestMD5Sum(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "5d41402abc4b2a76b9719d911017c592"},
		{[]byte("world"), "7d793037a0760186574b0282f2f435e7"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := MD5Sum(tt.input); got != tt.want {
				t.Errorf("MD5Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1Sum(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{[]byte("world"), "7c211433f02071597741e6ff5a8ea34789abbf43"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := SHA1Sum(tt.input); got != tt.want {
				t.Errorf("SHA1Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256Sum(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}
	for i, tt := range tests {
		var name = fmt.Sprintf("case-%d", i+1)
		t.Run(name, func(t *testing.T) {
			if got := SHA256Sum(tt.input); got != tt.want {
				t.Errorf("SHA1Sum() = %v, want %v", got, tt.want)
			}
		})
	}
}
