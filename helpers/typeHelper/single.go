package typeHelper

import (
	"reflect"
	"strconv"
	"unsafe"
)


// ForceUint64 强制转换为Uint64
func ForceUint64(value interface{}) uint64 {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case string(Int):
		return uint64(value.(int))
	case string(Int8):
		return uint64(value.(int8))
	case string(Int16):
		return uint64(value.(int16))
	case string(Int32):
		return uint64(value.(int32))
	case string(Int64):
		return uint64(value.(int64))
	case string(Uint):
		return uint64(value.(uint))
	case string(Uint8):
		return uint64(value.(uint8))
	case string(Uint16):
		return uint64(value.(uint16))
	case string(Uint32):
		return uint64(value.(uint32))
	case string(Uint64):
		return value.(uint64)
	case string(Uintptr):
		temp := (*uint64)(unsafe.Pointer(&value))
		return *temp
	case string(Float32):
		return uint64(value.(float32))
	case string(Float64):
		return uint64(value.(float64))
	case string(String):
		number, _ := strconv.ParseUint(value.(string), 10, 64)
		return number
	default:
		return 0
	}
}

// ForceFloat64 强制转换为Float64
func ForceFloat64(value interface{}) float64 {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case string(Int):
		return float64(value.(int))
	case string(Int8):
		return float64(value.(int8))
	case string(Int16):
		return float64(value.(int16))
	case string(Int32):
		return float64(value.(int32))
	case string(Int64):
		return float64(value.(int64))
	case string(Uint):
		return float64(value.(uint))
	case string(Uint8):
		return float64(value.(uint8))
	case string(Uint16):
		return float64(value.(uint16))
	case string(Uint32):
		return float64(value.(uint32))
	case string(Uint64):
		return float64(value.(uint64))
	case string(Uintptr):
		temp := (*uint64)(unsafe.Pointer(&value))
		return float64(*temp)
	case string(Float32):
		return float64(value.(float32))
	case string(Float64):
		return value.(float64)
	case string(String):
		number, _ := strconv.ParseFloat(value.(string), 64)
		return number
	default:
		return 0
	}
}

// ForceInt 强制转换为int
func ForceInt(value interface{}) int {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case string(Int):
		return value.(int)
	case string(Int8):
		return int(value.(int8))
	case string(Int16):
		return int(value.(int16))
	case string(Int32):
		return int(value.(int32))
	case string(Int64):
		return int(value.(int64))
	case string(Uint):
		return int(value.(uint))
	case string(Uint8):
		return int(value.(uint8))
	case string(Uint16):
		return int(value.(uint16))
	case string(Uint32):
		return int(value.(uint32))
	case string(Uint64):
		return int(value.(uint64))
	case string(Uintptr):
		temp := (*uint64)(unsafe.Pointer(&value))
		return int(*temp)
	case string(Float32):
		return int(value.(float32))
	case string(Float64):
		return int(value.(float64))
	case string(String):
		temp, _ := strconv.Atoi(value.(string))
		return temp
	default:
		return 0
	}
}

// ForceInt8 强制转换为int8
func ForceInt8(value interface{}) int8 {
	valueType := reflect.TypeOf(value).String()
	switch valueType {
	case string(Int):
		return int8(value.(int))
	case string(Int8):
		return value.(int8)
	case string(Int16):
		return int8(value.(int16))
	case string(Int32):
		return int8(value.(int32))
	case string(Int64):
		return int8(value.(int64))
	case string(Uint):
		return int8(value.(uint))
	case string(Uint8):
		return int8(value.(uint8))
	case string(Uint16):
		return int8(value.(uint16))
	case string(Uint32):
		return int8(value.(uint32))
	case string(Uint64):
		return int8(value.(uint64))
	case string(Uintptr):
		temp := (*uint64)(unsafe.Pointer(&value))
		return int8(*temp)
	case string(Float32):
		return int8(value.(float32))
	case string(Float64):
		return int8(value.(float64))
	case string(String):
		temp, _ := strconv.Atoi(value.(string))
		return int8(temp)
	default:
		return 0
	}
}

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
