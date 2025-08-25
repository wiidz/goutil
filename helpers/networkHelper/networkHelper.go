package networkHelper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/kataras/iris/v12"
	"github.com/wiidz/goutil/helpers/sliceHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/mngs/validatorMng"
	"github.com/wiidz/goutil/structs/networkStruct"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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
	resp, _ := client.Do(request)
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
	resp, _ := client.Do(request)
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
	resp, _ := client.Do(request)
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
	resp, _ := client.Do(request)
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

	ctx.StatusCode(400)

	ctx.JSON(iris.Map{
		"msg":  msg,
		"data": nil,
	})
	return
}

// ParamsInvalid json格式返回参数错误
func ParamsInvalid(ctx iris.Context, err error) {

	ctx.StatusCode(400)

	_ = ctx.JSON(iris.Map{
		"msg":  err.Error(),
		"data": nil,
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

// BuildParams 构建参数
func BuildParams(ctx iris.Context, params networkStruct.ParamsInterface, contentType networkStruct.ContentType) (err error) {

	//【1】把值填充进结构体
	err = fillParams(ctx, params, contentType)
	if err != nil {
		return
	}

	//【3】验证参数
	err = validatorMng.GetError(params)
	if err != nil {
		return
	}

	//【4】处理和填充
	err = handleParams(params)
	if err != nil {
		return
	}

	//【5】返回
	return
}

// getValueFromCtx 获取参数
func fillParams(ctx iris.Context, params networkStruct.ParamsInterface, contentType networkStruct.ContentType) (err error) {
	switch contentType {
	case networkStruct.Query:

		//【1】写入结构体
		err = ctx.ReadQuery(params)
		if err != nil {
			return
		}

		//【2】获取RawMap
		temp := ctx.URLParams()
		params.SetRawMap(typeHelper.StrMapToInterface(temp))
		break

	case networkStruct.BodyJson:

		//【1】写入结构体

		body := ctx.Request().Body
		buf, _ := ioutil.ReadAll(body)
		jsonStr := string(buf)
		//log.Println("【jsonStr】",jsonStr)
		jsonMap := typeHelper.JsonDecodeMap(jsonStr) // 映射到map
		params.SetRawMap(jsonMap)

		err = typeHelper.JsonDecodeWithStruct(jsonStr, params) // 映射到结构体
		//log.Println("【params-beforeHandle】",params)
		if err != nil {
			log.Println("err", err, params)
			return
		}

		//【2】获取RawMap
		//body := ctx.Request().Body
		//buf, _ := ioutil.ReadAll(body)

		////【2-1】写入到map[string]interface{},主要是看下前端发了哪些字段过来，保证0值不会被刷掉
		//tempMap := typeHelper.JsonDecodeMap(string(buf))
		//sendFields := mapHelper.GetKeys(tempMap)
		//params.SetParamFields(sendFields)

		//【2-2】写入到结构体

		break
	case networkStruct.BodyForm:
		formMap := make(map[string]interface{})

		// PostForm 包含 body form 的字段（不包含 URL 查询参数）
		for k, v := range ctx.Request().PostForm {
			if len(v) == 1 {
				formMap[k] = v[0]
			} else {
				formMap[k] = v
			}
		}

		params.SetRawMap(formMap)

		break
	case networkStruct.XWWWForm:
		err = ctx.ReadForm(params)
		if err != nil {
			return
		}

		// params.SetRawMap(typeHelper.StrMapToInterface(temp))
		break

	default:
		err = errors.New("未能匹配数据类型")
	}
	return
}

// handleParams 处理和填充数据
func handleParams(params networkStruct.ParamsInterface) (err error) {

	//【1】获取结构体反射结构
	structType := reflect.TypeOf(params)
	structValues := reflect.ValueOf(params)
	rawMap := params.GetRawMap()

	//【2】初始化变量
	//RawMap := map[string]interface{}{}

	condition := map[string]interface{}{}
	value := map[string]interface{}{}
	etc := map[string]interface{}{}

	for i := 0; i < structType.Elem().NumField(); i++ {

		//【1】获取标签
		field := structType.Elem().Field(i)
		fieldType := field.Type // 字段的类型

		jsonTag := field.Tag.Get("json")    // body体中使用
		urlTag := field.Tag.Get("url")      // query中使用
		fieldName := field.Tag.Get("field") // 字段写入condition、value的名称

		belong := field.Tag.Get("belong") // 值的归属，例如value、condition、etc
		kind := field.Tag.Get("kind")     // 值类型，例如between、like

		defaultValue := field.Tag.Get("default")

		if jsonTag == "" && urlTag == "" {
			// 没有jsonTag，视为内部填充，则跳过
			continue
		}

		//【2】确定填充的键名
		if fieldName == "" {
			if jsonTag != "" {
				fieldName = jsonTag
			} else if urlTag != "" {
				fieldName = urlTag
			}
		}

		//【3】填充默认值
		currentValue := reflect.Indirect(structValues).FieldByName(field.Name) // 当前结构体设置的值 reflect.Value 类型
		// 注意这里判断不要用 reflect.Value == reflect.Value，会一直false
		if fieldType.Kind() == reflect.Struct {
			//log.Println("Struct")
			//log.Println(reflect.Zero(fieldType).Interface())
		} else if fieldType.Kind() == reflect.Slice {
			//log.Println("Slice")
			//log.Println(reflect.Zero(fieldType).Interface())
			//log.Println("slice")
			//if len(currentValue.Interface()) == 0 {
			//	log.Println("Slice zero")
			//}
			//continue // 结构体类型和切片类型 默认不予填充和计算
		} else if fieldType.Kind() == reflect.Map {
			//log.Println("Slice")
			//log.Println(reflect.Zero(fieldType).Interface())
			//log.Println("slice")
			//if len(currentValue.Interface()) == 0 {
			//	log.Println("Slice zero")
			//}
			//continue // 结构体类型和切片类型 默认不予填充和计算
		} else if currentValue.Interface() == reflect.Zero(fieldType).Interface() {
			if defaultValue != "" {
				var formattedDefaultValue interface{}
				formattedDefaultValue, err = getFormattedValue(field.Type.String(), defaultValue)
				if err != nil {
					return
				}
				currentValue = reflect.ValueOf(formattedDefaultValue)
				structValues.Elem().Field(i).Set(currentValue) // 记得填充回结构体
			} else if _, ok := rawMap[fieldName]; !ok {
				// 判断rawMap里面有没有值
				continue
			}
			// 给的就是零值，继续
		}

		var formattedValue interface{}
		formattedValue, err = getFormattedValue(field.Type.String(), currentValue.Interface()) // 格式化后的当前值
		if err != nil {
			log.Println("field.Type.String()", field.Type.String(), err)
			return
		}

		//【4】按照belong进行填充
		switch belong {
		case "condition":
			switch kind {
			case "like":
				condition[fieldName] = []interface{}{"like", "%" + formattedValue.(string) + "%"}
			case "between":
				tempSlice := typeHelper.Explode(formattedValue.(string), ",")
				if len(tempSlice) == 2 {
					condition[fieldName] = []interface{}{"between", tempSlice[0], tempSlice[1]}
				}
			case "in":
				condition[fieldName] = []interface{}{"in"}
				condition[fieldName] = append(condition[fieldName].([]interface{}), formattedValue)
			case "not in":
				condition[fieldName] = []interface{}{"not in"}
				condition[fieldName] = append(condition[fieldName].([]interface{}), formattedValue)
			case "!=":
				condition[fieldName] = []interface{}{"!="}
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

	//【5】其他填充
	params.SetCondition(condition)
	params.SetValue(value)
	params.SetEtc(etc)

	// Tips：特别注意这里必须要set，外部的pageNow等字段，嵌入本结构体，外部的pageSize有值，但是取的GetPageSize还是内部的，注意
	if pageNow, ok := etc["page_now"].(int); ok {
		params.SetPageNow(pageNow)
	} else {
		params.SetPageNow(0)
	}
	if pageSize, ok := etc["page_size"].(int); ok {
		params.SetPageSize(pageSize)
	} else {
		params.SetPageSize(10)
	}
	if order, ok := etc["order"].(string); ok {
		params.SetOrder(order)
	} else {
		params.SetOrder("id asc")
	}
	return
}

// getFormattedValue 获取指定格式的数值
func getFormattedValue(t string, value interface{}) (data interface{}, err error) {
	switch t {
	case "string":
		data = typeHelper.ForceString(value)
	case "int":
		data = typeHelper.ForceInt(value)
	case "int8":
		data = typeHelper.ForceInt8(value)
	//case "int64":
	//	return typeHelper.ForceInt64(value)
	case "uint64":
		data = typeHelper.ForceUint64(value)
	case "float64":
		data = typeHelper.ForceFloat64(value)
	case "[]int":
		data = typeHelper.ForceIntSlice(value)
	case "[]int8":
		data = typeHelper.ForceInt8Slice(value)
	case "[]uint64":
		data = typeHelper.ForceUint64Slice(value)
	case "[]float64":
		data = typeHelper.ForceFloat64Slice(value)
	case "[]string":
		data = typeHelper.ForceStrSlice(value)
	default:
		// 其他情况默认是结构体，还有 结构体 slice的情况
		//temp, _ := typeHelper.JsonEncode(value)
		//err = typeHelper.JsonDecodeWithStruct(t, value)
		data = value
	}
	return
}

// SetRouterFlag 根据第一个斜线内的自符，记录模块（例如 v1,v2）
// router_flag 和 router_key 来共同判断路由是否合法
func SetRouterFlag(app *iris.Application) {
	app.Use(func(ctx iris.Context) {
		requestURL := ctx.Request().RequestURI
		reg := regexp.MustCompile(`/([a-z][0-9]*)/`)
		result := reg.FindAllStringSubmatch(requestURL, -1)
		//log.Println("result", result[0][1])
		ctx.Values().Set("router_flag", result[0][1])
		ctx.Next()
	})
}

// CheckMixedRouter 检查混合项目的路由准入
func CheckMixedRouter(app *iris.Application, requestRouterFlag string, requestRouterKey int) {
	app.Use(func(ctx iris.Context) {
		routerFlag := ctx.Values().Get("router_flag")
		routerKeys := ctx.Values().Get("router_keys") // 注意 这里已经改成了slice
		routerKeySlice, flag := routerKeys.([]int)
		if !flag {
			ReturnError(ctx, "登陆体结构有误")
			return
		}
		if len(routerKeySlice) == 0 {
			ReturnError(ctx, "登陆主体为空")
			return
		}

		if routerFlag.(string) == requestRouterFlag {
			// 查找slice里面有没有requestRouterKey
			if !sliceHelper.Exist(requestRouterKey, routerKeySlice) {
				ReturnError(ctx, "越界操作")
				return
			}
		}
		ctx.Next()
	})
}

// CheckMixedRouterWithHandler 检查混合项目的路由准入(自定义方法)
func CheckMixedRouterWithHandler(app *iris.Application, requestRouterFlag string, handler func(ctx iris.Context, routerFlag string, routerKeys []int)) {
	app.Use(func(ctx iris.Context) {
		routerFlag := ctx.Values().Get("router_flag")
		routerKeys := ctx.Values().Get("router_keys") // 注意 这里已经改成了slice
		routerKeySlice, flag := routerKeys.([]int)
		if !flag {
			// 这里不返回错误，因为可能不存在
			//ReturnError(ctx, "登陆体结构有误")
			//return
			routerKeySlice = []int{}
		}
		//if len(routerKeySlice) == 0 {
		//ReturnError(ctx, "登陆主体为空")
		//return
		//}
		handler(ctx, routerFlag.(string), routerKeySlice)
		//ctx.Next()
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
	resp, _ := client.Do(request)
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
	resp, _ := client.Do(request)
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
	defer resp.Body.Close()
	if err != nil {
		return
	}

	resData.StatusCode = resp.StatusCode
	resData.Headers = resp.Header

	//【7】读取body
	resStr, err := ioutil.ReadAll(resp.Body)
	resData.ResStr = string(resStr)

	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json2.Unmarshal(resStr, params.ResStruct)
	if err != nil {
		err = nil // 解析失败不报错
	} else {
		resData.ResStruct = params.ResStruct
		resData.IsParsedSuccess = true
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
