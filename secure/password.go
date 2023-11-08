// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package secure

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SALT_CHARS                = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	DEFAULT_PBKDF2_ITERATIONS = 310_000
)

// code taken from werkzeug
// 	https://github.com/pallets/werkzeug/blob/master/src/werkzeug/security.py

func generateSalt(length int) []byte {
	if length <= 0 {
		length = 16
	}
	var salt = make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		log.Panicf("rand.Read: %v", err)
	}
	return salt
}

// GeneratePasswordHash
//
//	Hash a password with the given method and salt with a string of
//	the given length. The format of the string returned includes the method
//	that was used so that :func:`check_password_hash` can check the hash.
//
//	The format for the hashed string looks like this::
//
//	method$salt$hash
func GeneratePasswordHash(password, method string) string {
	if method == "" {
		method = "default"
	}
	var saltText, passwdText string
	switch method {
	case "plain":
		passwdText = password

	case "default", "pbkdf2:sha256":
		var salt = generateSalt(32)
		var dk = pbkdf2.Key([]byte(password), salt, DEFAULT_PBKDF2_ITERATIONS, 32, sha256.New)
		saltText = hex.EncodeToString(salt)
		passwdText = hex.EncodeToString(dk)
	}
	return fmt.Sprintf("%s$%s$%s", method, saltText, passwdText)
}

// VerifyPasswordHash
//
//	check a password against a given salted and hashed password value.
//	In order to support unsalted legacy passwords this method supports
//	plain text passwords, md5 and sha1 hashes (both salted and unsalted).
func VerifyPasswordHash(hashText, password string) bool {
	var idx = strings.Index(hashText, "$")
	if idx <= 0 {
		return false
	}
	var method = hashText[:idx]
	hashText = hashText[idx+1:]
	idx = strings.Index(hashText, "$")
	if idx < 0 {
		return false
	}
	var saltText = hashText[:idx]
	hashText = hashText[idx+1:]

	switch method {
	case "plain":
		return hashText == password

	case "default", "pbkdf2:sha256":
		salt, err := hex.DecodeString(saltText)
		if err != nil {
			return false
		}
		var dk = pbkdf2.Key([]byte(password), salt, DEFAULT_PBKDF2_ITERATIONS, 32, sha256.New)
		return hashText == hex.EncodeToString(dk)
	}
	return false
}

// SignAccessToken 注册签名
func SignAccessToken(node, gameId, key string) string {
	var buf bytes.Buffer
	buf.WriteString(node)
	buf.WriteString(gameId)
	h := hmac.New(sha256.New, []byte(key))
	h.Write(buf.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

// SignEncryptSignature 签名
func SignEncryptSignature(method string, encrypt AESCryptor, priKey *rsa.PrivateKey) ([]byte, error) {
	if method == "" {
		return nil, nil
	}
	key, iv := encrypt.Key()
	var hash = sha256.New()
	hash.Write([]byte(method))
	hash.Write(key)
	hash.Write(iv)
	var digest = hash.Sum(nil)
	return rsa.SignPSS(rand.Reader, priKey, crypto.SHA256, digest, nil)
}

// VerifyEncryptSignature 校验签名
func VerifyEncryptSignature(method string, signature []byte, encrypt AESCryptor, pubKey *rsa.PublicKey) error {
	if method == "" {
		return nil
	}
	key, iv := encrypt.Key()
	var hash = sha256.New()
	hash.Write([]byte(method))
	hash.Write(key)
	hash.Write(iv)
	var digest = hash.Sum(nil)
	return rsa.VerifyPSS(pubKey, crypto.SHA256, digest, signature, nil)
}
