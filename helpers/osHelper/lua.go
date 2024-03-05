package osHelper

import (
	"os"
)

// ReadByteFromFile 从文件中读取byte数据
func ReadByteFromFile(filePath string) (byteData []byte, err error) {
	file, _ := os.Open(filePath)
	defer file.Close()
	byteData, err = os.ReadFile(filePath)
	return
}
