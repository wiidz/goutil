package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/mngs/mysqlMng"
	"github.com/wiidz/goutil/mngs/validatorMng"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type NetworkHelper struct{}

type Method int8

const (
	Get     Method = 1
	Post    Method = 2
	Put     Method = 3
	Delete  Method = 4
	Options Method = 5
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
	case Options:
		return "OPTIONS"
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
func (networkHelper *NetworkHelper) GetRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, *http.Header, error) {

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
func (networkHelper *NetworkHelper) PostRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, *http.Header, error) {

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
func (*NetworkHelper) DownloadFileWithFormat(targetURL, localPath, format string, headers map[string]string) (fileName, tempPath string, header *http.Header, err error) {

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
	defer resp.Body.Close()

	header = &resp.Header

	if err != nil {
		return
	}

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

func (*NetworkHelper) RequestRaw(method Method, targetURL string, params map[string]interface{}, headers map[string]string) (string, *http.Header, int, error) {

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
	//【5-1】增加content-Length
		if method != 1{
		request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	}

	//【6】发送请求
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)

	//【8】返回
	return string(data), &resp.Header, resp.StatusCode, err

}

func (*NetworkHelper) RequestJson(method Method, targetURL string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, *http.Header, int, error) {

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

	//【5-1】增加content-Length
		if method != 1{
		request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	}

	//【6】发送请求
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)

	var netReturn map[string]interface{}
	json.Unmarshal(data, &netReturn)

	//【8】返回
	return netReturn, &resp.Header, resp.StatusCode, err

}

func (*NetworkHelper) RequestRawTest(method Method, targetURL string, params map[string]interface{}, headers map[string]string) (string, *http.Header, int, error) {

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
	if method == Get {
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
		if method != 1{
		request.Header.Add("Content-Length", strconv.Itoa(len(param)))
	}

	fmt.Println("\n***************************")
	fmt.Println("Request:")
	fmt.Println("【method,apiURL】", method.String(), targetURL)
	fmt.Println("【parsedURL】", parsedURL.String())
	fmt.Println("【body】", body)
	fmt.Println("【header】", request.Header)
	//【6】发送请求
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)

	//【8】返回
	return string(data), &resp.Header, resp.StatusCode, err

}

