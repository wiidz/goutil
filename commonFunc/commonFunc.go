package goutil

//【1】计算偏移量
func GetOffset(pageNow, pageSize int) int {
	offset := 0
	if pageNow > 1 {
		offset = (pageNow - 1) * pageSize
	}
	return offset
}
