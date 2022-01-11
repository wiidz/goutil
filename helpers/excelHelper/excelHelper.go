package excelHelper

import (
	"github.com/wiidz/goutil/helpers/osHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/timeHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/xuri/excelize/v2"
	"log"
	"math"
	"time"
)

// NewExcelHelper 创建一个单页Excel助手
func NewExcelHelper(sheetName string) (helper *ExcelHelper){
	f := excelize.NewFile()
	_ = f.NewSheet(sheetName)

	helper = &ExcelHelper{
		ExcelFile: f,
		SheetName: sheetName,
	}

	return
}

// GetSimpleCellStyle 获取简单单元格样式
func (helper *ExcelHelper)GetSimpleCellStyle(styleObj *SimpleCellStyle) (cellStyle int, err error) {

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

	cellStyle, err = helper.ExcelFile.NewStyle(style)
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

/********* Style 样式相关  *************/

// SetRowHeight 设置行高
func (helper *ExcelHelper) SetRowHeight(rowNo int,rowHeight float64) error {
	return helper.ExcelFile.SetRowHeight(helper.SheetName, rowNo, rowHeight)
}

// SetCellStyle 直接设置单元格格式
func (helper *ExcelHelper) SetCellStyle(rowNo int,fromColumnNum,endColumnNum int,styleObj *SimpleCellStyle) (err error){

	//【1】获取样式
	var cellStyle int
	cellStyle,err = helper.GetSimpleCellStyle(styleObj)
	if err != nil {
		return
	}

	//【2】确定开始和结束的列
	startLetter := GetLetter(fromColumnNum)
	endLetter := GetLetter(endColumnNum)

	log.Println("start:",startLetter+typeHelper.Int2Str(rowNo))
	log.Println("endLetter:",endLetter+typeHelper.Int2Str(rowNo))

	//【3】设置
	err = helper.ExcelFile.SetCellStyle(helper.SheetName, startLetter+typeHelper.Int2Str(rowNo), endLetter+typeHelper.Int2Str(rowNo), cellStyle)
	return
}




/********* Style & Value 同时设置了样式和值  *************/

// SetTableTitle 设置表头的列名（ID、姓名、手机...）
func (helper *ExcelHelper) SetTableTitle(rowNo int, slice []HeaderSlice) (err error) {

	for _, v := range slice {
		err = helper.ExcelFile.SetColWidth(helper.SheetName, v.ColumnLetter, v.ColumnLetter, v.Width)
		if err != nil {
			return
		}

		err  = helper.ExcelFile.SetCellValue(helper.SheetName, v.ColumnLetter+typeHelper.Int2Str(rowNo), v.Label)
		if err != nil {
			return
		}
	}

	return
}




/********* Value 值相关  *************/

// SetCellValue 设置一个单元格的值
func (helper *ExcelHelper) SetCellValue(rowNo int, columnNum int,value string)(err error) {
	columnLetter := GetLetter(columnNum)
	err = helper.ExcelFile.SetCellValue(helper.SheetName, columnLetter+typeHelper.Int2Str(rowNo), value)
	return
}

// SetCellValues 批量设置数据（每个单元格占一行）
func (helper *ExcelHelper) SetCellValues(rowNo int,valueSlice []string) (err error) {

	var letter string
	for index,value := range valueSlice {
		letter = GetLetter(index)
		err = helper.ExcelFile.SetCellValue(helper.SheetName, letter+typeHelper.Int2Str(rowNo), value)
		if err != nil {
			return
		}
	}

	return
}

// SetMergedCellValue 设置一个占多个单元格的值（例如标题）
func (helper *ExcelHelper) SetMergedCellValue(rowNo int, fromColumnNum,endColumnNum int,value string) (err error) {

	//【1】确定开始和结束的列
	startLetter := GetLetter(fromColumnNum)
	endLetter := GetLetter(endColumnNum)

	//【2】合并单元格
	err = helper.ExcelFile.MergeCell(helper.SheetName, startLetter+typeHelper.Int2Str(rowNo), endLetter+typeHelper.Int2Str(rowNo))
	if err != nil {
		return
	}

	//【3】设置值
	err = helper.ExcelFile.SetCellValue(helper.SheetName, startLetter+typeHelper.Int2Str(rowNo), value)
	return
}

// SaveLocal 保存到本地
// dirPath：绝对路径，末尾不要接/，例子："/home/go_project/space-api/excel"
// 主意这个文件夹要777
func (helper *ExcelHelper) SaveLocal(dirPath string) (filePath,fileName,ymdStr string,err error){

	//【1】提取时间
	nowStr := timeHelper.MyJsonTime(time.Now()).GetPureNumberStr()
	ymdStr = nowStr[0:8] // 年月日的数字，方便按目录分割
	dirPath += "/" + ymdStr

	//【2】判断目录是否存在
	err = osHelper.CreateIfNotExist(dirPath,777)
	if err != nil {
		return
	}

	//【3】保存
	fileName = nowStr + "-" + strHelper.GetRandomString(4) + ".xlsx"
	filePath = dirPath + "/" + fileName
	err = helper.ExcelFile.SaveAs(filePath)

	return
}
