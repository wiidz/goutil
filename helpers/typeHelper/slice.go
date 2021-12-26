package typeHelper

import (
	"reflect"
)


// ForceIntSlice 强制转换成int64切片
func ForceIntSlice(value interface{}) []int {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case "string":
		slice := ExplodeInt(value.(string), ",")
		return slice
	case "[]int":
		return value.([]int)
	case "[]uint64":
		return Uint64Slice2Int(value.([]uint64))
	case "[]float64":
		return Float64Slice2Int(value.([]float64))
	case "[]string":
		return StrSlice2Int(value.([]string))
	case "[]int8":
		return Int8Slice2Int(value.([]int8))
	default :
		return []int{}
	}
}

// ForceUint64Slice 强制转换成uint64切片
func ForceUint64Slice(value interface{}) []uint64 {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case "string":
		return ExplodeUint64(value.(string), ",")
	case "[]int":
		return IntSlice2Uint64(value.([]int))
	case "[]uint64":
		return value.([]uint64)
	case "[]float64":
		return Float64Slice2Uint64(value.([]float64))
	case "[]string":
		return StrSlice2Uint64(value.([]string))
	case "[]int8":
		return Int8Slice2Uint64(value.([]int8))
	default :
		return []uint64{}
	}
}

// ForceFloat64Slice 强制转换成float64切片
func ForceFloat64Slice(value interface{}) []float64 {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case "string":
		return ExplodeFloat64(value.(string), ",")
	case "[]int":
		return IntSlice2Float64(value.([]int))
	case "[]float64":
		return value.([]float64)
	case "[]uint64":
		return Uint64Slice2Float64(value.([]uint64))
	case "[]string":
		return StrSlice2Float64(value.([]string))
	case "[]int8":
		return Int8Slice2Float64(value.([]int8))
	default :
		return []float64{}
	}
}

// ForceStrSlice 强制转换成str切片
func ForceStrSlice(value interface{}) []string {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case "string":
		return ExplodeStr(value.(string), ",")
	case "[]int":
		return IntSlice2Str(value.([]int))
	case "[]float64":
		return Float64Slice2Str(value.([]float64))
	case "[]uint64":
		return Uint64Slice2Str(value.([]uint64))
	case "[]string":
		return value.([]string)
	case "[]int8":
		return Int8Slice2Str(value.([]int8))
	default :
		return []string{}
	}
}

// ForceInt8Slice 强制转换成int8切片
func ForceInt8Slice(value interface{}) []int8 {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case "string":
		return ExplodeInt8(value.(string), ",")
	case "[]int":
		return IntSlice2Int8(value.([]int))
	case "[]float64":
		return Float64Slice2Int8(value.([]float64))
	case "[]uint64":
		return Uint64Slice2Int8(value.([]uint64))
	case "[]string":
		return StrSlice2Int8(value.([]string))
	case "[]int8":
		return value.([]int8)
	default :
		return []int8{}
	}
}








// Int

// IntSlice2Int8  Int切片转int8切片
func IntSlice2Int8(data []int) []int8 {
	tmp := make([]int8, 0)
	for _, v := range data {
		tmp = append(tmp, int8(v))
	}
	return tmp
}
// IntSlice2Float64  Int切片转float64切片
func IntSlice2Float64(data []int) []float64 {
	tmp := make([]float64, 0)
	for _, v := range data {
		tmp = append(tmp, float64(v))
	}
	return tmp
}
// IntSlice2Str  Int切片转str切片
func IntSlice2Str(data []int) []string {
	tmp := make([]string, 0)
	for _, v := range data {
		tmp = append(tmp, Int2Str(v))
	}
	return tmp
}
// IntSlice2Uint64  Int切片转uint64切片
func IntSlice2Uint64(data []int) []uint64 {
	tmp := make([]uint64, 0)
	for _, v := range data {
		tmp = append(tmp, uint64(v))
	}
	return tmp
}


// Int8

