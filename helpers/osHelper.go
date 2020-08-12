package goutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type OsHelper struct{}

/**
 * @func: ExistFile 判断文件是否已存在
 * @author Wiidz
 * @date   2019-11-16
 */
func (*OsHelper) ExistFile(filename string, filesize int64) bool {
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

/**
 * @func: DownloadFile 从url中下载文件
 * @author Wiidz
 * @date   2019-11-16
 */
func (osHelper *OsHelper) DownloadFile(url string, localPath string, fb func(length, downLen int64)) error {
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
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	if osHelper.ExistFile(localPath, fsize) {
		return err
	}
	fmt.Println("fsize", fsize)
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
		fb(fsize, written)
	}
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		fmt.Println(err)
	}
	return err
}

//
///**
// * @func: OpenImageFile 打开图像文件
// * @author Wiidz
// * @date   2019-11-16
// */
//func OpenImageFile(localUri string) (image.Image, error) {
//	var m image.Image
//	ff, _ := ioutil.ReadFile(localUri) //读取文件 要先下载
//	bbb := bytes.NewBuffer(ff)
//
//	datatype, err := imgtype.Get(localUri)
//
//	if err != nil {
//		fmt.Println(err)
//		return m, err
//	}
//	fmt.Println("【datatype】", datatype)
//
//	switch datatype {
//	case "image/jpeg":
//		m, err = jpeg.Decode(bbb)
//	case "image/png":
//		m, err = png.Decode(bbb)
//	case "image/gif":
//		m, err = gif.Decode(bbb)
//	default:
//		fmt.Println("不支持的格式", reflect.TypeOf(datatype).String())
//	}
//	return m, nil
//}
//
//func Buff2Image(bytes []byte) {
//	_ = ioutil.WriteFile("/tmp/test.jpg", bytes, 0666)
//}

/**
 * @func: ReadJsonFile 读取json格式的文件
 * @author Wiidz
 * @date   2019-11-16
 */
func (*OsHelper) ReadJsonFile(filePath string) map[string]interface{} {
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

func (*OsHelper) GetFileBuf(uri string) []byte {
	buf, _ := ioutil.ReadFile(uri)
	return buf
}
