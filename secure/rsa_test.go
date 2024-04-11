// Copyright © Johnnie Chen ( ki7chen@github ). All rights reserved.
// See accompanying LICENSE file

package secure

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

var RSATestPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0wugAtwSpvDgiYKi6GC5390KY
Qy4bAC2jBO13zVW5aQ83WPUHyvhXnj1N1xujGHMJyNGwEYA9voxmPxyYn83D4cRM
Bga/GaJtLzbJwakpFMaEzUtIq8bCgPSTXtxuUx+spw6G/yl6MxO9O+RhScDrQPmp
jvB4Z/u0Dl5tdwJPqQIDAQAB
-----END PUBLIC KEY-----`)

var RSATestPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQC0wugAtwSpvDgiYKi6GC5390KYQy4bAC2jBO13zVW5aQ83WPUH
yvhXnj1N1xujGHMJyNGwEYA9voxmPxyYn83D4cRMBga/GaJtLzbJwakpFMaEzUtI
q8bCgPSTXtxuUx+spw6G/yl6MxO9O+RhScDrQPmpjvB4Z/u0Dl5tdwJPqQIDAQAB
AoGATEj9NGAIvcFLR2bXjkHqSoK1PiEL8iUvHV9VAHxNs0PdQhRuxG0qRX/oi1M+
vKPy2KxBojagkm46PmRgIyE96rkI94boLKfctuMVsqg22GQDtcvBuSVrYPNfgDLw
1EbzQihFqgxO/QYnuakn7GAE4N9x1R5gAQr7Wy00aekhHkkCQQDlZgzAYAyFtZA4
A6NOGGPVM8/FLYwUZVyb9jh1uXJiOEj1j7p5bJhUrXRRduJ+Z2t4OP993OprTV86
slO/QVkvAkEAybkBY2JIK+nDxdxCEmbMcQRolTL/l/MQayBF0lbOVHb5svDdpWbm
q9Y6PwfVK8jbp8bJWYovDJ2wEQF3d0R+pwJAH+wzmhHDrFe32hOnhhaezeyH3UiZ
Vb1FRe7drIRCBqkOfh2iNYOHL0F0DmIc4rpBmllUNI+pj4UU23Y1cUgGwQJALBhq
+0Siria9iuTo9IjQK+xgyCyLvrV9Y018tcwP8lrHnpwUd3GU/v8nYFvf92BC09wa
a55PRpy5vh3p9YJdhQJBAKtEEaC8EB7ghXvvo1O+MJotd+EqO330JsLTUf0GcsjI
zV/yJu951ELuzMZTfemh6l8stjjDYlRvZVPbjwrZP8g=
-----END RSA PRIVATE KEY-----`)

// Run command below to generate test key files:
// 	openssl genrsa -out rsa_prikey.pem 1024
// 	openssl rsa -in rsa_prikey.pem -pubout -out rsa_pubkey.pem

func TestRSADecrypt(t *testing.T) {
	prikey, err := LoadRSAPrivateKey(RSATestPrivateKey)
	if err != nil {
		t.Fatalf("load private key: %v", err)
	}
	pubkey, err := LoadRSAPublicKey(RSATestPublicKey)
	if err != nil {
		t.Fatalf("load public key: %v", err)
	}
	var maxSize = MaxEncryptSize(pubkey)
	var data = []byte("a quick brown fox jumps over the lazy dog")
	encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubkey, data, nil)
	if err != nil {
		t.Fatalf("RSAEncrypt: %v, %d/%d", err, len(data), maxSize)
	}
	decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, prikey, encrypted, nil)
	if err != nil {
		t.Fatalf("RSADecrypt: %v", err)
	}
	if !bytes.Equal(data, decrypted) {
		t.Fatalf("data not equal after encryption/decription")
	}
	t.Logf("RSA encryption OK")
}
