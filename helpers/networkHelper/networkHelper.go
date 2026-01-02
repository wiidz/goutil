package networkHelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/wiidz/goutil/helpers/sliceHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/validatorMng"
	"github.com/wiidz/goutil/structs/networkStruct"
)

/**
 * @func: GetParsedURL 获取get地址
 * @author: Wiidz
 * @date: 2021-06-20
 */
func GetParsedURL(apiURL string, params map[string]interface{}) (string, error) {

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
func GetRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, *http.Header, error) {

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
	return netReturn, &resp.Header, err
}

/**
 * @func: PostRequest 发送post请求
 * @author: Wiidz
 * @date:  2021-6-20
 */
func PostRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, *http.Header, error) {

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
	return netReturn, &resp.Header, err
}

/**
 * @func: DownloadFile 下载目标文件到本地
 * @author: Wiidz
 * @date:  2021-6-20
 */
func DownloadFile(targetURL, localPath string) (fileName, pathString string, err error) {

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
func DownloadFileWithFormat(targetURL, localPath, format string, headers map[string]string) (fileName, tempPath string, header *http.Header, err error) {

	if localPath == "" {
		localPath = "/tmp/download/"
	}
	//【1】解析URL
	//var parsedURL *url.URL
	//parsedURL, err = url.Parse(targetURL)
	//if err != nil {
	//	fmt.Printf("解析url错误:\r\n%v", err)
	//	return
	//}

	//【1】下载文件
	client := &http.Client{}
	var body io.Reader

	//【4】创建请求
	request, err := http.NewRequest("GET", targetURL, body)
	if err != nil {
		return
	}

	//【5】增加header选项
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	header = &resp.Header

	//【2】创建一个文件用于保存
	fileName = strHelper.GetRandomString(10) + "." + format
	tempPath = localPath + fileName
	out, err := os.Create(tempPath)
	if err != nil {
		return
	}
	defer out.Close()

	//【3】然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	return
}

/**
 * @func: PostRequest 发送post请求
 * @author Wiidz
 * @date   2019-11-16
 */
func PostJsonRequest(apiURL string, params map[string]interface{}) ([]byte, error) {

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

func RequestRaw(method networkStruct.Method, targetURL string, params map[string]interface{}, headers map[string]string) (string, *http.Header, int, error) {

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return "", nil, 0, err
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
	if method == networkStruct.Get {
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
	//【5-1】增加content-Length
	if method != networkStruct.Get {
		request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	}

	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return "", nil, 0, err
	}
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)

	//【8】返回
	return string(data), &resp.Header, resp.StatusCode, err

}

func RequestJson(method networkStruct.Method, targetURL string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, *http.Header, int, error) {

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil, 0, err
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
	if method == networkStruct.Get {
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

	//【5-1】增加content-Length
	if method != networkStruct.Get {
		request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	}

	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	b, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)

	var netReturn map[string]interface{}
	json.Unmarshal(data, &netReturn)

	//【8】返回
	return netReturn, &resp.Header, resp.StatusCode, err

}

func RequestJsonWithStruct(
	method networkStruct.Method,
	targetURL string,
	params interface{}, // 可结构体/Map
	headers map[string]string,
	iStruct interface{},
	debug bool, // 新增调试开关
) (interface{}, *http.Header, int, error) {
	if debug {
		fmt.Println("【method】", method.String())
		fmt.Println("【targetURL】", targetURL)
		fmt.Printf("【params】 %+v\n", params)
		fmt.Println("【headers】", headers)
	}

	// 1. 解析URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		if debug {
			fmt.Printf("解析url错误: %v\n", err)
		}
		return nil, nil, 0, err
	}

	var body io.Reader
	var reqBodyJson string

	if method == networkStruct.Get {
		// GET: params拼URL
		if params != nil {
			m := map[string]interface{}{}
			switch v := params.(type) {
			case map[string]interface{}:
				m = v
			default:
				b, _ := json.Marshal(v)
				json.Unmarshal(b, &m)
			}
			q := parsedURL.Query()
			for k, v := range m {
				// 建议可以用 ToString
				q.Set(k, fmt.Sprintf("%v", v))
			}
			parsedURL.RawQuery = q.Encode()
		}
	} else {
		// POST/PUT: params转json做body
		if params != nil {
			jsonBytes, err := json.Marshal(params)
			if err != nil {
				if debug {
					fmt.Println("【marshal params error】", err)
				}
				return nil, nil, 0, err
			}
			reqBodyJson = string(jsonBytes)
			body = bytes.NewReader(jsonBytes)
		}
	}

	if debug {
		fmt.Println("【parsedURL】", parsedURL.String())
		if method != networkStruct.Get {
			fmt.Println("【body】", body != nil)
			fmt.Println("【reqBodyJson】", reqBodyJson)
		}
	}

	// 构建请求
	request, err := http.NewRequest(method.String(), parsedURL.String(), body)
	if err != nil {
		if debug {
			fmt.Printf("创建http请求失败: %v\n", err)
		}
		return nil, nil, 0, err
	}
	// 设置headers
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	if method != networkStruct.Get && params != nil {
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		if debug {
			fmt.Printf("http请求失败: %v\n", err)
		}
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	// 读取返回结果
	resBytes, err := io.ReadAll(resp.Body)
	if debug {
		fmt.Println("【StatusCode】", resp.StatusCode)
		fmt.Println("【resStr】", string(resBytes))
	}
	if err != nil {
		return nil, &resp.Header, resp.StatusCode, err
	}

	// 反序列化
	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json2.Unmarshal(resBytes, iStruct)

	return iStruct, &resp.Header, resp.StatusCode, err
}

