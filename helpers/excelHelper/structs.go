package excelHelper

import "github.com/xuri/excelize/v2"

// ExcelHelper Excel助手
type ExcelHelper struct {
	ExcelFile *excelize.File // Excel文件
	SheetName string // 左下角页面名
}

type HeaderSlice struct {
	Label        string
	Width        float64
}

// SimpleCellStyle 简单单元格格式
type SimpleCellStyle struct {
	BgColor     string // 背景色，默认 #E4E7ED
	BgColorFill bool   // 是否填充背景色，默认false

	FontSize  float64 // 字体大小，默认14
	FontBold  bool    // 字体加粗，默认false
	FontColor string  // 字体颜色，默认#303133

	BorderFill  bool   // 是否需要边框，默认false
	BorderColor string // 边框颜色，默认#303133

	IsMoneyField bool // 是否填充背景色，默认false
}
