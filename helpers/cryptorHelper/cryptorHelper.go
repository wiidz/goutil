package cryptorHelper

import (
	"crypto/md5"
	"fmt"
	"io"
)

func MD5Encrypt(raw string) string {
	w := md5.New()
	io.WriteString(w, raw)
	return fmt.Sprintf("%x", w.Sum(nil))
}
