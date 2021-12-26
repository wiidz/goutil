package typeHelper

import (
	"strconv"
	"unsafe"
	"reflect"
)

// ForceString 强制转换为string
func ForceString(value interface{}) string {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case string(Int):
		return strconv.Itoa(value.(int))
	case string(Int8):
		return strconv.Itoa(int(value.(int8)))
	case string(Int16):
		return strconv.Itoa(int(value.(int16)))
	case string(Int32):
		return strconv.Itoa(int(value.(int32)))
	case string(Int64):
		return strconv.FormatInt(value.(int64), 10)
	case string(Uint):
		return strconv.FormatUint(uint64(value.(uint)), 10)
	case string(Uint8):
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case string(Uint16):
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case string(Uint32):
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case string(Uint64):
		return strconv.FormatUint(value.(uint64), 10)
	case string(Uintptr):
		temp := (*uint64)(unsafe.Pointer(&value))
		return strconv.FormatUint(*temp, 10)
	case string(Float32):
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32)
	case string(Float64):
		return strconv.FormatFloat(value.(float64), 'f', -1, 32)
	case string(String):
		return value.(string)
	default:
		return ""
	}
}
