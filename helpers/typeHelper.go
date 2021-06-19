package helpers

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"go/types"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type TypeHelper struct{}

/**
 * @func: Implode 将slice转换成字符串
 * @author Wiidz
 * @date   2019-11-16
 */
func (typeHelper *TypeHelper) Implode(data interface{}, glue string) string {
	var tmp []string
	for _, item := range data.([]interface{}) {
		tmp = append(tmp, typeHelper.ToString(item))
	}

	return strings.Join(tmp, glue)
}

/**
 * @func: Explode  字符串转slice, 接受混合类型
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Explode(data string, sep string) []interface{} {

	old := strings.Split(data, sep)

	newS := make([]interface{}, len(old))

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
func (*TypeHelper) ExplodeInt(data string, sep string) []int {
	old := strings.Split(data, sep)
	newS := make([]int, len(old))
	for i, v := range old {
		newS[i], _ = strconv.Atoi(v)
	}
	return newS
}

/**
 * @func: Explode  字符串转int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) ExplodeFloat64(data string, sep string) []float64 {
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
func (*TypeHelper) ExplodeInt64(data string, sep string) []int64 {
	old := strings.Split(data, sep)
	newS := make([]int64, len(old))
	for i, v := range old {
		newS[i], _ = strconv.ParseInt(v, 64, 10)
	}
	return newS
}

/**
 * @func: GetType 获取目标的数据类型
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) GetType(params interface{}) string {
	//数据初始化
	v := reflect.ValueOf(params)
	//获取传递参数类型
	vT := v.Type()
	//类型名称对比
	return vT.String()
}

/**
 * @func: Empty 判断目标是否为空值或0
 * @author Wiidz
 * @date   2019-11-16
 */
func (typeHelper *TypeHelper) Empty(arg interface{}) bool {
	switch arg.(type) {

	case int:
		return typeHelper.If(arg.(int) == int(0), true, false).(bool)
	case int64:
		return typeHelper.If(arg.(int64) == int64(0), true, false).(bool)

	case float64:
		return typeHelper.If(arg.(float64) == float64(0.00), true, false).(bool)

	case []byte:

		return typeHelper.If(len(arg.([]byte)) == 0, true, false).(bool)

	case string:
		return typeHelper.If(arg.(string) == " " || arg.(string) == "" || arg.(string) == "0" || arg.(string) == "NULL", true, false).(bool)

	case map[string]interface{}:

		return typeHelper.If(len(arg.(map[string]interface{})) == 0, true, false).(bool)

	case []interface{}:

		return typeHelper.If(len(arg.([]interface{})) == 0, true, false).(bool)

	case []string:

		return typeHelper.If(len(arg.([]string)) == 0, true, false).(bool)

	case []int64:

		return typeHelper.If(len(arg.([]int64)) == 0, true, false).(bool)

	case []float64:

		return typeHelper.If(len(arg.([]float64)) == 0, true, false).(bool)

	case []int:

		return typeHelper.If(len(arg.([]int)) == 0, true, false).(bool)

	case []map[string]interface{}:

		return typeHelper.If(len(arg.([]map[string]interface{})) == 0, true, false).(bool)

	case types.Nil:

		return typeHelper.If(arg == nil, true, false).(bool)

	case bool:

		return typeHelper.If(!arg.(bool), true, false).(bool)

	default:

		return true
	}
}

/**
 * @func: IsType 判断目标是否为指定类型的值
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) IsType(needle interface{}, type_name string) bool {
	if reflect.TypeOf(needle).String() == type_name {
		return true
	}
	return false
}

/**
 * @func: StrSlice2InterfaceSlice 字符串slice转interface slice
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) StrSlice2InterfaceSlice(data []string) []interface{} {
	tmp := make([]interface{}, 0)
	for _, v := range data {
		tmp = append(tmp, v)
	}
	return tmp
}

/**
 * @func: ToString 将任何参数转换为字符串
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) ToString(data interface{}) string {
	switch data.(type) {
	case int:
		return strconv.Itoa(data.(int))
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
 * @func: Map2InterfaceSlice 将map转换成interface slice
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Map2InterfaceSlice(data map[string]interface{}) []interface{} {
	islice := []interface{}{}
	for _, v := range data {
		islice = append(islice, v)
	}
	return islice
}

/**
 * @func: Str2Int64 将字符串转为int64
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Str2Int64(str string) int64 {
	number, _ := strconv.ParseInt(str, 10, 64)
	return number
}

/**
 * @func: Str2Int 将字符串转为int
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Str2Int(str string) int {
	ints, _ := strconv.Atoi(str)
	return ints
}

/**
 * @func: Str2Int8 将字符串转为int8
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Str2Int8(str string) int8 {
	ints, _ := strconv.Atoi(str)
	return int8(ints)
}

/**
 * @func: Str2Int 将字符串转为int
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Float64ToStr(number float64) string {
	return strconv.FormatFloat(number, 'f', -1, 64)
}

/**
 * @func: Int64ToStr int64转为str
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Int64ToStr(number int64) string {
	str := strconv.FormatInt(number, 10)
	return str
}

/**
 * @func: Str2Float64 str转为float64
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Str2Float64(str string) float64 {
	v, _ := strconv.ParseFloat(str, 64)
	return v
}

/**
 * @func: JsonEncode 编码json
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) JsonEncode(data interface{}) (string, error) {
	res, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(res), err
}

/**
 * @func: JsonDecode 解码json
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) JsonDecode(json_str string) map[string]interface{} {

	var data map[string]interface{}

	json.Unmarshal([]byte(json_str), &data)

	return data
}

/**
 * @func: JsonDecode 解码json
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) JsonDecodeInt64Slice(json_str string) []int64 {

	var data []int64

	json.Unmarshal([]byte(json_str), &data)

	return data
}

/**
 * @func: JsonDecode 解码json
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) JsonDecodeWithStruct(json_str string, istruct interface{}) interface{} {
	json.Unmarshal([]byte(json_str), &istruct)
	return istruct
}

/**
 * @func: InterfaceSlice2MapSlice interface slice 转换 map slice
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) InterfaceSlice2MapSlice(inter []interface{}) []map[string]interface{} {

	tmp := make([]map[string]interface{}, 0)

	for _, v := range inter {

		tmp = append(tmp, v.(map[string]interface{}))

	}

	return tmp

}

/**
 * @func: Float64SliceToInt float64 slice转换成 int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Float64Slice2IntSlice(slice []interface{}) []int {
	int_slice := []int{}
	for _, v := range slice {
		//fmt.Println("v",v)
		//str:=strconv.FormatFloat(v.(float64),'E',-1,64)
		//fmt.Println("str",str)
		//int64,_:=strconv.ParseInt(str,10,64)
		int_slice = append(int_slice, int(int64(v.(float64))))
	}
	return int_slice
}

/**
 * @func: Float64SliceToInt float64 slice转换成 int slice
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) Int2Str(number int) string {
	return strconv.Itoa(number)
}

/**
 * @func: If 三元运算符
 * @author Wiidz
 * @date   2019-11-16
 */
func (*TypeHelper) If(conditions bool, trueVal, falseVal interface{}) interface{} {
	if conditions {
		return trueVal
	}
	return falseVal
}

func (*TypeHelper) Int64Slice2IntSlice(int64_slice []int64) []int {
	res := []int{}
	for _, v := range int64_slice {
		res = append(res, int(v))
	}
	return res
}

func (*TypeHelper) IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}


// 整形转换成字节
func (*TypeHelper)Int2Bytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
// 字节转换成整形
func (*TypeHelper)Bytes2Int(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}