package pdfHelper

import (
	"github.com/jung-kurt/gofpdf"
	"strconv"
)

const (
	FontName             = "MyFont"
	PortraitValidWidth   = 190.0
	PortraitValidHeight  = 277.0
	LandscapeValidWidth  = 277.0
	LandscapeValidHeight = 190.0
	Margin               = 10
)

// NewPDFHelper 创建一个pdf助手
func NewPDFHelper(fontOption *FontOption, headerOption *HeaderOption, footerOption *FooterOption, waterMarkOption *WaterMarkOption) (helper *PDFHelper) {

	helper = &PDFHelper{
		PDF:             gofpdf.New("P", "mm", "A4", ""),
		FontOption:      fontOption,
		HeaderOption:    headerOption,
		FooterOption:    footerOption,
		WaterMarkOption: waterMarkOption,
	}

	helper.addFonts() //添加预设字体

	return
}

// addFonts 添加字体
func (helper *PDFHelper) addFonts() {
	helper.PDF.AddUTF8Font(FontName, FontRegular, helper.FontOption.RegularTTFURL)
	helper.PDF.AddUTF8Font(FontName, FontBold, helper.FontOption.BoldTTFURL)
	helper.PDF.AddUTF8Font(FontName, FontLight, helper.FontOption.LightTTFURL)
}

// getValidWidth 获取当前页有效宽度
func (helper *PDFHelper) getValidWidth() float64 {
	totalWidth := PortraitValidWidth
	if !helper.isPortraitHeader() {
		totalWidth = LandscapeValidWidth
	}
	return totalWidth
}

// isPortraitHeader 判断是否是竖直的
func (helper *PDFHelper) isPortraitHeader() bool {
	w, h := helper.PDF.GetPageSize()
	return w < h
}

// AddHeader 添加页眉（左边logo，右边文件编号）
func (helper *PDFHelper) AddHeader() {
	helper.PDF.SetHeaderFunc(func() {

		helper.PDF.SetXY(Margin, Margin)

		//【1】Logo
		helper.PDF.Image(helper.HeaderOption.LeftImgURL, Margin, Margin, 30, 0, false, "", 0, "") //插图

		//【2】合同编号
		helper.PDF.SetXY(Margin, Margin)
		helper.PDF.SetTextColor(144, 147, 153)          //设置字体
		helper.PDF.SetFont(FontName, TextAlignLeft, 12) //设置字体
		helper.PDF.CellFormat(helper.getValidWidth(), Margin, helper.HeaderOption.RightText, "", 1, TextAlignRight, false, 0, "")
		//【3】添加水印
		if helper.WaterMarkOption != nil {
			helper.AddWaterMark()
		}
		helper.PDF.SetXY(Margin, 30)
	})
}

// AddFooter 添加页脚（两端文字，中间页码）
func (helper *PDFHelper) AddFooter() {
	helper.PDF.SetFooterFunc(func() {
		var height []float64
		if helper.isPortraitHeader() {
			helper.PDF.SetXY(Margin, 277)
			height = []float64{60, 60, 70}
		} else {
			helper.PDF.SetXY(Margin, 190)
			height = []float64{90, 97, 90}
		}

		//【1】编号（左）
		helper.PDF.SetTextColor(96, 98, 102)
		helper.PDF.SetFont(FontName, FontLight, 10) //设置字体
		helper.PDF.CellFormat(height[0], Margin, helper.FooterOption.LeftText, "", 0, TextAlignLeft, false, 0, "")

		//【2】页码（中）
		helper.PDF.SetTextColor(48, 49, 51)
		helper.PDF.SetFont(FontName, FontBold, 12) //设置字体

		pageNow := helper.PDF.PageNo()
		//pageTotal := pdf.PageCount()
		//log.Println(pageNow, pageTotal)
		//pdf.CellFormat(height[1], 10, strconv.Itoa(pageNow)+" / "+strconv.Itoa(pageTotal), "", 0, TextAlignCenter, false, 0, "")
		helper.PDF.CellFormat(height[1], Margin, "第"+strconv.Itoa(pageNow)+"页", "", 0, TextAlignCenter, false, 0, "")
		//pageTotal = 99
		//【3】编号（右）
		helper.PDF.SetTextColor(96, 98, 102)
		helper.PDF.SetFont(FontName, TextAlignLeft, 10) //设置字体
		helper.PDF.CellFormat(height[2], Margin, helper.FooterOption.RightText, "", 0, TextAlignRight, false, 0, "")
	})
}

// AddWaterMark 添加水印（两种文字交替）
func (helper *PDFHelper) AddWaterMark() {

	markLineHt := helper.PDF.PointToUnitConvert(96)
	ctrX := LandscapeValidWidth / 2.0
	ctrY := LandscapeValidHeight / 2.0

	helper.PDF.SetTextColor(helper.WaterMarkOption.Color.R, helper.WaterMarkOption.Color.G, helper.WaterMarkOption.Color.B)
	helper.PDF.SetFont(FontName, "", helper.WaterMarkOption.FontSize)
	helper.PDF.SetXY(30, 0)
	helper.PDF.TransformBegin()
	helper.PDF.TransformRotate(15, ctrX, ctrY)

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 1, TextAlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 0, TextAlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 0, TextAlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 0, TextAlignCenter, false, 0, "")
		} else {
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 1, TextAlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 0, TextAlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 0, TextAlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 0, TextAlignCenter, false, 0, "")
		}
	}

	helper.PDF.TransformEnd()
}

// MainTitle 设置主标题
func (helper *PDFHelper) MainTitle(text string) {

	helper.PDF.SetXY(Margin, 15)
	helper.PDF.SetFont(FontName, FontBold, 24) // 设置字体

	totalWidth := PortraitValidWidth
	if !helper.isPortraitHeader() {
		totalWidth = LandscapeValidWidth
	}

	helper.PDF.CellFormat(totalWidth, 24, text, "0", 2, TextAlignCenter, false, 0, "")
}

// FirstTitle 一级标题
func (helper *PDFHelper) FirstTitle(text string) {
	helper.PDF.SetFont(FontName, "B", 14)
	helper.PDF.SetTextColor(0, 0, 0)
	helper.PDF.MultiCell(190, 14, text, "", TextAlignLeft, false)
}

// SecondTitle 二级标题
func (helper *PDFHelper) SecondTitle(text string) {
	helper.PDF.SetFont(FontName, "B", 16) // 设置字体
	helper.PDF.CellFormat(helper.getValidWidth(), 16, text, "1", 2, TextAlignCenter, false, 0, "")
}

// NormalContent 常规正文内容（前面会有两格缩进）
func (helper *PDFHelper) NormalContent(text string, doIndent bool, textAlign TextAlign) {

	//【1】首行是否需要缩进两格
	if doIndent {
		text = "        " + text
	}

	//【2】写入
	helper.PDF.SetFont(FontName, "", 10)
	helper.PDF.SetTextColor(48, 49, 51)
	helper.PDF.MultiCell(190, 8, text, "", string(textAlign), false)
}
