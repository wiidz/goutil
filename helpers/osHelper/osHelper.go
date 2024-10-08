package osHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ExistSameNameFile 判断是否已存在同名文件
func ExistSameNameFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	return true
}

/**
 * @func: ExistFile 判断文件是否已存在
 * @author Wiidz
 * @date   2019-11-16
 */
func ExistSameSizeFile(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		fmt.Println(info)
		return false
	}
	if filesize == info.Size() {
		fmt.Println("文件已存在！", info.Name(), info.Size(), info.ModTime())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

// ReadJsonFile 读取json格式的文件
func ReadJsonFile(filePath string, targetStruct interface{}) error {
	file, _ := os.Open(filePath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(&targetStruct)
}

// ReadJsonFileMap 读取json格式的文件
func ReadJsonFileMap(filePath string) map[string]interface{} {
	file, _ := os.Open(filePath)
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := make(map[string]interface{}, 0)
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return conf
}

func GetFileBuf(uri string) []byte {
	buf, _ := ioutil.ReadFile(uri)
	return buf
}

func GetFileString(uri string) string {
	buf, _ := ioutil.ReadFile(uri)
	return string(buf)
}

// DownloadFileFromContext 从请求体中保存文件
func DownloadFileFromContext(ctx iris.Context, fieldName, targetPath string) (fileName, filePath string, err error) {

	// Get the file from the request.
	file, info, err := ctx.FormFile(fieldName)
	if err != nil {
		err = errors.New("上传文件为空")
		return
	}
	defer file.Close()

	fileName = typeHelper.Int64ToStr(time.Now().Unix()) + strHelper.GetRandomString(4) + "-" + info.Filename
	filePath = targetPath + fileName
	//创建一个具有相同名称的文件 假设你有一个名为'uploads'的文件夹
	// mkdir uploads
	// chomod -R 777 uploads

	out, err := os.OpenFile(filePath,
		os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		err = errors.New("下载时出错")
		return
	}
	defer out.Close()
	io.Copy(out, file)

	return
}

// IsDirExist 判断目录是否存在
// dirPath 绝对路径，不要以/结尾
func IsDirExist(dirPath string) bool {
	s, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
// perm：755,777
func CreateDir(dirName string, perm os.FileMode) error {
	return os.Mkdir(dirName, perm)
}

// CreateIfNotExist 如果目录不存在，则创建
func CreateIfNotExist(dirName string, perm os.FileMode) (err error) {
	exist := IsDirExist(dirName)
	if exist {
		return
	}

	return CreateDir(dirName, perm)
}

// Delete 删除指定路径文件夹/文件
func Delete(filePath string) error {
	return os.Remove(filePath)
}

// DeleteFiles 删除指定路径文件夹/文件
func DeleteFiles(filePaths []string) (err error) {
	for _, v := range filePaths {
		err = os.Remove(v)
		if err != nil {
			break
		}
	}
	return
}

// DownloadFile 从url中下载文件
func DownloadFile(url string, localPath string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	fmt.Println(tmpFilePath)
	//创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	//get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	//读取服务器返回的文件大小
	contentLen := resp.Header.Get("Content-Length")
	// 如果能获取文件长度 则判断是否存在
	if contentLen != "" {
		fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		if ExistSameSizeFile(localPath, fsize) {
			return err
		}
		fmt.Println("fsize", fsize)
	}

	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		if fb != nil {
			fb(fsize, written)
		}
	}
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		fmt.Println(err)
	}
	return err
}

// SimpleDownloadFile 简单下载（下载到当前工程目录的/tmp下，记得给权限）
func SimpleDownloadFile(url string) (finalFilePath string, err error) {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)

	dir, _ := os.Getwd()

	_, _, fileType := networkHelper.GetFileNameFromURL(url)

	tempName := strHelper.GetRandomString(8)

	tempPath := dir + "/temp/" + tempName + ".download"        // 临时路径
	finalFilePath = dir + "/temp/" + tempName + "." + fileType // 最终路径

	//创建一个http client
	client := new(http.Client)

	// client.Timeout = time.Second * 60 //设置超时时间
	// get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return
	}

	// 读取服务器返回的文件大小
	contentLen := resp.Header.Get("Content-Length")
	// 如果能获取文件长度 则判断是否存在
	if contentLen != "" {
		fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
		if err != nil {
			fmt.Println(err)
		}
		if ExistSameSizeFile(tempPath, fsize) {
			return
		}
		fmt.Println("fsize", fsize)
	}

	// 创建文件
	file, err := os.Create(tempPath)
	if err != nil {
		return
	}
	defer file.Close()
	if resp.Body == nil {
		err = errors.New("body is null")
		return
	}
	defer resp.Body.Close()

	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	if err == nil {
		file.Close()
		err = os.Rename(tempPath, finalFilePath)
		fmt.Println(err)
	}
	return
}

// GetFileNameFromLocalPath 从本地路径中提取文件
// localFilePath 应该是 /home/www/wwwroot/aaaaa.html 这样
func GetFileNameFromLocalPath(localFilePath string) (targetDir, wholeName, fileName, fileType string) {

	reg := regexp.MustCompile(`([^/]\w*\.*\w*)$`)
	result := reg.FindStringSubmatch(localFilePath)
	if len(result) == 2 {
		wholeName = result[1]
		temp := typeHelper.ExplodeStr(wholeName, ".")

		if len(temp) == 2 {

			fileName = temp[0]
			fileType = temp[1]

			targetDir = strings.ReplaceAll(localFilePath, "/"+wholeName, "")
		}
	}
	return
}

// SaveStringToFile 将字符串输出到文件中
func SaveStringToFile(content string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
