package typeHelper

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Kind string

const (
	Invalid       Kind = "invalid"
	Bool          Kind = "bool"
	Int           Kind = "int"
	Int8          Kind = "int8"
	Int16         Kind = "int16"
	Int32         Kind = "int32"
	Int64         Kind = "int64"
	Uint          Kind = "uint"
	Uint8         Kind = "uint8"
	Uint16        Kind = "uint16"
	Uint32        Kind = "uint32"
	Uint64        Kind = "uint64"
	Uintptr       Kind = "uintptr"
	Float32       Kind = "float32"
	Float64       Kind = "float64"
	Complex64     Kind = "complex64"
	Complex128    Kind = "complex128"
	Array         Kind = ""
	Chan          Kind = ""
	Func          Kind = ""
	Interface     Kind = "interface{}"
	Map           Kind = ""
	Ptr           Kind = ""
	Slice         Kind = ""
	String        Kind = "string"
	Struct        Kind = ""
	UnsafePointer Kind = ""
)

/**
 * @func: Implode 将slice转换成字符串
 * @author Wiidz
 * @date   2019-11-16
 */
func Implode(data interface{}, glue string) string {
	var tmp []string
	for _, item := range data.([]interface{}) {
		tmp = append(tmp, ToString(item))
	}

	return strings.Join(tmp, glue)
}

// ImplodeInt8 将Int8切片转换成string
func ImplodeInt8(data []int8, glue string) string {
	var tmp []string
	for _, item := range data {
		tmp = append(tmp, ToString(item))
	}

	return strings.Join(tmp, glue)
}

// ImplodeUint64 将slice转换成字符串
func ImplodeUint64(data []uint64, glue string) string {
	var tmp []string
	for _, item := range data {
		tmp = append(tmp, strconv.FormatUint(item, 10))
	}

	return strings.Join(tmp, glue)
}

/**
 * @func: ImplodeInt 将int slice转换成字符串
 * @author Wiidz
 * @date   2019-11-16
 */
func ImplodeInt(data []int, glue string) string {
	var tmp []string
	for _, item := range data {
		tmp = append(tmp, strconv.Itoa(item))
	}

	return strings.Join(tmp, glue)
}

// ImplodeStr 将 str slice转换成字符串
func ImplodeStr(data []string, glue string) string {
	var tmp []string
	for _, item := range data {
		tmp = append(tmp, item)
	}

	return strings.Join(tmp, glue)
}

/**
 * @func: Explode  字符串转slice, 接受混合类型
 * @author Wiidz
 * @date   2019-11-16
 */
func Explode(data string, sep string) []interface{} {

	if len(data) == 0 {
		return []interface{}{}
	}

	old := strings.Split(data, sep)

	newS := make([]interface{}, len(old))

	for i, v := range old {
		newS[i] = v
	}

	return newS

}

/**
 * @func: ExplodeStr  字符串转str clise
 * @author Wiidz
 * @date   2019-11-16
 */
func ExplodeStr(data string, sep string) []string {
	if len(data) == 0 {
		return []string{}
	}
	old := strings.Split(data, sep)
	newS := make([]string, len(old))
	for i, v := range old {
		newS[i] = v
	}
	return newS
}

/**
 * @func: Explode  字符串转int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func ExplodeInt(data string, sep string) []int {
	if len(data) == 0 {
		return []int{}
	}
	old := strings.Split(data, sep)
	newS := make([]int, len(old))
	for i, v := range old {
		newS[i], _ = strconv.Atoi(v)
	}
	return newS
}

/**
 * @func: ExplodeUint64  字符串转int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func ExplodeUint64(data string, sep string) []uint64 {
	if len(data) == 0 {
		return []uint64{}
	}
	old := strings.Split(data, sep)
	newS := make([]uint64, len(old))
	for i, v := range old {
		newS[i], _ = strconv.ParseUint(v, 10, 64)
	}
	return newS
}

/**
 * @func: Explode  字符串转int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func ExplodeFloat64(data string, sep string) []float64 {
	if len(data) == 0 {
		return []float64{}
	}
	old := strings.Split(data, sep)
	newS := make([]float64, len(old))
	for i, v := range old {
		newS[i], _ = strconv.ParseFloat(v, 64)
	}
	return newS
}

/**
 * @func: Explode  字符串转int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func ExplodeInt64(data string, sep string) []int64 {
	if len(data) == 0 {
		return []int64{}
	}
	old := strings.Split(data, sep)
	newS := make([]int64, len(old))
	for i, v := range old {
		newS[i], _ = strconv.ParseInt(v, 64, 10)
	}
	return newS
}

// ExplodeInt8 字符串转int8 slice
func ExplodeInt8(data string, sep string) []int8 {
	if len(data) == 0 {
		return []int8{}
	}
	old := strings.Split(data, sep)
	newS := make([]int8, len(old))
	for i, v := range old {
		temp, _ := strconv.Atoi(v)
		newS[i] = int8(temp)
	}
	return newS
}

/**
 * @func: GetType 获取目标的数据类型
 * @author Wiidz
 * @date   2019-11-16
 */
