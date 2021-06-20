package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type NetworkHelper struct{}


/**
 * @func: GetRequest 发送get请求
 * @author Wiidz
 * @date   2019-11-16
 */
func (networkHelper *NetworkHelper) GetParamsUrl(apiURL string, params map[string]interface{}) (string, error) {

	param := url.Values{}
	var typeHelper TypeHelper
	for key, value := range params {

		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}

	var Url *url.URL
	Url, err := url.Parse(apiURL)
	if err != nil {
		fmt.Printf("解析url错误:\r\n%v", err)
		return "", err
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = param.Encode()
	return Url.String(),err
}

/**
 * @func: GetRequest 发送get请求
 * @author Wiidz
 * @date   2019-11-16
 */
func (networkHelper *NetworkHelper) GetRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, error) {

	param := url.Values{}
	var typeHelper TypeHelper
	for key, value := range params {

		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}

	var Url *url.URL
	Url, er := url.Parse(apiURL)
	if er != nil {
		fmt.Printf("解析url错误:\r\n%v", er)
		return nil, er
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = param.Encode()

	log.Println("networkHelper.GetRequest:",Url)

	resp, err := http.Get(Url.String())
	if err != nil {
		fmt.Println("err:", err)
		return nil, err
	}
	defer resp.Body.Close()

	data, e := ioutil.ReadAll(resp.Body)

	var netReturn map[string]interface{}

	json.Unmarshal(data, &netReturn)

	return netReturn, e
}

/**
 * @func: PostRequest 发送post请求
 * @author Wiidz
 * @date   2019-11-16
 */
func (networkHelper *NetworkHelper) PostRequest(apiURL string, params map[string]interface{}) (map[string]interface{}, error) {

	param := url.Values{}
	var typeHelper TypeHelper
	for key, value := range params {

		k := typeHelper.ToString(key)
		v := typeHelper.ToString(value)
		param.Set(k, v)
	}

	log.Println("networkHelper.PostRequest:",apiURL)
	log.Println("params:",param)

	resp, err := http.PostForm(apiURL, param)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, e := ioutil.ReadAll(resp.Body)

	var netReturn map[string]interface{}

	fmt.Println(data)

	json.Unmarshal(data, &netReturn)

	return netReturn, e
}

//
///**
// * @func: ReturnResult json格式返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func ReturnResult(ctx iris.Context, message string, data interface{}, statusCode int) {
//
//	ctx.StatusCode(statusCode)
//
//	ctx.JSON(iris.Map{
//		"msg":  message,
//		"data": data,
//	})
//	return
//}
//
///**
// * @func: ReturnResult json格式返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func ParamsError(ctx iris.Context) {
//
//	ctx.StatusCode(422)
//
//	ctx.JSON(iris.Map{
//		"msg":  "参数错误",
//		"data": nil,
//	})
//	return
//}
//
///**
// * @func: ReturnResult json格式返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func Forbidden(ctx iris.Context) {
//
//	ctx.StatusCode(403)
//
//	ctx.JSON(iris.Map{
//		"msg":  "无权访问",
//		"data": nil,
//	})
//	return
//}
//
///**
// * @func: GetJsonArrayParams 从body体中提取参数，以 []map[string]interface{} 返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func GetJsonArrayParams(ctx iris.Context) []map[string]interface{} {
//
//	post_data := make([]map[string]interface{}, 0)
//	data := ctx.Request().Body
//
//	js, _ := ioutil.ReadAll(data)
//	json.Unmarshal(js, &post_data)
//
//	return post_data
//}
//
///**
// * @func: GetJsonArrayParams 从body体中提取参数，以 map[string]interface{} 返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func GetJsonParams(ctx iris.Context) map[string]interface{} {
//
//	post_data := make(map[string]interface{}, 0)
//
//	data := ctx.Request().Body
//
//	js, _ := ioutil.ReadAll(data)
//
//	json.Unmarshal(js, &post_data)
//
//	return post_data
//}
//
///**
// * @func: GetJsonArrayParams 从body体中提取参数，以 map[string]interface{} 返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func GetJsonParamsWithStruct(ctx iris.Context, istruct interface{}) interface{} {
//
//	//post_data := make(map[string]interface{}, 0)
//
//	data := ctx.Request().Body
//
//	js, _ := ioutil.ReadAll(data)
//
//	json.Unmarshal(js, &istruct)
//
//	return istruct
//}

///**
// * @func: GetFilteredParams 从query中获取指定字段的值，以 map[string]interface{} 返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func GetFilteredParams(ctx iris.Context, fields []string) map[string]interface{} {
//	temp := ""
//	container := make(map[string]interface{})
//	for _, v := range fields {
//		temp = ctx.URLParam(v)
//		if temp != "" {
//			container[v] = temp
//		}
//	}
//	return container
//}

func (networkHelper *NetworkHelper) DownloadFile(url string, fb func(string) error) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	var typeHelper TypeHelper
	var strHelper StrHelper
	tempSlice := typeHelper.Explode(url, ".")
	format := tempSlice[len(tempSlice)-1]
	tempPath := "/tmp/download/" + strHelper.GetRandomString(10) + "." + format.(string)
	out, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return fb(tempPath)
}


func (*NetworkHelper) DownloadFileWithFormat(url,format string, fb func(string) error) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	var strHelper StrHelper
	tempPath := "/tmp/download/" + strHelper.GetRandomString(10) + "." + format
	out, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return fb(tempPath)
}

//
///**
// * @func: ReturnResult json格式返回
// * @author Wiidz
// * @date   2019-11-16
// */
//func (*NetworkHelper) ReturnResult(ctx iris.Context, message string, data interface{}, statusCode int) {
//
//	ctx.StatusCode(statusCode)
//
//	ctx.JSON(iris.Map{
//		"msg":  message,
//		"data": data,
//	})
//	return
//}
//