func RequestRawTest(method networkStruct.Method, targetURL string, params map[string]interface{}, headers map[string]string) (string, *http.Header, int, error) {

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return "", nil, 0, err
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
	if method == networkStruct.Get {
		parsedURL.RawQuery = param.Encode() //如果参数中有中文参数,这个方法会进行URLEncode
		log.Println("【parsedURL.RawQuery】", parsedURL.RawQuery)
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
	//【5-1】增加content-Length
	if method != networkStruct.Get {
		request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	}

	fmt.Println("\n***************************")
	fmt.Println("Request:")
	fmt.Println("【method,apiURL】", method.String(), targetURL)
	fmt.Println("【parsedURL】", parsedURL.String())
	fmt.Println("【body】", body)
	fmt.Println("【header】", request.Header)
	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return "", nil, 0, err
	}
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)

	//【8】返回
	return string(data), &resp.Header, resp.StatusCode, err

}

func RequestJsonTest(method networkStruct.Method, targetURL string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, *http.Header, int, error) {

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil, 0, err
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
	if method == networkStruct.Get {
		parsedURL.RawQuery = param.Encode() //如果参数中有中文参数,这个方法会进行URLEncode
		log.Println("【parsedURL.RawQuery】", parsedURL.RawQuery)
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
	//【5-1】增加content-Length
	if method != networkStruct.Get {
		if method != networkStruct.Get {
			request.Header.Add("Content-Length", strconv.Itoa(len(param)))
		}
	}

	fmt.Println("\n***************************")
	fmt.Println("Request:")
	fmt.Println("【method,apiURL】", method.String(), targetURL)
	fmt.Println("【parsedURL】", parsedURL.String())
	fmt.Println("【body】", body)
	fmt.Println("【header】", request.Header)
	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println("\n\nResponse:")
	fmt.Println("【body-data】", string(data))

	var netReturn map[string]interface{}
	json.Unmarshal(data, &netReturn)

	fmt.Println("【body-json】", netReturn)
	fmt.Println("***************************")
	//【8】返回
	return netReturn, &resp.Header, resp.StatusCode, err

}

func ReturnResult(w http.ResponseWriter, message string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"msg":  message,
		"data": data,
	})
}

/**
 * @func: ReturnResult json格式返回
 * @author Wiidz
 * @date   2019-11-16
 */
func ReturnError(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"msg":  msg,
		"data": nil,
	})
}

// ParamsInvalid json格式返回参数错误
func ParamsInvalid(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"msg":  err.Error(),
		"data": nil,
	})
}

// GetValidatedJson 获取通过validatMng验证过的json数据
func GetValidatedJson(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return validatorMng.GetError(target)
}

// GetJson 获取body体中的json数据
func GetJson(r *http.Request, target interface{}) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(target)
}

// SetRouterFlag 根据第一个斜线内的自符，记录模块（例如 v1,v2）
// router_flag 和 router_key 来共同判断路由是否合法
func SetRouterFlag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURL := r.RequestURI
		reg := regexp.MustCompile(`/([a-z][0-9]*)/`)
		result := reg.FindAllStringSubmatch(requestURL, -1)
		if len(result) > 0 && len(result[0]) > 1 {
			r = r.WithContext(context.WithValue(r.Context(), "router_flag", result[0][1]))
		}
		next.ServeHTTP(w, r)
	})
}

