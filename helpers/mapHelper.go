package helpers

import (
	"fmt"
	"strconv"
)

type MapHelper struct{}

/**
 * @func: MapEqual 判断两个map内容相同
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) Equal(m1, m2 map[string]interface{}) bool {
	for k1, v1 := range m1 {
		if v2, has := m2[k1]; has {
			if v1 != v2 {
				return false
			}
		} else {
			return false
		}
	}
	for k2, v2 := range m2 {
		if v1, has := m1[k2]; has {
			if v1 != v2 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

/**
 * @func: Exist  判断map中是否包含某键值
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) Exist(data map[string]interface{}, key string) bool {

	if _, ok := data[key]; ok {
		return true
	}

	return false
}

/**
 * @func: FilterMap 过滤map 删除fields以外的键
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) Filter(target map[string]interface{}, fields []string) {
	var sliceHelper SliceHelper
	for k, _ := range target {
		index := sliceHelper.IndexOfStrSlice(k, fields)
		if index != -1 {
			sliceHelper.StrSliceDelete(fields, index)
			continue
		}
		delete(target, k)
	}
}

/**
 * @func: HasIllegalKey 判断map是否包含规定以外的键
 * @author Wiidz
 * @date   2019-11-06
 */
func (*MapHelper) HasIllegalKey(target map[string]interface{}, fields []string) bool {
	var sliceHelper SliceHelper
	for k, _ := range target {
		index := sliceHelper.IndexOfStrSlice(k, fields)
		if index != -1 {
			sliceHelper.StrSliceDelete(fields, index)
			continue
		}
		fmt.Println("illegal field", k)
		return true
	}
	return false
}

/**
 * @func: GetKeys 从 map[string]interface{} 中获取键名
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) GetKeys(imap map[string]interface{}) []string {
	var tmp []string
	if len(imap) > 0 {
		for k, _ := range imap {
			tmp = append(tmp, k)
		}
	}
	return tmp
}

/**
 * @func: GetKeysPlus  从 map[string]map[string]interface{} 中获取一级键名
 * @author Wiidz
 * @date  2019-11-16
 */
func (*MapHelper) GetKeysPlus(imap map[string]map[string]interface{}) []string {
	var keyNames []string
	if len(imap) > 0 {
		for k, _ := range imap {
			keyNames = append(keyNames, k)
		}
	}
	return keyNames
}

/**
 * @func: GetKeys 从 map[string]interface{} 中获取整数类型的键名
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) GetIntKeys(imap map[string]interface{}) []int {
	var tmp []int
	if len(imap) > 0 {
		for k, _ := range imap {
			key, _ := strconv.Atoi(k)
			tmp = append(tmp, key)
		}
	}
	return tmp
}

/**
 * @func: GetKeys 从 map[string]interface{} 中获取整数类型的键名
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) GetIntKeysFromIntKeyMap(imap map[int]bool) []int {
	var tmp []int
	if len(imap) > 0 {
		for k, _ := range imap {
			tmp = append(tmp, k)
		}
	}
	return tmp
}

/**
 * @func: GetKeys 从 map[string]interface{} 中获取int64类型的键名
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) GetInt64Keys(imap map[string]interface{}) []int64 {
	var tmp []int64
	if len(imap) > 0 {
		for k, _ := range imap {
			key, _ := strconv.ParseInt(k, 64, 10)
			tmp = append(tmp, key)
		}
	}
	return tmp
}

/**
 * @func: GetKeys 从 map[string]interface{} 中获取float64类型的键名
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) GetFloat64Keys(imap map[string]interface{}) []float64 {
	var tmp []float64
	if len(imap) > 0 {
		for k, _ := range imap {
			key, _ := strconv.ParseFloat(k, 64)
			tmp = append(tmp, key)
		}
	}
	return tmp
}

/**
 * @func: GetValues  从 map[string]interface{} 中输出所有的键值到一个[]interface{}数组中
 * @author Wiidz
 * @date   2019-11-16
 */
func (*MapHelper) GetValues(imap map[string]interface{}) []interface{} {
	var tmp []interface{}
	if len(imap) > 0 {
		for _, v := range imap {
			tmp = append(tmp, v)
		}
	}
	return tmp
}