func GetType(params interface{}) string {
	vT := reflect.TypeOf(params)
	log.Println("vT", vT)
	if params == nil {
		return "nil"
	}
	return vT.String()

	////数据初始化
	//v := reflect.ValueOf(params)
	////log.Println("")
	//
	////获取传递参数类型
	//vT := v.Type()
	////类型名称对比
	//return vT.String()
}

/**
 * @func: Empty 判断目标是否为空值或0
 * @author Wiidz
 * @date   2019-11-16
 */
func Empty(arg interface{}) bool {

	if arg == nil {
		return true
	}

	var dataType = reflect.TypeOf(arg)

	return reflect.Zero(dataType) == arg

	//switch arg.(type) {
	//
	//case int:
	//	return If(arg.(int) == int(0), true, false).(bool)
	//case int64:
	//	return If(arg.(int64) == int64(0), true, false).(bool)
	//
	//case float64:
	//	return If(arg.(float64) == float64(0.00), true, false).(bool)
	//
	//case []byte:
	//
	//	return If(len(arg.([]byte)) == 0, true, false).(bool)
	//
	//case string:
	//	return If(arg.(string) == " " || arg.(string) == "" || arg.(string) == "0" || arg.(string) == "NULL", true, false).(bool)
	//
	//case map[string]interface{}:
	//
	//	return If(len(arg.(map[string]interface{})) == 0, true, false).(bool)
	//
	//case []interface{}:
	//
	//	return If(len(arg.([]interface{})) == 0, true, false).(bool)
	//
	//case []string:
	//
	//	return If(len(arg.([]string)) == 0, true, false).(bool)
	//
	//case []int64:
	//
	//	return If(len(arg.([]int64)) == 0, true, false).(bool)
	//
	//case []float64:
	//
	//	return If(len(arg.([]float64)) == 0, true, false).(bool)
	//
	//case []int:
	//
	//	return If(len(arg.([]int)) == 0, true, false).(bool)
	//
	//case []map[string]interface{}:
	//
	//	return If(len(arg.([]map[string]interface{})) == 0, true, false).(bool)
	//
	//case types.Nil:
	//
	//	return If(arg == nil, true, false).(bool)
	//
	//case bool:
	//
	//	return !arg.(bool)
	//
	//default:
	//
	//	return true
	//}
}

/**
 * @func: IsType 判断目标是否为指定类型的值
 * @author Wiidz
 * @date   2019-11-16
 */
func IsType(needle interface{}, type_name string) bool {
	if reflect.TypeOf(needle).String() == type_name {
		return true
	}
	return false
}

/**
 * @func: ToString 将任何参数转换为字符串
 * @author Wiidz
 * @date   2019-11-16
 */
