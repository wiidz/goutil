package cryptorHelper

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
)

const (
	PemBegin = "-----BEGIN RSA PRIVATE KEY-----\n"
	PemEnd   = "\n-----END RSA PRIVATE KEY-----"
)

// RsaSign 使用rsa进行签名(RsaSign(queryStr, PrivateKey, crypto.SHA256))
func RsaSign(signContent string, privateKey string, hash crypto.Hash) (signStr string, err error) {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		return
	}
	signStr = base64.StdEncoding.EncodeToString(signature)
	return
}

func ParsePrivateKey(privateKey string) (priKey *rsa.PrivateKey, err error) {

	privateKey = FormatPrivateKey(privateKey)

	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, errors.New("私钥信息错误！")
	}

	// 3、解析DER编码的私钥，生成私钥对象
	priKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	return
}

func FormatPrivateKey(privateKey string) string {
	if !strings.HasPrefix(privateKey, PemBegin) {
		privateKey = PemBegin + privateKey
	}
	if !strings.HasSuffix(privateKey, PemEnd) {
		privateKey = privateKey + PemEnd
	}
	return privateKey
}
