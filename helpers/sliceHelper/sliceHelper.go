package sliceHelper

import (
	"fmt"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

/**
 * @func: GetRange  类似php的函数，获取范围内的整数集合
 * @author Wiidz
 * @date   2019-11-06
 */
func GetRange(min float64, max float64, step float64) []float64 {
	arr := make([]float64, 0)
	for min < max {
		arr = append(arr, min)
		min += step
	}
	return arr
}

/**
 * @func: ArrayIntersect 取出两个slice的交集
 * @author Wiidz
 * @date   2019-11-06
 */
func Intersect(arr_1 []interface{}, arr_2 []interface{}) []interface{} {
	m := make(map[interface{}]int)
	intersectArr := make([]interface{}, 0)
	//【1】统计arr_1中值出现的次数
	for _, v := range arr_1 {
		m[v]++
	}
	for _, v := range arr_2 {
		times, _ := m[v]
		if times != 0 {
			intersectArr = append(intersectArr, v)
		}
	}
	return intersectArr
}

/**
 * @func: ArrayIntersect 取出两个slice的交集
 * @author Wiidz
 * @date   2019-11-06
 */
func IntersectInt(arr_1 []int64, arr_2 []int64) []int64 {
	m := make(map[interface{}]int64)
	intersectArr := make([]int64, 0)
	//【1】统计arr_1中值出现的次数
	for _, v := range arr_1 {
		m[v]++
	}
	for _, v := range arr_2 {
		times, _ := m[v]
		if times != 0 {
			intersectArr = append(intersectArr, v)
		}
	}
	return intersectArr
}

/**
 * @func: ArrayIntersect 取出两个slice的交集
 * @author Wiidz
 * @date   2019-11-06
 */
func IntersectFloat64(arr_1 []float64, arr_2 []float64) []float64 {
	m := make(map[interface{}]float64)
	intersectArr := make([]float64, 0)
	//【1】统计arr_1中值出现的次数
	for _, v := range arr_1 {
		m[v]++
	}
	for _, v := range arr_2 {
		times, _ := m[v]
		if times != 0 {
			intersectArr = append(intersectArr, v)
		}
	}
	return intersectArr
}

/**
 * @func: ArrayIntersect 取出两个slice的差集
 * @author Wiidz
 * @date   2019-11-06
 */
func Diffrence(arr_1 []int64, arr_2 []int64) []int64 {
	m := make(map[interface{}]int64)
	diffrence := make([]int64, 0)
	//【1】统计arr_1中值出现的次数
	for _, v := range arr_1 {
		m[v]++
	}
	for _, v := range arr_2 {
		times, _ := m[v]
		if times == 0 {
			diffrence = append(diffrence, v)
		}
	}
	return diffrence
}

/**
 * @func: ArrayIntersect 取出两个slice的差集
 * @author Wiidz
 * @date   2019-11-06
 */
func DiffrenceFloat64(arr_1 []float64, arr_2 []float64) []float64 {
	m := make(map[interface{}]float64)
	diffrence := make([]float64, 0)
	//【1】统计arr_1中值出现的次数
	for _, v := range arr_1 {
		m[v]++
	}
	for _, v := range arr_2 {
		times, _ := m[v]
		if times == 0 {
			diffrence = append(diffrence, v)
		}
	}
	return diffrence
}

/**
 * @func: sliceMerge 合并两个slice
 * @author Wiidz
 * @date   2019-11-06
 */
func Merge(arr_1 []map[string]float64, arr_2 []map[string]float64) []map[string]float64 {
	for _, v := range arr_2 {
		arr_1 = append(arr_1, v)
	}
	return arr_1
}

/**
 * @func: SliceDelete 返回新slice（不包含指定的那个index键）
 * @author Wiidz
 * @date   2019-11-06
 */
func SliceDelete(slice []interface{}, index int) []interface{} {
	return append(slice[:index], slice[index+1:]...)
}

/**
 * @func: SliceDelete 返回新slice（不包含指定的那个index键）
 * @author Wiidz
 * @date   2019-11-06
 */
func IntSliceDelete(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

/**
 * @func: SliceDelete 返回新slice（不包含指定的那个index键）
 * @author Wiidz
 * @date   2019-11-06
 */
func Int64SliceDelete(slice []int64, index int) []int64 {
	return append(slice[:index], slice[index+1:]...)
}

/**
 * @func: SliceDelete 返回新slice（不包含指定的那个index键）
 * @author Wiidz
 * @date   2019-11-06
 */
func Float64SliceDelete(slice []float64, index int) []float64 {
	return append(slice[:index], slice[index+1:]...)
}

/**
 * @func: StrSliceDelete 返回新slice（不包含指定的那个index键）
 * @author Wiidz
 * @date   2019-11-06
 */
func StrSliceDelete(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

/**
 * @func: StrSliceDelete 返回新slice（不包含指定的那个index键）
 * @author Wiidz
 * @date   2019-11-06
 */
func MapSliceDelete(slice []map[string]interface{}, index int) []map[string]interface{} {
	return append(slice[:index], slice[index+1:]...)
}

/**
 * @func: Float64SliceSum float64切片求和
 * @author Wiidz
 * @date   2019-11-06
 */
func Float64SliceSum(arr []float64) float64 {
	sum := float64(0)
	for _, v := range arr {
		sum += v
	}
	return sum
}

/**
 * @func: InterfaceSliceSum float64切片求和
 * @author Wiidz
 * @date   2019-11-06
 */
func InterfaceSliceSum(arr []interface{}) float64 {
	sum := float64(0)
	for _, v := range arr {
		sum += v.(float64)
	}
	return sum
}

/**
 * @func: NarrowSlice 缩短slice，取中间部分
 * @author Wiidz
 * @date   2019-11-06
 */
func NarrowSlice(arr []map[string]interface{}, amount int) []map[string]interface{} {
	length := len(arr)
	flag_length, flag_amount := length%2, amount%2
	half_length, half_amount := 0, 0
	switch {
	case length <= amount:
		return arr
	case flag_length == 0:
		half_length = (length + 1) / 2
		if flag_amount == 0 {
			half_amount = amount / 2
		} else {
			half_amount = (amount - 1) / 2
		}
	default:
		half_length = (length + 1) / 2
		if flag_amount == 0 {
			half_amount = amount / 2
		} else {
			half_amount = (amount + 1) / 2
		}
	}
	return arr[half_length-half_amount : half_length+half_amount]
}

// UniqueMapByUint64Key 根据某一个uint64键来去重
func UniqueMapByUint64Key(sourceSlice []map[string]interface{}, keyName string) []map[string]interface{} {
	uint64Map := map[uint64]bool{}
	var handledSlice []map[string]interface{}
	for k := range sourceSlice {
		var key = sourceSlice[k][keyName].(uint64)
		if _, ok := uint64Map[key]; ok {
			continue
		}
		handledSlice = append(handledSlice, sourceSlice[k])
		uint64Map[key] = true
	}
	return handledSlice
}

// UniqueMapByFloat64Key 根据某一个float64键来去重
func UniqueMapByFloat64Key(sourceSlice []map[string]interface{}, keyName string) []map[string]interface{} {
	float64Map := map[float64]bool{}
	var handledSlice []map[string]interface{}
	for k := range sourceSlice {
		var key = sourceSlice[k][keyName].(float64)
		if _, ok := float64Map[key]; ok {
			continue
		}
		handledSlice = append(handledSlice, sourceSlice[k])
		float64Map[key] = true
	}
	return handledSlice
}

// UniqueMapByStrKey 根据某一个str键来去重
func UniqueMapByStrKey(sourceSlice []map[string]interface{}, keyName string) []map[string]interface{} {
	strMap := map[string]bool{}
	var handledSlice []map[string]interface{}
	for k := range sourceSlice {
		var key = sourceSlice[k][keyName].(string)
		if _, ok := strMap[key]; ok {
			continue
		}
		handledSlice = append(handledSlice, sourceSlice[k])
		strMap[key] = true
	}
	return handledSlice
}

// UniqueMapByIntKey 根据某一个int键来去重
func UniqueMapByIntKey(sourceSlice []map[string]interface{}, keyName string) []map[string]interface{} {
	intMap := map[int]bool{}
	var handledSlice []map[string]interface{}
	for k := range sourceSlice {
		var key = sourceSlice[k][keyName].(int)
		if _, ok := intMap[key]; ok {
			continue
		}
		handledSlice = append(handledSlice, sourceSlice[k])
		intMap[key] = true
	}
	return handledSlice
}

// UniqueMapByInt64Key 根据某一个int64键来去重
func UniqueMapByInt64Key(sourceSlice []map[string]interface{}, keyName string) []map[string]interface{} {
	int64Map := map[int64]bool{}
	var handledSlice []map[string]interface{}
	for k := range sourceSlice {
		var key = sourceSlice[k][keyName].(int64)
		if _, ok := int64Map[key]; ok {
			continue
		}
		handledSlice = append(handledSlice, sourceSlice[k])
		int64Map[key] = true
	}
	return handledSlice
}

/**
 * @func: UniqueIntSlice int slice去重
 * @author Wiidz
 * @date   2019-11-06
 */
func UniqueIntSlice(slc []int) []int {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return UniqueByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return UniqueByMap(slc)
	}
}

// UniqueInt8Slice int8 slice去重
func UniqueInt8Slice(slc []int8) []int8 {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return UniqueByLoopInt8(slc)
	} else {
		// 大于的时候，通过map来过滤
		return UniqueByLoopInt8(slc)
	}
}

// UniqueUint64Slice int slice去重
func UniqueUint64Slice(slc []uint64) []uint64 {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return UniqueByLoopUint64(slc)
	} else {
		// 大于的时候，通过map来过滤
		return UniqueByMapUint64(slc)
	}
}