// Int8Slice2Int  Int8切片转int切片
func Int8Slice2Int(data []int8) []int {
	tmp := make([]int, 0)
	for _, v := range data {
		tmp = append(tmp, int(v))
	}
	return tmp
}
// Int8Slice2Float64  Int8切片转float64切片
func Int8Slice2Float64(data []int8) []float64 {
	tmp := make([]float64, 0)
	for _, v := range data {
		tmp = append(tmp, float64(v))
	}
	return tmp
}
// Int8Slice2Str  Int8切片转str切片
func Int8Slice2Str(data []int8) []string {
	tmp := make([]string, 0)
	for _, v := range data {
		tmp = append(tmp, Int8ToStr(v))
	}
	return tmp
}
// Int8Slice2Uint64  Int8切片转uint64切片
func Int8Slice2Uint64(data []int8) []uint64 {
	tmp := make([]uint64, 0)
	for _, v := range data {
		tmp = append(tmp, uint64(v))
	}
	return tmp
}


// Float64

// Float64Slice2Int  Float64切片转int切片
func Float64Slice2Int(data []float64) []int {
	tmp := make([]int, 0)
	for _, v := range data {
		tmp = append(tmp, int(v))
	}
	return tmp
}
// Float64Slice2Int8  Float64切片转int8切片
func Float64Slice2Int8(data []float64) []int8 {
	tmp := make([]int8, 0)
	for _, v := range data {
		tmp = append(tmp, Float64ToInt8(v))
	}
	return tmp
}
// Float64Slice2Uint64  Float64切片转uint64切片
func Float64Slice2Uint64(data []float64) []uint64 {
	tmp := make([]uint64, 0)
	for _, v := range data {
		tmp = append(tmp, Float64ToUint64(v))
	}
	return tmp
}
// Float64Slice2Str  Float64切片转str切片
func Float64Slice2Str(data []float64) []string {
	tmp := make([]string, 0)
	for _, v := range data {
		tmp = append(tmp, Float64ToStr(v))
	}
	return tmp
}


// Uint64

// Uint64Slice2Int  Uint64切片转int切片
func Uint64Slice2Int(data []uint64) []int {
	tmp := make([]int, 0)
	for _, v := range data {
		tmp = append(tmp, int(v))
	}
	return tmp
}
// Uint64Slice2Int8  Uint64切片转int8切片
func Uint64Slice2Int8(data []uint64) []int8 {
	tmp := make([]int8, 0)
	for _, v := range data {
		tmp = append(tmp, int8(v))
	}
	return tmp
}
// Uint64Slice2Float64  Uint64切片转float64切片
func Uint64Slice2Float64(data []uint64) []float64 {
	tmp := make([]float64, 0)
	for _, v := range data {
		tmp = append(tmp, float64(v))
	}
	return tmp
}
// Uint64Slice2Str  Uint64切片转str切片
func Uint64Slice2Str(data []uint64) []string {
	tmp := make([]string, 0)
	for _, v := range data {
		tmp = append(tmp, Uint64ToStr(v))
	}
	return tmp
}


// Str

// StrSlice2Int  Str切片转int切片
func StrSlice2Int(data []string) []int {
	tmp := make([]int, 0)
	for _, v := range data {
		tmp = append(tmp, Str2Int(v))
	}
	return tmp
}
// StrSlice2Int8  Str切片转int8切片
func StrSlice2Int8(data []string) []int8 {
	tmp := make([]int8, 0)
	for _, v := range data {
		tmp = append(tmp, Str2Int8(v))
	}
	return tmp
}
// StrSlice2Float64  Str切片转float64切片
func StrSlice2Float64(data []string) []float64 {
	tmp := make([]float64, 0)
	for _, v := range data {
		tmp = append(tmp, Str2Float64(v))
	}
	return tmp
}
// StrSlice2Uint64  Str切片转uint64切片
func StrSlice2Uint64(data []string) []uint64 {
	tmp := make([]uint64, 0)
	for _, v := range data {
		tmp = append(tmp, Str2Uint64(v))
	}
	return tmp
}





/**
 * @func: StrSlice2InterfaceSlice 字符串slice转interface slice
 * @author Wiidz
 * @date   2019-11-16
 */
func StrSlice2InterfaceSlice(data []string) []interface{} {
	tmp := make([]interface{}, 0)
	for _, v := range data {
		tmp = append(tmp, v)
	}
	return tmp
}


// InterfaceSlice2MapSlice
func InterfaceSlice2MapSlice(inter []interface{}) []map[string]interface{} {

	tmp := make([]map[string]interface{}, 0)

	for _, v := range inter {

		tmp = append(tmp, v.(map[string]interface{}))

	}

	return tmp

}