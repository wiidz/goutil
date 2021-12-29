package typeHelper

// Map2InterfaceSlice 将map转换成interface slice
func Map2InterfaceSlice(data map[string]interface{}) []interface{} {
	iSlice := make([]interface{}, 0)
	for _, v := range data {
		iSlice = append(iSlice, v)
	}
	return iSlice
}

// StrMapToInterface 把map[string]string 转换成 map[string]interface{}
func StrMapToInterface(data map[string]string) map[string]interface{} {
	res := map[string]interface{}{}
	for k := range data {
		res[k] = data[k]
	}
	return res
}

// InterfaceMapToStr 把map[string]interface{} 转换成 map[string]string
func InterfaceMapToStr(data map[string]interface{}) map[string]string {
	res := map[string]string{}
	for k := range data {
		res[k] = ForceString(data[k])
	}
	return res
}
