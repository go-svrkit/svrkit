// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.txt

package secure

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

// https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
type AESCryptor interface {
	Key() ([]byte, []byte)
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

func GetEncryptionMethod(mode string, keySize int) string {
	return fmt.Sprintf("aes-%d-%s", keySize*8, strings.ToLower(mode))
}

// `method` is like "aes-128-cbc", "aes-256-ctr", "aes-192-gcm"
func CreateAESCryptor(method string) (AESCryptor, error) {
	var parts = strings.Split(method, "-")
	if strings.ToUpper(parts[0]) != "AES" {
		return nil, fmt.Errorf("AES encryption only, got %s", method)
	}
	blockSize, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	var n = blockSize / 8
	switch n {
	case 16, 24, 32:
	default:
		return nil, fmt.Errorf("invalid AES block size %d", blockSize)
	}
	var key = make([]byte, n)
	var iv = make([]byte, aes.BlockSize)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return NewAESCryptor(parts[2], key, iv)
}

func NewAESCryptor(mode string, key, iv []byte) (AESCryptor, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	var cryptor AESCryptor
	switch strings.ToUpper(mode) {
	case "CBC":
		if len(iv) != block.BlockSize() {
			return nil, fmt.Errorf("IV length must equal block size")
		}
		cryptor = &aesBlockCryptor{
			key:          key,
			iv:           iv,
			encryptBlock: cipher.NewCBCEncrypter(block, iv),
			decryptBlock: cipher.NewCBCDecrypter(block, iv),
		}
		return cryptor, nil

	case "CFB":
		if len(iv) != block.BlockSize() {
			return nil, fmt.Errorf("IV length must equal block size")
		}
		cryptor = &aesStreamCryptor{
			key:           key,
			iv:            iv,
			encryptStream: cipher.NewCFBEncrypter(block, iv),
			decryptStream: cipher.NewCFBDecrypter(block, iv),
		}
		return cryptor, nil

	case "CTR":
		if len(iv) != block.BlockSize() {
			return nil, fmt.Errorf("IV length must equal block size")
		}
		cryptor = &aesStreamCryptor{
			key:           key,
			iv:            iv,
			encryptStream: cipher.NewCTR(block, iv),
			decryptStream: cipher.NewCTR(block, iv),
		}
		return cryptor, nil

	case "GCM":
		encrypt, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		decrypt, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		cryptor = &aesGCMCryptor{
			key:        key,
			iv:         iv[:12],
			encryptGCM: encrypt,
			decryptGCM: decrypt,
		}
		return cryptor, nil

	case "OFB":
		if len(iv) != block.BlockSize() {
			return nil, fmt.Errorf("IV length must equal block size")
		}
		cryptor = &aesStreamCryptor{
			key:           key,
			iv:            iv,
			encryptStream: cipher.NewOFB(block, iv),
			decryptStream: cipher.NewOFB(block, iv),
		}
		return cryptor, nil
	}
	return nil, fmt.Errorf("invalid AES mode %s", mode)
}

type aesStreamCryptor struct {
	key, iv       []byte
	encryptStream cipher.Stream
	decryptStream cipher.Stream
}

func (c *aesStreamCryptor) Key() ([]byte, []byte) {
	return c.key, c.iv
}

func (c *aesStreamCryptor) Encrypt(plainData []byte) ([]byte, error) {
	var encrypted = make([]byte, len(plainData))
	c.encryptStream.XORKeyStream(encrypted, plainData)
	return encrypted, nil
}

func (c *aesStreamCryptor) Decrypt(encrypted []byte) ([]byte, error) {
	var decrypted = make([]byte, len(encrypted))
	c.decryptStream.XORKeyStream(decrypted, encrypted)
	return decrypted, nil
}

type aesBlockCryptor struct {
	key, iv      []byte
	encryptBlock cipher.BlockMode
	decryptBlock cipher.BlockMode
}

func (c *aesBlockCryptor) Key() ([]byte, []byte) {
	return c.key, c.iv
}

func (c *aesBlockCryptor) Encrypt(plainData []byte) ([]byte, error) {
	var blockSize = c.encryptBlock.BlockSize()
	var paddedData = PKCS5Pad(plainData, blockSize)
	var encrypted = make([]byte, len(paddedData))
	c.encryptBlock.CryptBlocks(encrypted, paddedData)
	return encrypted, nil
}

func (c *aesBlockCryptor) Decrypt(encrypted []byte) ([]byte, error) {
	var decrypted = make([]byte, len(encrypted))
	c.decryptBlock.CryptBlocks(decrypted, encrypted)
	return PKCS5Unpad(decrypted), nil
}

func PKCS5Pad(ciphertext []byte, blockSize int) []byte {
	var padding = blockSize - len(ciphertext)%blockSize
	var padtext = bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Unpad(encrypted []byte) []byte {
	var padding = encrypted[len(encrypted)-1]
	return encrypted[:len(encrypted)-int(padding)]
}

type aesGCMCryptor struct {
	key, iv    []byte
	encryptGCM cipher.AEAD
	decryptGCM cipher.AEAD
}

func (c *aesGCMCryptor) Key() ([]byte, []byte) {
	return c.key, c.iv
}

func (c *aesGCMCryptor) Encrypt(plainData []byte) ([]byte, error) {
	var encrypted = c.encryptGCM.Seal(nil, c.iv, plainData, nil)
	return encrypted, nil
}

func (c *aesGCMCryptor) Decrypt(encrypted []byte) ([]byte, error) {
	return c.decryptGCM.Open(nil, c.iv, encrypted, nil)
}
