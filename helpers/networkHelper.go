package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type NetworkHelper struct{}

type Method int8

const (
	Get    Method = 1
	Post   Method = 2
	Put    Method = 3
	Delete Method = 4
)

func (p Method) String() string {
	switch p {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Delete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

var typeHelper = TypeHelper{}
var strHelper = StrHelper{}

/**
 * @func: GetParsedURL 获取get地址
 * @author: Wiidz
 * @date: 2021-06-20
 */
func (networkHelper *NetworkHelper) GetParsedURL(apiURL string, params map[string]interface{}) (string, error) {

	//【1】解析URL
	var targetURL *url.URL
	targetURL, err := url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return "", err
	}

	//【2】构造参数
	param := url.Values{}
	for key, value := range params {
		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}
	targetURL.RawQuery = param.Encode() //如果参数中有中文参数,这个方法会进行URLEncode

	//【3】返回
	return targetURL.String(), err
}

/**
 * @func: GetRequest 发送get请求
 * @author: Wiidz
 * @date: 2021-06-20
 */
func (networkHelper *NetworkHelper) GetRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, http.Header, error) {

	//【1】解析URL
	var targetURL *url.URL
	targetURL, err := url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil, err
	}

	//【2】构造参数
	param := url.Values{}
	for key, value := range params {
		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}
	targetURL.RawQuery = param.Encode() //如果参数中有中文参数,这个方法会进行URLEncode
	log.Println("networkHelper.GetRequest:", targetURL)

	//【3】发送请求
	resp, err := http.Get(targetURL.String())
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	//【4】读取body体
	data, err := ioutil.ReadAll(resp.Body)
	var netReturn map[string]interface{}
	json.Unmarshal(data, &netReturn)

	//【5】返回
	return netReturn, resp.Header, err
}

/**
 * @func: PostRequest 发送post请求
 * @author: Wiidz
 * @date:  2021-6-20
 */
func (networkHelper *NetworkHelper) PostRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, http.Header, error) {

	//【1】解析URL
	var targetURL *url.URL
	targetURL, err := url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil, err
	}

	//【2】构造参数
	param := url.Values{}
	for key, value := range params {
		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}
	log.Println("networkHelper.PostRequest:", apiURL)
	log.Println("params:", param)

	//【3】发送请求
	resp, err := http.PostForm(targetURL.String(), param)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	//【4】读取body
	data, err := ioutil.ReadAll(resp.Body)
	var netReturn map[string]interface{}
	fmt.Println(data)
	json.Unmarshal(data, &netReturn)

	//【5】返回
	return netReturn, resp.Header, err
}

/**
 * @func: DownloadFile 下载目标文件到本地
 * @author: Wiidz
 * @date:  2021-6-20
 */
func (networkHelper *NetworkHelper) DownloadFile(targetURL, localPath string) (fileName, pathString string, err error) {

	if localPath == "" {
		localPath = "/tmp/download/"
	}

	//【1】下载文件
	resp, err := http.Get(targetURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	//【2】创建一个文件用于保存
	tempSlice := typeHelper.Explode(targetURL, ".")
	format := tempSlice[len(tempSlice)-1]
	fileName = strHelper.GetRandomString(10) + "." + format.(string)
	tempPath := localPath + fileName
	out, err := os.Create(tempPath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	//【3】然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", "", err
	}

	//【4】返回地址
	return fileName, tempPath, err
}

/**
 * @func: DownloadFileWithFormat 下载目标文件到本地，强制命名格式
 * @author: Wiidz
 * @date:  2021-6-20
 */
func (*NetworkHelper) DownloadFileWithFormat(targetURL, localPath, format string) (fileName, pathString string, err error) {

	if localPath == "" {
		localPath = "/tmp/download/"
	}

	//【1】下载文件
	resp, err := http.Get(targetURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	//【2】创建一个文件用于保存
	fileName = strHelper.GetRandomString(10) + "." + format
	tempPath := localPath + fileName
	out, err := os.Create(tempPath)
	if err != nil {
		return "", "", err
	}
	defer out.Close()

	//【3】然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", "", err
	}

	//【4】返回地址
	return fileName, tempPath, err
}


/**
 * @func: PostRequest 发送post请求
 * @author Wiidz
 * @date   2019-11-16
 */
func (*NetworkHelper) PostJsonRequest(apiURL string, params map[string]interface{}) ([]byte, error) {

	param := url.Values{}

	for key, value := range params {

		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}

	jsonByte, _ := json.Marshal(params)

	resp, err := http.Post(apiURL, "application/json", bytes.NewReader(jsonByte))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, e := ioutil.ReadAll(resp.Body)

	return data, e
}

func (*NetworkHelper) Request(method Method, targetURL string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, http.Header, int,error) {

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil,0, err
	}

	//【2】创建client
	client := &http.Client{}

	//【3】构造参数
	param := url.Values{}
	for key, value := range params {
		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}

	var body io.Reader
	if method == Get {
		parsedURL.RawQuery = param.Encode() //如果参数中有中文参数,这个方法会进行URLEncode
	} else {
		body = strings.NewReader(param.Encode())
	}

	//【4】创建请求
	request, err := http.NewRequest(method.String(), parsedURL.String(), body)
	if err != nil {
		panic(err)
	}

	//【5】增加header选项
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	log.Println("apiURL", targetURL, body)

	//【6】发送请求
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)
	log.Println("【resp.Body】",resp.Body)
	var netReturn map[string]interface{}
	json.Unmarshal(data, &netReturn)

	//【8】返回
	return netReturn, resp.Header,resp.StatusCode, err

}