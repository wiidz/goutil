package luaHelper

import (
	"encoding/json"
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// ReadMapSliceFromFile 从lua文件中读取数据并转换为[]map[string]interface{}
// 记得在外面关闭 defer L.Close()
func ReadMapSliceFromFile(filePath string, dataName string, doCamel bool) (L *lua.LState, mapSlice []map[string]lua.LValue, err error) {

	//【1】读byte
	file, _ := os.Open(filePath)
	defer file.Close()
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	//【2】创建 Lua 解释器
	L = lua.NewState()

	//【3】执行 Lua 脚本，判断有无错误
	if err = L.DoString(string(byteData)); err != nil {
		fmt.Println("Error executing Lua script:", err)
		return
	}

	//【4】获取 Lua 全局变量
	tableItem := L.GetGlobal(dataName)
	if tableItem == lua.LNil {
		err = errors.New("Lua global variable " + dataName + " not found")
		return
	}

	//【5】解析 Lua 表格数据
	luaTable := tableItem.(*lua.LTable)

	//【6】转换
	mapSlice = []map[string]lua.LValue{}
	luaTable.ForEach(func(_, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			row := extractRow(tbl, doCamel)
			mapSlice = append(mapSlice, row)
		}
	})

	return
}

// extractRow 提取一行数据为map[string]lua.LValue
func extractRow(tbl *lua.LTable, doCamel bool) map[string]lua.LValue {
	row := make(map[string]lua.LValue)
	tbl.ForEach(func(key, value lua.LValue) {
		if key.Type() == lua.LTString {
			var columnName = key.String()
			if doCamel {
				columnName = camelToSnake(key.String())
			}
			row[columnName] = value
		}
	})
	return row
}

// camelToSnake 驼峰式命名转换为下划线间隔
func camelToSnake(s string) string {
	var re = regexp.MustCompile(`(?m)([a-z])([A-Z])`)
	return strings.ToLower(re.ReplaceAllString(s, "${1}_${2}"))
}

// LuaValueToInterface 将 Lua 值转换为对应的 Go 数据类型
func LuaValueToInterface(lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTString:
		return lv.String()
	case lua.LTNumber:
		return lv.(lua.LNumber)
	case lua.LTBool:
		return lv.(lua.LBool)
	case lua.LTTable:
		tbl := lv.(*lua.LTable)
		// 不这么做会返回一个带有Metatable的空map
		if tbl.Len() == 0 {
			return make(map[string]interface{})
		}
		goMap := make(map[string]interface{})
		tbl.ForEach(func(key, value lua.LValue) {
			goMap[key.String()] = LuaValueToInterface(value)
		})
		return goMap
	default:
		return lv.String()
	}
}

// LuaValueToInterfaceNoTable 将 Lua 值转换为对应的 Go 数据类型(舍弃table)
func LuaValueToInterfaceNoTable(L *lua.LState, lv lua.LValue) interface{} {
	switch lv.Type() {
	case lua.LTString:
		return lv.String()
	case lua.LTNumber:
		return lv.(lua.LNumber)
	case lua.LTBool:
		return lv.(lua.LBool)
	case lua.LTTable:
		return tableToJSON(L, lv.(*lua.LTable))
	default:
		return lv.String()
	}
}

func LuaValueToGoValue(lv lua.LValue) reflect.Value {
	switch lv.Type() {
	case lua.LTString:
		return reflect.ValueOf(lv.String())
	case lua.LTNumber:
		return reflect.ValueOf(float64(lv.(lua.LNumber)))
	case lua.LTBool:
		return reflect.ValueOf(bool(lv.(lua.LBool)))
	case lua.LTTable:
		tbl := lv.(*lua.LTable)
		// 判断是否为标准数组（key为1..N且没有别的key）
		arrayLen := tbl.Len()
		// 记录表项总数（用于区分kv类型table和稀疏/混合型table）
		itemCount := 0
		isArray := true
		arrayValues := make([]interface{}, arrayLen)

		tbl.ForEach(func(key, value lua.LValue) {
			itemCount++
			if key.Type() == lua.LTNumber {
				kf := float64(key.(lua.LNumber))
				ki := int(key.(lua.LNumber))
				// 检查key必须是1~arrayLen的整数，并且没有小数部分
				if kf == float64(ki) && ki >= 1 && ki <= arrayLen {
					arrayValues[ki-1] = LuaValueToGoValue(value).Interface()
				} else {
					isArray = false
				}
			} else {
				isArray = false
			}
		})

		// 必须所有key都是连续1..N 且数量等于Len才认定为数组
		if isArray && itemCount == arrayLen && arrayLen > 0 {
			return reflect.ValueOf(arrayValues)
		} else {
			goMap := make(map[string]interface{})
			tbl.ForEach(func(key, value lua.LValue) {
				goMap[key.String()] = LuaValueToGoValue(value).Interface()
			})
			return reflect.ValueOf(goMap)
		}
	case lua.LTNil:
		return reflect.Zero(reflect.TypeOf((*interface{})(nil)).Elem())
	default:
		return reflect.ValueOf(lv.String())
	}
}

// tableToJSON 将 Lua 表转换为 JSON 格式的字符串
func tableToJSON(L *lua.LState, table *lua.LTable) string {
	result := make(map[string]interface{})
	table.ForEach(func(key, value lua.LValue) {
		keyStr := fmt.Sprintf("%v", key)
		result[keyStr] = LuaValueToInterfaceNoTable(L, value)
	})
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf("Error converting table to JSON: %v", err)
	}
	return string(jsonBytes)
}

// ReadMapStringSliceFromFile 从lua文件中读取数据并转换为[]map[string]string
// 记得在外面关闭 defer L.Close()
func ReadMapStringSliceFromFile(filePath string, dataName string, doCamel bool) (L *lua.LState, mapSlice []map[string]string, err error) {

	//【1】读byte
	file, _ := os.Open(filePath)
	defer file.Close()
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	//【2】创建 Lua 解释器
	L = lua.NewState()

	//【3】执行 Lua 脚本，判断有无错误
	if err = L.DoString(string(byteData)); err != nil {
		fmt.Println("Error executing Lua script:", err)
		return
	}

	//【4】获取 Lua 全局变量
	tableItem := L.GetGlobal(dataName)
	if tableItem == lua.LNil {
		err = errors.New("Lua global variable " + dataName + " not found")
		return
	}

	//【5】解析 Lua 表格数据
	luaTable := tableItem.(*lua.LTable)

	//【6】转换
	mapSlice = []map[string]string{}
	luaTable.ForEach(func(_, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			row := extractRowString(tbl, doCamel)
			mapSlice = append(mapSlice, row)
		}
	})

	return
}

// extractRowString 提取一行数据为map[string]string
func extractRowString(tbl *lua.LTable, doCamel bool) map[string]string {
	row := make(map[string]string)
	tbl.ForEach(func(key, value lua.LValue) {
		var columnName = key.String()
		if doCamel {
			columnName = camelToSnake(key.String())
		}
		row[columnName] = value.String()
	})
	return row
}
