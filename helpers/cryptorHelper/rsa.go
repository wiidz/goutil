package cryptorHelper

import "io"
import  "github.com/forgoer/openssl"

func MakeSign() {
	openssl.RSAGenerateKey(bits int, out io.Writer)
	openssl.RSAGeneratePublicKey(priKey []byte, out io.Writer)

	openssl.RSAEncrypt(src, pubKey []byte) ([]byte, error)
	openssl.RSADecrypt(src, priKey []byte) ([]byte, error)

	openssl.RSASign(src []byte, priKey []byte) ([]byte, error)
	openssl.RSAVerify(src, sign, pubKey []byte) error
}
