package cryptorHelper

import (
	"crypto/md5"
	"crypto/sha1"
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

// SHA1Encrypt 对字符串进行SHA1哈希
func SHA1Encrypt(data string) (result string, err error) {
	hasher := sha1.New()
	_, err = io.WriteString(hasher, data)
	if err != nil {
		// 处理错误
		return
	}
	hashBytes := hasher.Sum(nil)
	result = hex.EncodeToString(hashBytes)
	return
}
