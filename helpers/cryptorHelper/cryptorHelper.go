package cryptorHelper

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"io"
)

// MD5Encrypt MD5加密
func MD5Encrypt(raw string) string {
	w := md5.New()
	_, _ = io.WriteString(w, raw)
	return fmt.Sprintf("%x", w.Sum(nil))
}

// SHA1Encrypt 对字符串进行SHA1哈希
func SHA1Encrypt(data string) string {
	t := sha1.New()
	_, _ = io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}
