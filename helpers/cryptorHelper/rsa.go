package cryptorHelper

import "io"
import "github.com/forgoer/openssl"

// RSASign 使用 crypto.SHA256 创建rsa密钥
func RSASign(src []byte, priKey []byte) ([]byte, error) {
	return openssl.RSASign(src, priKey)
}

// RSAVerify 验签
func RSAVerify(src, sign, priKey []byte) error {
	return openssl.RSAVerify(src, sign, priKey)
}

// RSAEncrypt RSA加密
func RSAEncrypt(src []byte, priKey []byte) ([]byte, error) {
	return openssl.RSAEncrypt(src, priKey)
}

// RSADecrypt RSA解密
func RSADecrypt(src []byte, priKey []byte) ([]byte, error) {
	return openssl.RSADecrypt(src, priKey)
}

// RSAGenerateKey 创建 RSA private key 私钥
func RSAGenerateKey(bits int, out io.Writer) error {
	return openssl.RSAGenerateKey(bits, out)
}

// RSAGeneratePublicKey 创建 RSA public key 公钥
func RSAGeneratePublicKey(priKey []byte, out io.Writer) error {
	return openssl.RSAGeneratePublicKey(priKey, out)
}