func ToString(data interface{}) string {
	switch data.(type) {
	case int:
		return strconv.Itoa(data.(int))
	case int8:
		return strconv.Itoa(int(data.(int8)))
	case int64:
		return strconv.FormatInt(data.(int64), 10)
	case int32:
		return strconv.FormatInt(int64(data.(int32)), 10)
	case uint32:
		return strconv.FormatUint(uint64(data.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(data.(uint64), 10)
	case float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case string:
		return data.(string)
	case time.Time:
		return data.(time.Time).Format("2006-01-02 15:04:05")
	default:
		return ""
	}
}

/**
 * @func: Str2Int64 将字符串转为int64
 * @author Wiidz
 * @date   2019-11-16
 */
func Str2Int64(str string) int64 {
	number, _ := strconv.ParseInt(str, 10, 64)
	return number
}

/**
 * @func: Str2Int 将字符串转为int
 * @author Wiidz
 * @date   2019-11-16
 */
func Str2Int(str string) int {
	ints, _ := strconv.Atoi(str)
	return ints
}

/**
 * @func: Str2Int8 将字符串转为int8
 * @author Wiidz
 * @date   2019-11-16
 */
func Str2Int8(str string) int8 {
	ints, _ := strconv.Atoi(str)
	return int8(ints)
}

/**
 * @func: Str2Int8 将字符串转为int8
 * @author Wiidz
 * @date   2019-11-16
 */
func Str2Uint64(str string) uint64 {
	number, _ := strconv.ParseUint(str, 10, 64)
	return number
}

// Str2Uint 将字符串转为uint
// @author Wiidz
// @date   2023-12-19
func Str2Uint(str string) (res uint, err error) {

	num64, _ := strconv.ParseUint(str, 10, 0)

	// 检查是否在 uint 范围内
	if num64 > uint64(^uint(0)) {
		err = errors.New("超出 uint 范围")
		return
	}

	// 将 uint64 转换为 uint
	res = uint(num64)

	return
}

/**
 * @func: Str2Int 将字符串转为int
 * @author Wiidz
 * @date   2019-11-16
 */
func Float64ToStr(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}

/**
 * @func: Int64ToStr int64转为str
 * @author Wiidz
 * @date   2019-11-16
 */
func Int64ToStr(number int64) string {
	str := strconv.FormatInt(number, 10)
	return str
}

/**
 * @func: Str2Float64 str转为float64
 * @author Wiidz
 * @date   2019-11-16
 */
func Str2Float64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

/**
 * @func: JsonEncode 编码json
 * @author Wiidz
 * @date   2019-11-16
 */
func JsonEncode(data interface{}) (string, error) {
	res, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(res), err
}

func JsonDecode(jsonStr string) (parsedData interface{}) {
	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	_ = json2.Unmarshal([]byte(jsonStr), &parsedData)
	return
}

func JsonDecodeMap(jsonStr string) (parsedData map[string]interface{}) {
	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	_ = json2.Unmarshal([]byte(jsonStr), &parsedData)
	return
}

func JsonDecodeMapStr(jsonStr string) (parsedData map[string]string) {
	var json2 = jsoniter.ConfigCompatibleWithStandardLibrary
	_ = json2.Unmarshal([]byte(jsonStr), &parsedData)
	return
}

// JsonDecodeInt64Slice json解码至int64切片
func JsonDecodeInt64Slice(jsonStr string) []int64 {

	var data []int64

	_ = json.Unmarshal([]byte(jsonStr), &data)

	return data
}

// JsonDecodeUint64Slice json解码至uint64切片
func JsonDecodeUint64Slice(jsonStr string) []uint64 {

	var data []uint64

	_ = json.Unmarshal([]byte(jsonStr), &data)

	return data
}

// JsonDecodeStrSlice json解码至str切片
func JsonDecodeStrSlice(jsonStr string) []string {

	var data []string

	_ = json.Unmarshal([]byte(jsonStr), &data)

	return data
}

// JsonDecodeIntSlice json解码至int切片
func JsonDecodeIntSlice(jsonStr string) []int {

	var data []int

	_ = json.Unmarshal([]byte(jsonStr), &data)

	return data
}

// JsonDecodeWithStruct 带结构体的json解码
//
//	temp := ReGeoRes{}
//	typeHelper.JsonDecodeWithStruct(tempStr,&temp)
func JsonDecodeWithStruct(jsonStr string, iStruct interface{}) error {
	return json.Unmarshal([]byte(jsonStr), &iStruct)
}

/**
 * @func: Float64SliceToInt float64 slice转换成 int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func Int2Str(number int) string {
	return strconv.Itoa(number)
}

// UintToStr uint 转换成 str
func UintToStr(number uint) string {
	return strconv.FormatUint(uint64(number), 10)
}

/**
 * @func: Float64SliceToInt float64 slice转换成 int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func Uint64ToStr(number uint64) string {
	return strconv.FormatUint(number, 10)
}

/**
 * @func: If 三元运算符
 * @author Wiidz
 * @date   2019-11-16
 */
func If(conditions bool, trueVal, falseVal interface{}) interface{} {
	if conditions {
		return trueVal
	}
	return falseVal
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

// 整形转换成字节
func Int2Bytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// 字节转换成整形
func Bytes2Int(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

// Float64ToInt float64转int
func Float64ToInt(number float64) int {
	return int(number)
}

/**
 * @func: Float64ToInt8 将字符串转为int8
 * @author Wiidz
 * @date   2019-11-16
 */
func Float64ToInt8(number float64) int8 {
	return int8(number)
}

/**
 * @func: Float64ToUint64 将字符串转为int8
 * @author Wiidz
 * @date   2019-11-16
 */
func Float64ToUint64(number float64) uint64 {
	return uint64(number)
}

// Float64ToInt64 float64转为int64
func Float64ToInt64(number float64) int64 {
	return int64(number)
}

// Int8ToStr int8转字符串
func Int8ToStr(number int8) string {
	return strconv.Itoa(int(number))
}
