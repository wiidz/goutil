package cryptorHelper

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// MD5Encrypt MD5加密
func MD5Encrypt(raw string) string {
	w := md5.New()
	_, _ = io.WriteString(w, raw)
	return fmt.Sprintf("%x", w.Sum(nil))
}

// SHA1EncryptOld 对字符串进行SHA1哈希
func SHA1EncryptOld(data string) string {
	t := sha1.New()
	_, _ = io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

// SHA1Hash 对字符串进行SHA1哈希
func SHA1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SHA256Hash 对字符串进行SHA256哈希
func SHA256Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