func (*NetworkHelper) RequestJsonTest(method Method, targetURL string, params map[string]interface{}, headers map[string]string) (map[string]interface{}, *http.Header, int, error) {

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
	if method == Get {
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
	if method != 1{
			if method != 1{
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
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	//【7】读取body
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Println("\n\nResponse:")
	fmt.Println("【body-data】", string(data))

	var netReturn map[string]interface{}
	json.Unmarshal(data, &netReturn)

	fmt.Println("【body-json】", netReturn)
	fmt.Println("***************************\n")
	//【8】返回
	return netReturn, &resp.Header, resp.StatusCode, err

}

func ReturnResult(ctx iris.Context, message string, data interface{}, statusCode int) {

	ctx.StatusCode(statusCode)

	ctx.JSON(iris.Map{
		"msg":  message,
		"data": data,
	})
	return
}

/**
 * @func: ReturnResult json格式返回
 * @author Wiidz
 * @date   2019-11-16
 */
func ReturnError(ctx iris.Context, msg string) {

	ctx.StatusCode(404)

	ctx.JSON(iris.Map{
		"msg":  msg,
		"data": nil,
	})
	return
}

/**
 * @func: ReturnResult json格式返回
 * @author Wiidz
 * @date   2019-11-16
 */
func ParamsInvalid(ctx iris.Context, err error) {

	ctx.StatusCode(404)

	ctx.JSON(iris.Map{
		"msg":  "参数无效",
		"data": err.Error(),
	})
	return
}

// GetValidatedJson 获取通过validatMng验证过的json数据
func GetValidatedJson(ctx iris.Context, target interface{}) error {
	err := ctx.ReadJSON(target)
	if err != nil {
		return err
	}
	err = validatorMng.GetError(target)
	return err
}

// GetJson 获取body体中的json数据
func GetJson(ctx iris.Context, target interface{}) error {
	err := ctx.ReadJSON(target)
	return err
}


// QueryParamsFilter get参数过滤器+验证
func QueryParamsFilter(ctx iris.Context, params interface{}) (condition,etc map[string]interface{},err error) {

	//【1】提取字段
	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)

	condition = map[string]interface{}{}
	etc = map[string]interface{}{}

	//【3】遍历过滤参数
	for i := 0; i < t.Elem().NumField(); i++ {

		//【3-1】获取标签值
		field := t.Elem().Field(i)
		belong := field.Tag.Get("belong")
		json := field.Tag.Get("json")
		kind := field.Tag.Get("kind")
		fieldName := field.Tag.Get("field")
		defaultValue := field.Tag.Get("default")

		if fieldName == "" {
			fieldName = json
		}

		//【3-2】取值
		temp := ctx.URLParam(json)
		if temp == "" {
			if defaultValue == "" {
				// 即没有默认值也没有值传递过来的，跳过
				continue
			}
			temp = defaultValue
		}

		//【3-3】将值处理成需要的格式
		formattedValue := getFormattedValue(field.Type.String(), temp)

		//【3-4】将值赋值给param结构体的对应字段
		val := reflect.ValueOf(formattedValue)
		v.Elem().Field(i).Set(val)

		//【3-5】填充到condition、etcMap
		switch belong {
		case "condition":
			switch kind {
			case "like":
				condition[fieldName] = []interface{}{"like", "%" + temp + "%"}
			case "between":
				tempSlice := typeHelper.Explode(temp, ",")
				condition[fieldName] = []interface{}{"between", tempSlice[0], tempSlice[1]}
			default:
				condition[fieldName] = temp
			}
		case "etc":
			etc[json] = temp
		}
	}

	//【4】验证参数是否合法
	err = validatorMng.GetError(params)
	if err != nil {
		return
	}

	return
}

// JsonParamsFilter 依据json格式从body体中过滤参数+验证
func JsonParamsFilter(params mysqlMng.JsonInterface) (condition,value,etc map[string]interface{},err error){

	//【1】提取字段
	t := reflect.TypeOf(params)
	v := reflect.ValueOf(params)
	rawJsonMap := params.GetRawJsonMap()

	// 【2】初始化变量
	condition = map[string]interface{}{}
	value = map[string]interface{}{}
	etc = map[string]interface{}{}

	//【3】遍历过滤参数
	for i := 0; i < t.Elem().NumField(); i++ {

		//【3-1】获取标签值
		field := t.Elem().Field(i)
		belong := field.Tag.Get("belong")
		jsonTag := field.Tag.Get("json")
		kind := field.Tag.Get("kind")
		fieldName := field.Tag.Get("field")
		defaultValue := field.Tag.Get("default")

		if fieldName == "" {
			fieldName = jsonTag
		}

		//【2-2】将值处理成需要的格式
		temp, ok := rawJsonMap[jsonTag]
		if !ok {
			if defaultValue == "" {
				//即没有默认值也没有值传递过来的，跳过
				continue
			}
			temp = defaultValue
		}
		formattedValue := getFormattedValue(field.Type.String(), temp)

		//【2-3】将值赋值给param结构体的对应字段
		val := reflect.ValueOf(formattedValue)
		v.Elem().Field(i).Set(val)

		//【3-3】填充
		switch belong {
		case "condition":
			switch kind {
			case "like":
				condition[fieldName] = []interface{}{"like", "%" + formattedValue.(string) + "%"}
			case "between":
				tempSlice := typeHelper.Explode(formattedValue.(string), ",")
				condition[fieldName] = []interface{}{"between", tempSlice[0], tempSlice[1]}
			case "in":
				fieldName := field.Tag.Get("field")
				condition[fieldName] = []interface{}{"in"}
				condition[fieldName] = append(condition[fieldName].([]interface{}), formattedValue)
			default:
				condition[fieldName] = formattedValue
			}
		case "value":
			value[fieldName] = formattedValue
		case "etc":
			etc[fieldName] = formattedValue
		}
	}

	return
}

// getFormattedValue 获取指定格式的数值
func getFormattedValue(t string, value interface{}) interface{} {
	switch t {
	case "int":
		if typeHelper.GetType(value) == "float64" {
			return typeHelper.Float64ToInt(value.(float64))
		} else {
			return typeHelper.Str2Int(value.(string))
		}
	case "int8":
		if typeHelper.GetType(value) == "string" {
			return typeHelper.Str2Int8(value.(string))
		} else {
			return typeHelper.Float64ToInt8(value.(float64))
		}
	case "float64":
		if str, ok := value.(string); ok {
			return typeHelper.Str2Float64(str)
		} else {
			return value.(float64)
		}
	case "[]int":
		if str, ok := value.(string); ok {
			slice := typeHelper.ExplodeInt(str, ",")
			return slice
		} else {
			slice := typeHelper.Float64ToIntSlice(value.([]interface{}))
			return slice
		}
	case "[]string":
		slice := typeHelper.ExplodeStr(value.(string), ",")
		return slice
	default:
		log.Println("t", t)
		log.Println("value", value)
		temp, _ := typeHelper.JsonEncode(value)
		typeHelper.JsonDecodeWithStruct(temp, value)
		return value
	}
}