// CheckMixedRouter 检查混合项目的路由准入
func CheckMixedRouter(requestRouterFlag string, requestRouterKey int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routerFlag, _ := r.Context().Value("router_flag").(string)
		routerKeys, _ := r.Context().Value("router_keys").([]int)
		if len(routerKeys) == 0 {
			ReturnError(w, "登陆主体为空")
			return
		}
		if routerFlag == requestRouterFlag {
			if !sliceHelper.Exist(requestRouterKey, routerKeys) {
				ReturnError(w, "越界操作")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// CheckMixedRouterWithHandler 检查混合项目的路由准入(自定义方法)
func CheckMixedRouterWithHandler(requestRouterFlag string, handler func(r *http.Request, routerFlag string, routerKeys []int) error, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		routerFlag, _ := r.Context().Value("router_flag").(string)
		routerKeys, _ := r.Context().Value("router_keys").([]int)
		if routerKeys == nil {
			routerKeys = []int{}
		}
		if routerFlag == requestRouterFlag {
			if err := handler(r, routerFlag, routerKeys); err != nil {
				ReturnError(w, err.Error())
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequestWithStructTest(method networkStruct.Method, contentType networkStruct.ContentType, targetURL string, params map[string]interface{}, headers map[string]string, iStruct interface{}) (interface{}, *http.Header, int, error) {

	log.Println("【method】", method.String())
	log.Println("【targetURL】", targetURL)
	log.Println("【params】", params)
	log.Println("【headers】", headers)
	log.Println("【iStruct】", iStruct)

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil, 0, err
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
	if method == networkStruct.Get {
		parsedURL.RawQuery = param.Encode() // 如果参数中有中文参数,这个方法会进行URLEncode
	} else if contentType == networkStruct.BodyJson {
		jsonByte, _ := json.Marshal(params)
		body = bytes.NewReader(jsonByte)
	} else {
		body = strings.NewReader(param.Encode())
	}

	log.Println("【parsedURL】", parsedURL.String())
	log.Println("【body】", body)

	//【4】创建请求
	request, err := http.NewRequest(method.String(), parsedURL.String(), body)
	if err != nil {
		panic(err)
	}

	//【5】增加header选项
	log.Println("len(headers)", len(headers))
	if len(headers) > 0 {
		for k, v := range headers {
			request.Header.Add(k, v)
		}
	}
	request.Header.Set("Content-Type", contentType.GetContentTypeStr())
	//request.Header.Set("Content-Type", "application/json")

	//【5-1】增加content-Length
	//if method != networkStruct.Get {
	//	request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	//}
	log.Println("【request.Header】", request.Header)
	log.Println("【request.Body】", request.Body)
	//outBodyStr, _ := ioutil.ReadAll(request.Body)
	//log.Println("【outBodyStr】", string(outBodyStr))

	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	//【7】读取body
	resStr, err := ioutil.ReadAll(resp.Body)

	log.Println("【StatusCode】", resp.StatusCode)
	log.Println("【resStr】", string(resStr))

	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json2.Unmarshal(resStr, iStruct)

	//【8】返回
	return iStruct, &resp.Header, resp.StatusCode, err
}

func SendRequest(method networkStruct.Method, contentType networkStruct.ContentType, targetURL string, params map[string]interface{}, headers map[string]string) ([]byte, *http.Header, int, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, nil, 0, err
	}

	var body io.Reader

	if params == nil {
		params = map[string]interface{}{}
	}

	switch {
	case method == networkStruct.Get || method == networkStruct.Delete || contentType == networkStruct.Query:
		query := parsedURL.Query()
		for key, value := range params {
			query.Set(typeHelper.ToString(key), typeHelper.ToString(value))
		}
		parsedURL.RawQuery = query.Encode()
	case contentType == networkStruct.BodyJson:
		jsonBytes, err := json.Marshal(params)
		if err != nil {
			return nil, nil, 0, err
		}
		body = bytes.NewReader(jsonBytes)
	default:
		form := url.Values{}
		for key, value := range params {
			form.Set(typeHelper.ToString(key), typeHelper.ToString(value))
		}
		body = strings.NewReader(form.Encode())
	}

	request, err := http.NewRequest(method.String(), parsedURL.String(), body)
	if err != nil {
		return nil, nil, 0, err
	}

	for k, v := range headers {
		request.Header.Set(k, v)
	}

	if request.Header.Get("Content-Type") == "" {
		if ct := contentType.GetContentTypeStr(); ct != "" {
			request.Header.Set("Content-Type", ct)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	resBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &resp.Header, resp.StatusCode, err
	}

	return resBytes, &resp.Header, resp.StatusCode, nil
}

func RequestWithStruct(method networkStruct.Method, contentType networkStruct.ContentType, targetURL string, params map[string]interface{}, headers map[string]string, iStruct interface{}) (interface{}, *http.Header, int, error) {

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return nil, nil, 0, err
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
	if method == networkStruct.Get {
		parsedURL.RawQuery = param.Encode() // 如果参数中有中文参数,这个方法会进行URLEncode
	} else if contentType == networkStruct.BodyJson {
		jsonByte, _ := json.Marshal(params)
		body = bytes.NewReader(jsonByte)
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
	request.Header.Set("Content-Type", contentType.GetContentTypeStr())
	//request.Header.Set("Content-Type", "application/json")

	//【5-1】增加content-Length
	//if method != networkStruct.Get {
	//	request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	//}
	//outBodyStr, _ := ioutil.ReadAll(request.Body)
	//log.Println("【outBodyStr】", string(outBodyStr))

	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	//【7】读取body
	resStr, err := ioutil.ReadAll(resp.Body)

	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json2.Unmarshal(resStr, iStruct)

	//【8】返回
	return iStruct, &resp.Header, resp.StatusCode, err
}

// MyRequest 自定义请求
func MyRequest(params *networkStruct.MyRequestParams) (resData *networkStruct.MyRequestResp, err error) {

	resData = &networkStruct.MyRequestResp{}

	//【1】解析URL
	var parsedURL *url.URL
	parsedURL, err = url.Parse(params.URL)
	if err != nil {
		return
	}

	//【2】创建client
	client := &http.Client{}

	//【3】构造参数
	param := url.Values{}
	for key, value := range params.Params {
		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}

	var body io.Reader
	if params.Method == networkStruct.Get {
		parsedURL.RawQuery = param.Encode() // 如果参数中有中文参数,这个方法会进行URLEncode
	} else if params.ContentType == networkStruct.BodyJson {
		jsonByte, _ := json.Marshal(params)
		body = bytes.NewReader(jsonByte)
	} else {
		body = strings.NewReader(param.Encode())
	}

	if params.Debug {
		fmt.Printf("[MyRequest] => %s %s\n", params.Method.String(), parsedURL.String())
		fmt.Printf("[MyRequest] headers=%v\n", params.Headers)
		fmt.Printf("[MyRequest] contentType=%s\n", params.ContentType.GetContentTypeStr())
		fmt.Printf("[MyRequest] params=%v\n", params.Params)
	}

	//【4】创建请求
	request, err := http.NewRequest(params.Method.String(), parsedURL.String(), body)
	if err != nil {
		panic(err)
	}

	//【5】增加header选项
	if len(params.Headers) > 0 {
		for k, v := range params.Headers {
			request.Header.Add(k, v)
		}
	}
	request.Header.Set("Content-Type", params.ContentType.GetContentTypeStr())

	//【6】发送请求
	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	resData.StatusCode = resp.StatusCode
	resData.Headers = resp.Header

	//【7】读取body
	resStr, err := ioutil.ReadAll(resp.Body)
	resData.ResStr = string(resStr)

	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json2.Unmarshal(resStr, params.ResStruct)
	if err != nil {
		if params.Debug {
			fmt.Printf("[MyRequest] unmarshal failed: %v\n", err)
		}
		err = nil // 解析失败不报错
	} else {
		resData.ResStruct = params.ResStruct
		resData.IsParsedSuccess = true
	}

	if params.Debug {
		fmt.Printf("[MyRequest] <= status=%d\n", resData.StatusCode)
		fmt.Printf("[MyRequest] body=%s\n", resData.ResStr)
		fmt.Printf("[MyRequest] parsed=%v\n", resData.IsParsedSuccess)
	}

	//【8】返回
	return
}

// GetDomainFromURL 从url中提取文件名
func GetDomainFromURL(url string) (domain string) {

	reg := regexp.MustCompile(`(\w+)(\.\w+)+`)
	result := reg.FindStringSubmatch(url)
	if len(result) > 2 {
		domain = result[0]
	}
	return
}

// GetFileNameFromURL 从url中提取文件名
func GetFileNameFromURL(url string) (wholeName, fileName string, fileType string) {
	reg := regexp.MustCompile(`\w+\.\w+\/([^?]*)\??`)
	result := reg.FindStringSubmatch(url)
	if len(result) == 2 {
		wholeName = result[1]
		temp := typeHelper.ExplodeStr(wholeName, ".")
		if len(temp) == 2 {
			fileName = temp[0]
			fileType = temp[1]
		}
	}
	return
}
