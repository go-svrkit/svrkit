// Copyright © 2021 ichenq@gmail.com All rights reserved.
// See accompanying files LICENSE.

package secure

import (
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadRSAPublicKeyFile(pemFile string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	return LoadRSAPublicKey(data)
}

// LoadRSAPublicKey 解析公钥文件
func LoadRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("incorrect public key file")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("unexpected key type %s", block.Type)
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PublicKey), nil
}

func LoadRSAPrivateKeyFile(pemFile string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	return LoadRSAPrivateKey(data)
}

// LoadRSAPrivateKey 解析私钥文件
func LoadRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("incorrect private key file")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("unexpected key type %s", block.Type)
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// MaxEncryptSize 最大加密内容大小
func MaxEncryptSize(pubkey *rsa.PublicKey) int {
	var k = pubkey.Size()
	var hash = sha256.New()
	return k - 2*hash.Size() - 2
}
