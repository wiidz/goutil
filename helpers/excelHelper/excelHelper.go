package excelHelper

import (
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/xuri/excelize/v2"
	"math"
)

// GetSimpleCellStyle 获取简单单元格样式
func GetSimpleCellStyle(f *excelize.File, styleObj *SimpleCellStyle) (cellStyle int, err error) {

	//【1】填充默认值
	if styleObj.FontSize == 0 {
		styleObj.FontSize = 14
	}
	if styleObj.BgColor == "" {
		styleObj.BgColor = "#E4E7ED"
	}
	if styleObj.BorderColor == "" {
		styleObj.BorderColor = "#303133"
	}
	if styleObj.FontColor == "" {
		styleObj.FontColor = "#303133"
	}

	//【2】构建数据
	style := &excelize.Style{
		Alignment: &excelize.Alignment{
			Vertical:   "center",
			Horizontal: "center",
		},
		Font: &excelize.Font{
			Bold:  styleObj.FontBold,
			Size:  styleObj.FontSize,
			Color: styleObj.FontColor,
			//Strike: false,
		},
	}

	if styleObj.BorderFill {
		style.Border = []excelize.Border{
			{Type: "left", Color: styleObj.BorderColor, Style: 1}, // 这里的color没有#号
			{Type: "top", Color: styleObj.BorderColor, Style: 1},
			{Type: "right", Color: styleObj.BorderColor, Style: 1},
			{Type: "bottom", Color: styleObj.BorderColor, Style: 1},
		}
	}

	if styleObj.IsMoneyField {
		style.NumFmt = 193 //NumFmt
	}

	if styleObj.BgColorFill {
		style.Fill = excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{styleObj.BgColor}, //  这里的color有#号
		}
	}

	cellStyle, err = f.NewStyle(style)
	return
}

// GetLetter 获取列名（字母）
func GetLetter(index int) string {
	var letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	if index < 26 {
		return letters[index]
	} else if index < 26*26 {
		part := math.Floor(float64(index) / float64(26))
		return letters[int(part)-1] + letters[index-26*int(part)]
	} else {
		return ""
	}
}

// SetSingleCell 设置一个单元格的值
func SetSingleCell(f *excelize.File, sheetName string, rowNo int, columnNum int,value string,cellStyle int) {
	columnLetter := GetLetter(columnNum)
	_ = f.SetCellValue(sheetName, columnLetter+typeHelper.Int2Str(rowNo), value)
	_ = f.SetCellStyle(sheetName, columnLetter+typeHelper.Int2Str(rowNo), columnLetter+typeHelper.Int2Str(rowNo), cellStyle)
}

// SetMultiCell 设置一个占多个单元格的值（例如标题）
func SetMultiCell(f *excelize.File, sheetName string, rowNo int, fromColumnNum,endColumnNum int,value string,cellStyle int) {

	//【1】大标题

	startLetter := GetLetter(fromColumnNum)
	endLetter := GetLetter(endColumnNum)

	_ = f.MergeCell(sheetName, startLetter+typeHelper.Int2Str(rowNo), endLetter+typeHelper.Int2Str(rowNo))
	_ = f.SetCellValue(sheetName, startLetter+typeHelper.Int2Str(rowNo), value)

	_ = f.SetCellStyle(sheetName, startLetter+typeHelper.Int2Str(rowNo), startLetter+typeHelper.Int2Str(rowNo), cellStyle)
	_ = f.SetRowHeight(sheetName, rowNo, 30)

}

// SetTableTitle 设置表头的列名（ID、姓名、手机...）
func SetTableTitle(f *excelize.File, sheetName string, rowNo int, slice []HeaderSlice, headerStyle int, rowHeight float64)  {

	for _, v := range slice {
		_ = f.SetColWidth(sheetName, v.ColumnLetter, v.ColumnLetter, v.Width)
		_ = f.SetCellValue(sheetName, v.ColumnLetter+typeHelper.Int2Str(rowNo), v.Label)
	}

	_ = f.SetCellStyle(sheetName, "A"+typeHelper.Int2Str(rowNo), GetLetter(len(slice)-1)+typeHelper.Int2Str(rowNo), headerStyle)
	_ = f.SetRowHeight(sheetName, rowNo, rowHeight)

	return
}