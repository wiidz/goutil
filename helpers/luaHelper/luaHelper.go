package luaHelper

import (
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"os"
	"regexp"
	"strings"
)

// ReadMapSliceFromFile 从lua文件中读取数据并转换为[]map[string]interface{}
func ReadMapSliceFromFile(filePath string, dataName string, doCamel bool) (mapSlice []map[string]interface{}, err error) {

	//【1】读byte
	file, _ := os.Open(filePath)
	defer file.Close()
	byteData, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	//【2】创建 Lua 解释器
	L := lua.NewState()
	defer L.Close()

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
	mapSlice = []map[string]interface{}{}
	luaTable.ForEach(func(_, value lua.LValue) {
		if tbl, ok := value.(*lua.LTable); ok {
			row := extractRow(tbl, doCamel)
			mapSlice = append(mapSlice, row)
		}
	})

	return
}

// extractRow 提取一行数据为map[string]interface{}
func extractRow(tbl *lua.LTable, doCamel bool) map[string]interface{} {
	row := make(map[string]interface{})
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