// UniqueByLoop 通过两重循环过滤重复元素
func UniqueByLoop(slc []int) []int {
	result := []int{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

func UniqueByLoopInt8(slc []int8) []int8 {
	result := []int8{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// UniqueByLoopUint64 通过两重循环过滤重复元素
func UniqueByLoopUint64(slc []uint64) []uint64 {
	result := []uint64{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

/**
 * @func: UniqueInterface interface slice去重
 * @author Wiidz
 * @date   2019-11-06
 */
func UniqueInterface(slc []interface{}) []interface{} {
	result := []interface{}{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

/**
 * @func: UniqueStrSlice str slice去重
 * @author Wiidz
 * @date   2019-11-16
 */
func UniqueStrSlice(strSlice []string) []string {
	result := make([]string, 0, len(strSlice))
	temp := map[string]struct{}{}
	for _, item := range strSlice {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

/**
 * @func: UniqueByMap  通过map主键唯一的特性过滤重复元素
 * @author Wiidz
 * @date   2019-11-16
 */
func UniqueByMap(slc []int) []int {
	result := []int{}
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

/**
 * @func: UniqueByMap  通过map主键唯一的特性过滤重复元素
 * @author Wiidz
 * @date   2019-11-16
 */
func UniqueByMapUint64(slc []uint64) []uint64 {
	result := []uint64{}
	tempMap := map[uint64]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

/**
 * @func: ArrayFilter 不知道
 * @author Wiidz
 * @date   2019-11-06
 */
func ArrayFilter(a []interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	fmt.Printf("va : %+v\n", va)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

/**
 * @func: MapSliceReverse 将map slice倒过来排
 * @author Wiidz
 * @date   2019-11-16
 */
func MapSliceReverse(arr []map[string]interface{}) ([]map[string]interface{}, error) {
	lenArr := len(arr)
	if lenArr == 0 {
		return arr, nil
	}
	var newArr []map[string]interface{}
	for i := lenArr - 1; i >= 0; i-- {
		newArr = append(newArr, arr[i])
	}
	return newArr, nil
}

/**
 * @func: GetParentIds 递归查找父级逗号隔开从高到低
 * @author Wiidz
 * @date   2019-11-16
 */
func GetParentIds(children_id, result_str, key string, original_data []map[string]interface{}) string {

	str_inn := result_str

	var tmp map[string]interface{}

	for _, val := range original_data {

		if val["id"] == children_id {

			tmp = val
		}

	}

	if typeHelper.Empty(tmp[key]) {

		str_inn = ""

	} else {

		str_inn = tmp["id"].(string) + "," + result_str

	}

	if typeHelper.Empty(tmp[key]) {

		return GetParentIds(tmp[key].(string), str_inn, key, original_data)

	}

	return strings.TrimRight(result_str, ",")

}

/**
 * @func: Paginator 组装分页
 * @author Wiidz
 * @date   2019-11-16
 */
func Paginator(page, pageSize int, data []map[string]interface{}) []map[string]interface{} {

	page = page + 1

	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 10
	}

	var tmpPageSize int

	size := len(data)

	span := (page - 1) * pageSize

	//超出了范围
	if span >= size {

		return make([]map[string]interface{}, 0)

		//不满整个长度需要特使处理

	} else if size-span < pageSize {

		tmpPageSize = size - span

		start := span

		end := start + tmpPageSize

		return data[start:end]

	} else {

		tmpPageSize = pageSize

	}

	start := (page - 1) * tmpPageSize

	if start < 0 || start == 0 {

		start = 0
	}

	end := start + tmpPageSize

	return data[start:end]
}

/**
 * @func: ArrayGroupByMapKey 以一个数组map 中的某一个key  进行分组
 * @author Wiidz
 * @date   2019-11-16
 */
func ArrayGroupByMapKey(key string, list []map[string]interface{}) [][]map[string]interface{} {

	returnData := make([][]map[string]interface{}, 0)
	i := 0
	var j int
	for {
		if i >= len(list) {
			break
		}
		for j = i + 1; j < len(list) && list[i][key] == list[j][key]; j++ {
		}

		returnData = append(returnData, list[i:j])
		i = j
	}

	return returnData
}

// Exist 判断slice中是否存在needle
func Exist(needle interface{}, scopeSlice interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range scopeSlice.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range scopeSlice.([]int) {
			if key == item {
				return true
			}
		}
	case int8:
		for _, item := range scopeSlice.([]int8) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range scopeSlice.([]int64) {
			if key == item {
				return true
			}
		}
	case float64:
		for _, item := range scopeSlice.([]float64) {
			if key == item {
				return true
			}
		}
	case uint64:
		for _, item := range scopeSlice.([]uint64) {
			if key == item {
				return true
			}
		}

	default:
		return false
	}
	return false
}

/**
 * @func: GetValuesFromInterface  从interface slice中获取键值slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetValuesFromInterfaceSlice(d []interface{}, column_key string) []interface{} {
	nd := make([]interface{}, 0)
	for _, v := range d {
		nd = append(nd, v.(map[string]interface{})[column_key])
	}
	return nd
}

/**
 * @func: GetValuesFromMapSlice  从map slice中获取键值slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetValuesFromMapSlice(d []map[string]interface{}, column_key string) []interface{} {
	nd := make([]interface{}, 0)
	for _, v := range d {
		nd = append(nd, v[column_key])
	}
	return nd
}

/**
 * @func: SkuList  从map slice中获取键值 int64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetInt64FromMapSlice(d []map[string]interface{}, column_key string) []int64 {
	nd := make([]int64, 0)
	for _, v := range d {
		nd = append(nd, v[column_key].(int64))
	}
	return nd
}

/**
 * @func: SkuList  从map slice中获取键值 int64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetIntFromMapSliceByInt64(d []map[string]interface{}, column_key string) []int {
	nd := make([]int, 0)
	for _, v := range d {
		temp := v[column_key].(int64)
		nd = append(nd, *(*int)(unsafe.Pointer(&temp)))
	}
	return nd
}

/**
 * @func: GetIntFromMapSlice  从map slice中获取键值 int slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetIntFromMapSlice(d []map[string]interface{}, column_key string) []int {
	nd := make([]int, 0)
	for _, v := range d {
		nd = append(nd, v[column_key].(int))
	}
	return nd
}

/**
 * @func: GetFloat64FromMapSlice  从map slice中获取键值float64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetFloat64FromMapSlice(d []map[string]interface{}, column_key string) []float64 {
	nd := make([]float64, 0)
	for _, v := range d {
		nd = append(nd, v[column_key].(float64))
	}
	return nd
}

/**
 * @func: GetFloat64FromMapSlice  从map slice中获取键值float64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetFloat64FromInterfaceSlice(d []interface{}, column_key string) []float64 {
	nd := make([]float64, 0)
	for _, v := range d {
		nd = append(nd, v.(map[string]interface{})[column_key].(float64))
	}
	return nd
}

/**
 * @func: GetFloat64FromMapSlice  从map slice中获取键值float64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetInt64FromInterfaceSlice(d []interface{}, column_key string) []int64 {
	nd := make([]int64, 0)
	for _, v := range d {
		nd = append(nd, v.(map[string]interface{})[column_key].(int64))
	}
	return nd
}

/**
 * @func: GetFloat64FromMapSlice  从map slice中获取键值float64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetIntFromInterfaceSlice(d []interface{}, column_key string) []int {
	nd := make([]int, 0)
	for _, v := range d {
		nd = append(nd, int(v.(map[string]interface{})[column_key].(float64)))
	}
	return nd
}

/**
 * @func: IndexOfStrSlice 从str slice中，寻找指定内容的键值
 * @author Wiidz
 * @date   2019-11-06
 */
func IndexOfStrSlice(needle string, fields_slice []string) int {
	for k, v := range fields_slice {
		if v == needle {
			return k
		}
		continue
	}
	return -1
}

/**
 * @func: IndexOfOld 从slice中，寻找指定内容的键值
 * @author Wiidz
 * @date   2019-11-06
 */
func IndexOfOld(needle interface{}, scopeSlice interface{}) int {
	switch key := needle.(type) {
	case string:
		for index, item := range scopeSlice.([]string) {
			if key == item {
				return index
			}
		}
	case int:
		for index, item := range scopeSlice.([]int) {
			if key == item {
				return index
			}
		}
	case int64:
		for index, item := range scopeSlice.([]int64) {
			if key == item {
				return index
			}
		}
	case uint64:
		for index, item := range scopeSlice.([]uint64) {
			if key == item {
				return index
			}
		}
	case int8:
		for index, item := range scopeSlice.([]int8) {
			if key == item {
				return index
			}
		}
	case uint:
		for index, item := range scopeSlice.([]uint) {
			if key == item {
				return index
			}
		}
	default:
		return -1
	}

	return -1
}

// IndexOf 从slice中，寻找指定内容的键值
func IndexOf(needle interface{}, scopeSlice []interface{}) int {
	for index, item := range scopeSlice {
		if needle == item {
			return index
		}
	}
	return -1
}

func Slice2MapByInt64ColumnAsKey(data []map[string]interface{}, index string) map[string]interface{} {
	res := map[string]interface{}{}
	//temp:=map[string]interface{}{}
	for _, v := range data {
		index := strconv.FormatInt(v[index].(int64), 10)
		res[index] = v
	}
	return res
}

func MapSlice2SimpleMap(data []map[string]interface{}, index, key string) map[string]interface{} {
	res := map[string]interface{}{}
	//temp:=map[string]interface{}{}
	for _, v := range data {
		index := strconv.FormatInt(v[index].(int64), 10)
		res[index] = v[key]
	}
	return res
}
func InterfaceSlice2SimpleMap(data []interface{}, index, key string) map[string]interface{} {
	res := map[string]interface{}{}
	//temp:=map[string]interface{}{}
	for _, v := range data {
		index := strconv.FormatFloat(v.(map[string]interface{})[index].(float64), 'f', -1, 64)
		res[index] = v.(map[string]interface{})[key]
	}
	return res
}

func Int64Slice2Float64Slice(data []int64) []float64 {
	res := []float64{}
	for _, v := range data {
		res = append(res, float64(v))
	}
	return res
}

/**
 * @func: SkuList  从map slice中获取键值 int64 slice
 * @author Wiidz
 * @date   2019-11-06
 */
func GetFloat64FromInt64MapSlice(d []map[string]interface{}, column_key string) []float64 {
	nd := make([]float64, 0)
	for _, v := range d {
		nd = append(nd, float64(v[column_key].(int64)))
	}
	return nd
}

func Join(islice []string, letter string) string {
	return strings.Join(islice, letter)
}

func JoinInterfaceSlice(islice []interface{}, letter string) string {
	str := ""
	for _, v := range islice {
		str += v.(string) + letter
	}
	return str[0 : len(str)-len(letter)]
}

/**
 * @func: Exsit 判断slice中是否存在needle
 * @author Wiidz
 * @date   2019-11-16
 */
func ExistInt(needle int, scopeSlice []int) bool {
	for _, item := range scopeSlice {
		if needle == item {
			return true
		}
	}
	return false
}

/**
 * @func: Exsit 判断slice中是否存在needle
 * @author Wiidz
 * @date   2019-11-16
 */
func ExistUint64(needle uint64, scopeSlice []uint64) bool {
	for _, item := range scopeSlice {
		if needle == item {
			return true
		}
	}
	return false
}

// IsElementEqual 判断参数的attrIDs是否和数据库里的要求吻合
func IsElementEqual(requiredIDs, paramsIDs []uint64) bool {

	if len(requiredIDs) != len(paramsIDs) {
		return false
	}

	for _, v := range paramsIDs {
		if Exist(v, requiredIDs) {
			continue
		}
		return false
	}
	return true
}

// Difference 获取s2比s1差异的部分
func Difference(s1, s2 []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	var diff []interface{}

	// 将s2中的元素添加到map中
	for _, item := range s2 {
		m[item] = true
	}

	// 遍历s1，检查元素是否在s2的map中
	for _, item := range s1 {
		if _, found := m[item]; !found {
			// 如果元素不在map中，则添加到差集结果中
			diff = append(diff, item)
		}
	}

	return diff
}

// Intersection 获取s1和s2的交集
func Intersection(s1, s2 []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	var intersection []interface{}

	// 遍历s1并将元素添加到map中
	for _, item := range s1 {
		m[item] = true
	}

	// 遍历s2，检查元素是否在s1的map中
	for _, item := range s2 {
		// 使用反射来比较元素
		if _, ok := m[item]; ok {
			// 如果元素在map中，则添加到交集结果中
			intersection = append(intersection, item)
			// 可以选择从map中删除该元素，以避免重复结果
			delete(m, item)
		}
	}

	return intersection
}

// Union 获取s1和s2的并集
func Union(s1, s2 []interface{}) []interface{} {
	m := make(map[interface{}]bool)
	var union []interface{}

	// 将s1中的所有元素添加到map中
	for _, item := range s1 {
		m[item] = true
	}

	// 遍历s2，将不在map中的元素添加到map和union中
	for _, item := range s2 {
		if _, found := m[item]; !found {
			m[item] = true
			union = append(union, item)
		}
	}

	// 将s1中剩余的元素添加到union中
	for item := range m {
		union = append(union, item)
	}

	return union
}
