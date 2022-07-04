package pdfHelper

import (
	"errors"
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/jung-kurt/gofpdf"
	"github.com/wiidz/goutil/helpers/imgHelper"
	"github.com/wiidz/goutil/helpers/mathHelper"
	"image"
	"image/jpeg"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

const (
	FontName             = "MyFont"
	PortraitValidWidth   = 190.0
	PortraitValidHeight  = 277.0
	LandscapeValidWidth  = 277.0
	LandscapeValidHeight = 190.0
	Margin               = 10

	HalfPortraitValidWidth = 95.0

	SignSpaceRowAmount      = 4    // 签字区域的空白行数
	SignerInfoRowAmount     = 7    // 签名人员信息高度（公司7 个人4）
	BlankRowHeight          = 8    // 签字区域单行高度
	SignFormPartyLineHeight = 10.0 // 签字区域 甲方乙方的行高
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

	helper.addFonts() // 添加预设字体
	helper.AddHeader()
	helper.AddFooter()
	helper.PDF.AddPage() // 添加一页
	// 这里不要忘记了，如果没有addPage，也能输出pdf，但是这个pdf的数据头不一样，就会导致fitz认不到格式

	return
}

// addFonts 添加字体
func (helper *PDFHelper) addFonts() {
	helper.PDF.AddUTF8Font(FontName, FontLight, helper.FontOption.LightTTFURL)
	helper.PDF.AddUTF8Font(FontName, FontRegular, helper.FontOption.RegularTTFURL)
	helper.PDF.AddUTF8Font(FontName, FontBold, helper.FontOption.BoldTTFURL)
	helper.PDF.AddUTF8Font(FontName, FontHeavy, helper.FontOption.HeavyTTFURL)
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
		helper.PDF.SetTextColor(144, 147, 153)             //设置字体
		helper.PDF.SetFont(FontName, gofpdf.AlignLeft, 12) //设置字体
		helper.PDF.CellFormat(helper.getValidWidth(), Margin, helper.HeaderOption.RightText, "", 1, gofpdf.AlignRight, false, 0, "")
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
		helper.PDF.CellFormat(height[0], Margin, helper.FooterOption.LeftText, "", 0, gofpdf.AlignLeft, false, 0, "")

		//【2】页码（中）
		helper.PDF.SetTextColor(48, 49, 51)
		helper.PDF.SetFont(FontName, FontBold, 12) //设置字体

		pageNow := helper.PDF.PageNo()
		//pageTotal := pdf.PageCount()
		//log.Println(pageNow, pageTotal)
		//pdf.CellFormat(height[1], 10, strconv.Itoa(pageNow)+" / "+strconv.Itoa(pageTotal), "", 0, gofpdf.AlignCenter, false, 0, "")
		helper.PDF.CellFormat(height[1], Margin, "第"+strconv.Itoa(pageNow)+"页", "", 0, gofpdf.AlignCenter, false, 0, "")
		//pageTotal = 99
		//【3】编号（右）
		helper.PDF.SetTextColor(96, 98, 102)
		helper.PDF.SetFont(FontName, gofpdf.AlignLeft, 10) //设置字体
		helper.PDF.CellFormat(height[2], Margin, helper.FooterOption.RightText, "", 0, gofpdf.AlignCenter, false, 0, "")
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
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 1, gofpdf.AlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 0, gofpdf.AlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 0, gofpdf.AlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 0, gofpdf.AlignCenter, false, 0, "")
		} else {
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 1, gofpdf.AlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 0, gofpdf.AlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextCn, "", 0, gofpdf.AlignCenter, false, 0, "")
			helper.PDF.CellFormat(100, markLineHt, helper.WaterMarkOption.TextEn, "", 0, gofpdf.AlignCenter, false, 0, "")
		}
	}

	helper.PDF.TransformEnd()
}

// MainTitle 设置主标题
func (helper *PDFHelper) MainTitle(text string) {

	helper.PDF.SetXY(Margin, 15)
	helper.PDF.SetFont(FontName, FontHeavy, 24) // 设置字体

	totalWidth := PortraitValidWidth
	if !helper.isPortraitHeader() {
		totalWidth = LandscapeValidWidth
	}

	helper.PDF.CellFormat(totalWidth, 24, text, "0", 2, gofpdf.AlignCenter, false, 0, "")
}

// FirstTitle 一级标题
func (helper *PDFHelper) FirstTitle(text string) {
	//helper.PDF.SetFont(FontName, "B", 14)
	helper.PDF.SetFontStyle(FontBold)
	helper.PDF.SetFontSize(14)

	helper.PDF.SetTextColor(0, 0, 0)
	helper.PDF.MultiCell(190, 14, text, "", gofpdf.AlignLeft, false)
}

// SecondTitle 二级标题
func (helper *PDFHelper) SecondTitle(text string) {
	helper.PDF.SetFontStyle(FontBold)
	helper.PDF.SetFontSize(12)
	helper.PDF.CellFormat(helper.getValidWidth(), 16, text, "", 2, gofpdf.AlignCenter, false, 0, "")
}

// NormalContent 常规正文内容（前面会有两格缩进）
func (helper *PDFHelper) NormalContent(text string, opt ...*ContentStyle) {

	//【1】默认样式
	var doIndent = true
	var fontSize = float64(10)
	var textAlign = gofpdf.AlignLeft
	var fontWeight = FontRegular
	var color = &RGBColor{
		R: 48,
		G: 49,
		B: 51,
	}
	var lineHeight = fontSize * 1

	//【2】判断有无设置的样式
	if len(opt) != 0 {
		doIndent = opt[0].DoIntent
		if opt[0].FontSize != 0 {
			fontSize = opt[0].FontSize
		}
		if opt[0].TextAlign != "" {
			textAlign = opt[0].TextAlign
		}
		if opt[0].FontWeight != "" {
			fontWeight = string(opt[0].FontWeight)
		}
		if opt[0].Color != nil {
			color = opt[0].Color
		}
		if opt[0].LineHeight != 0 {
			lineHeight = opt[0].LineHeight
		}
	}

	//【3】处理缩进
	if doIndent {
		text = "        " + text
	}

	//【3】写入
	helper.PDF.SetFont(FontName, fontWeight, fontSize)
	helper.PDF.SetTextColor(color.R, color.G, color.B)
	helper.PDF.MultiCell(190, lineHeight, text, "", textAlign, false)
}

// SaveAsPDF 保存为pdf
// dir : 要以斜杠 / 结尾
// fileName : 不要后缀名
func (helper *PDFHelper) SaveAsPDF(dir, fileName string) (filePath string, err error) {
	filePath = dir + fileName + ".pdf"
	err = helper.PDF.OutputFileAndClose(filePath)
	return
}

// SaveAsImgs 以图片格式保存
// dir : 要以斜杠 / 结尾
// fileName : 不要后缀名
func (helper *PDFHelper) SaveAsImgs(dir, fileName string) (imgFileNames []string, err error) {

	//【1】先导出为pdf
	pdfFilePath := dir + fileName + ".pdf"
	err = helper.PDF.OutputFileAndClose(pdfFilePath)
	if err != nil {
		return
	}
	defer os.Remove(pdfFilePath) // 完成后删除pdf文件

	//【2】打开pdf文件
	doc, err := fitz.New(pdfFilePath)
	if err != nil {
		return
	}

	//【3】循环将每页pdf转换成图片
	imgFileNames = []string{}
	for n := 0; n < doc.NumPage(); n++ {

		//【3-1】
		var img image.Image
		img, err = doc.Image(n)
		if err != nil {
			return
		}

		//【3-2】创建文件
		var file *os.File
		imgFileName := fmt.Sprintf(fileName+"-%02d.jpg", n)
		file, err = os.Create(filepath.Join(dir, imgFileName))
		if err != nil {
			return
		}

		//【3-3】写入图片信息
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			return
		}

		//【3-4】将文件名写入数组
		imgFileNames = append(imgFileNames, imgFileName)
		file.Close()
	}

	return
}

// AddSignForm 添加一个签字用的区域
func (helper *PDFHelper) AddSignForm(firstParty, secondParty SignerInterface, fillTime, fillIP bool) (err error) {

	//【1】获取两边的数据
	leftData := getPartyInfo(firstParty)
	log.Println("leftData", leftData)
	rightData := getPartyInfo(secondParty)
	log.Println("rightData", rightData)

	leftSignData := firstParty.GetSignData()
	rightSignData := secondParty.GetSignData()

	leftSignData.SignFormCellStyle = getTips(firstParty)
	rightSignData.SignFormCellStyle = getTips(secondParty)

	//【2】为完整的form创建空间
	helper.createSpaceForSignForm(fillTime, fillIP)

	//【2】填充甲乙双方信息
	helper.PDF.Ln(-1)
	helper.PDF.SetFont(FontName, FontBold, SignFormPartyLineHeight)
	helper.PDF.CellFormat(HalfPortraitValidWidth, SignFormPartyLineHeight, "甲方", "", 0, gofpdf.AlignCenter, false, 0, "")
	helper.PDF.CellFormat(HalfPortraitValidWidth, SignFormPartyLineHeight, "乙方", "", 1, gofpdf.AlignCenter, false, 0, "")

	helper.PDF.SetFont(FontName, FontRegular, 9)

	//【3】循环填充数据
	for k := range leftData {
		helper.PDF.CellFormat(HalfPortraitValidWidth, BlankRowHeight, leftData[k], "", 0, gofpdf.AlignLeft, false, 0, "")
		helper.PDF.CellFormat(HalfPortraitValidWidth, BlankRowHeight, rightData[k], "", 1, gofpdf.AlignLeft, false, 0, "")
	}

	//【4】下方签字盖章区域
	err = helper.drawSignArea(leftSignData, rightSignData, fillTime, fillIP)
	return
}

// getPartyInfo 获取甲方/乙方信息数据
func getPartyInfo(party SignerInterface) (fillData [SignerInfoRowAmount]string) {

	if party.GetKind() == Company {
		log.Println("company")
		temp, ok := party.(CompanySigner)
		log.Println("temp", temp)
		log.Println("ok", ok)
		ogBankName := temp.OgBankName
		if temp.OgBankNo != "" {
			ogBankName += "（行号" + temp.OgBankNo + "）"
		}

		fillData = [SignerInfoRowAmount]string{
			"单位名称：" + temp.OgName,
			"税        号：" + temp.OgLicenseNo,
			"单位地址：" + temp.OgAddress,
			"开户银行：" + ogBankName,
			"法定代表：" + temp.LawPersonName,
			"电        话：" + temp.OgTel,
			"传        真：" + temp.OgFax,
		}
	} else {
		log.Println("person")
		temp, _ := party.(PersonSigner)
		fillData = [SignerInfoRowAmount]string{
			"姓        名：" + temp.TrueName,
			"身份证号：" + temp.IDCardNo,
			"手        机：" + temp.Phone,
			"住        所：" + temp.Address,
			"",
			"",
			"",
		}
	}

	return
}

// getTips 获取签署区域的数据
func getTips(party SignerInterface) (fillData [4]*SignFormCellStyle) {
	fillData = [SignSpaceRowAmount]*SignFormCellStyle{{
		Content: "",
		Fill:    false,
	}, {
		Content: "",
		Fill:    false,
	}, {
		Content: "",
		Fill:    false,
	}, {
		Content: "",
		Fill:    false,
	}}
	if party.GetSignData().DoHint {

		fillData[0].Fill = true
		fillData[1].Fill = true
		fillData[2].Fill = true
		fillData[3].Fill = true

		fillData[0].Content = "请在此处红色区域内"
		fillData[1].Content = "签署本人姓名 \"" + party.GetHintName() + "\"，请勿冒名顶替"
		if party.GetKind() == Company {
			fillData[2].Content = "并加盖 本公司/单位公章"
		}
	}

	return
}

// getRandomImgCenter 根据区域和图形尺寸，获取一个随机的中心点
func getRandomImgCenter(area *RectArea, size *imgHelper.Size, overflowRate float64) (randomCenter *Point) {

	randomCenter = &Point{
		X: 0,
		Y: 0,
	}

	//【1】扩大一下区域
	if overflowRate != 0 {
		area.LeftTop.X *= 1 - overflowRate
		area.LeftTop.Y *= 1 - overflowRate
		area.RightTop.X *= 1 + overflowRate
		area.RightTop.Y *= 1 - overflowRate
		area.LeftBottom.X *= 1 - overflowRate
		area.LeftBottom.Y *= 1 + overflowRate
		area.RightBottom.X *= 1 + overflowRate
		area.RightBottom.Y *= 1 + overflowRate
	}

	//【2】寻找中心点
	distanceX := math.Abs(area.RightTop.X-area.LeftTop.X) - size.Width
	distanceY := math.Abs(area.LeftBottom.Y-area.LeftTop.Y) - size.Height

	log.Println("distanceX", distanceX)
	log.Println("distanceY", distanceY)
	//if distanceX < 0 {
	//	distanceX = 0
	//}
	//if distanceY < 0 {
	//	distanceY = 0
	//}

	randomCenter.X = area.LeftTop.X + distanceX*float64(mathHelper.GetRandomInt(0, 100))/100
	randomCenter.Y = area.LeftTop.Y + distanceY*float64(mathHelper.GetRandomInt(0, 100))/100

	return
}

// AddTableHead 添加一个表格头
func (helper *PDFHelper) AddTableHead(width float64, ln Ln, content string, opt ...*ContentStyle) {
	//【1】默认样式
	var fontSize = float64(10)
	var lineHeight = fontSize * 1
	var textAlign = gofpdf.AlignCenter
	var fontWeight = FontBold
	var color = &RGBColor{
		R: 48,
		G: 49,
		B: 51,
	}
	var bgColor *RGBColor

	//【2】判断有无设置的样式
	if len(opt) != 0 {
		if opt[0].FontSize != 0 {
			fontSize = opt[0].FontSize
		}
		if opt[0].TextAlign != "" {
			textAlign = opt[0].TextAlign
		}
		if opt[0].FontWeight != "" {
			fontWeight = string(opt[0].FontWeight)
		}
		if opt[0].Color != nil {
			color = opt[0].Color
		}
		if opt[0].BgColor != nil {
			bgColor = opt[0].BgColor
		}
		if opt[0].LineHeight != 0 {
			lineHeight = opt[0].LineHeight
		}
	}

	//【3】设置样式
	var fill bool
	if bgColor != nil {
		fill = true
		helper.PDF.SetFillColor(bgColor.R, bgColor.G, bgColor.B)
	}

	helper.PDF.SetFont(FontName, fontWeight, fontSize)
	helper.PDF.SetTextColor(color.R, color.G, color.B)

	helper.PDF.CellFormat(width, lineHeight, content, "LTRB", int(ln), textAlign, fill, 0, "")
}

// AddTableBody 添加一个表格体
func (helper *PDFHelper) AddTableBody(width float64, ln Ln, content string, opt ...*ContentStyle) {
	//【1】默认样式
	var fontSize = float64(10)
	var lineHeight = fontSize * 1
	var textAlign = gofpdf.AlignCenter
	var fontWeight = FontRegular
	var color = &RGBColor{
		R: 48,
		G: 49,
		B: 51,
	}
	var bgColor *RGBColor

	//【2】判断有无设置的样式
	if len(opt) != 0 {
		if opt[0].FontSize != 0 {
			fontSize = opt[0].FontSize
		}
		if opt[0].TextAlign != "" {
			textAlign = opt[0].TextAlign
		}
		if opt[0].FontWeight != "" {
			fontWeight = string(opt[0].FontWeight)
		}
		if opt[0].Color != nil {
			color = opt[0].Color
		}
		if opt[0].BgColor != nil {
			bgColor = opt[0].BgColor
		}
		if opt[0].LineHeight != 0 {
			lineHeight = opt[0].LineHeight
		}
	}

	//【3】设置样式
	var fill bool
	if bgColor != nil {
		fill = true
		helper.PDF.SetFillColor(bgColor.R, bgColor.G, bgColor.B)
	}

	helper.PDF.SetFont(FontName, fontWeight, fontSize)
	helper.PDF.SetTextColor(color.R, color.G, color.B)

	helper.PDF.CellFormat(width, lineHeight, content, "LTRB", int(ln), textAlign, fill, 0, "")
}

// AddTableHeadMulti 添加一个表格头
func (helper *PDFHelper) AddTableHeadMulti(width float64, startPoint *Point, content string, opt ...*ContentStyle) (thisLineEndPoint, nextLineStartPoint *Point) {

	thisLineEndPoint, nextLineStartPoint = &Point{}, &Point{}
	helper.PDF.SetXY(startPoint.X, startPoint.Y)
	//startPoint.X, startPoint.Y = helper.PDF.GetXY()

	//【1】默认样式
	var fontSize = float64(10)
	var lineHeight = fontSize * 1
	var textAlign = gofpdf.AlignCenter
	var fontWeight = FontBold
	var color = &RGBColor{
		R: 48,
		G: 49,
		B: 51,
	}
	var bgColor *RGBColor

	//【2】判断有无设置的样式
	if len(opt) != 0 {
		if opt[0].FontSize != 0 {
			fontSize = opt[0].FontSize
		}
		if opt[0].TextAlign != "" {
			textAlign = opt[0].TextAlign
		}
		if opt[0].FontWeight != "" {
			fontWeight = string(opt[0].FontWeight)
		}
		if opt[0].Color != nil {
			color = opt[0].Color
		}
		if opt[0].BgColor != nil {
			bgColor = opt[0].BgColor
		}
		if opt[0].LineHeight != 0 {
			lineHeight = opt[0].LineHeight
		}
	}

	//【3】设置样式
	var fill bool
	if bgColor != nil {
		fill = true
		helper.PDF.SetFillColor(bgColor.R, bgColor.G, bgColor.B)
	}

	helper.PDF.SetFont(FontName, fontWeight, fontSize)
	helper.PDF.SetTextColor(color.R, color.G, color.B)

	//helper.PDF.CellFormat(width, lineHeight, content, "LTRB", int(ln), textAlign, fill, 0, "")
	helper.PDF.MultiCell(width, lineHeight, content, "LTRB", textAlign, fill)

	thisLineEndPoint.X = startPoint.X + width
	thisLineEndPoint.Y = startPoint.Y
	nextLineStartPoint.X, nextLineStartPoint.Y = helper.PDF.GetXY()
	return
}

// AddTableBodyMulti 添加一个表格体
func (helper *PDFHelper) AddTableBodyMulti(width float64, startPoint *Point, content string, opt ...*ContentStyle) (thisLineEndPoint, nextLineStartPoint *Point) {

	thisLineEndPoint, nextLineStartPoint = &Point{}, &Point{}
	helper.PDF.SetXY(startPoint.X, startPoint.Y)

	//【1】默认样式
	var fontSize = float64(10)
	var lineHeight = fontSize * 1
	var textAlign = gofpdf.AlignCenter
	var fontWeight = FontRegular
	var color = &RGBColor{
		R: 48,
		G: 49,
		B: 51,
	}
	var bgColor *RGBColor

	//【2】判断有无设置的样式
	if len(opt) != 0 {
		if opt[0].FontSize != 0 {
			fontSize = opt[0].FontSize
		}
		if opt[0].TextAlign != "" {
			textAlign = opt[0].TextAlign
		}
		if opt[0].FontWeight != "" {
			fontWeight = string(opt[0].FontWeight)
		}
		if opt[0].Color != nil {
			color = opt[0].Color
		}
		if opt[0].BgColor != nil {
			bgColor = opt[0].BgColor
		}
		if opt[0].LineHeight != 0 {
			lineHeight = opt[0].LineHeight
		}
	}

	//【3】设置样式
	var fill bool
	if bgColor != nil {
		fill = true
		helper.PDF.SetFillColor(bgColor.R, bgColor.G, bgColor.B)
	}

	helper.PDF.SetFont(FontName, fontWeight, fontSize)
	helper.PDF.SetTextColor(color.R, color.G, color.B)

	//helper.PDF.CellFormat(width, lineHeight, content, "LTRB", int(ln), textAlign, fill, 0, "")
	helper.PDF.MultiCell(width, lineHeight, content, "LTRB", textAlign, fill)

	thisLineEndPoint.X = startPoint.X + width
	thisLineEndPoint.Y = startPoint.Y

	nextLineStartPoint.X, nextLineStartPoint.Y = helper.PDF.GetXY()
	return
}

// GetTotalHeight 获取多行文字的行高
func (helper *PDFHelper) GetTotalHeight(content string, width float64, weight FontWeight, fontSize, lineHeight float64) (totalHeight float64) {
	helper.PDF.SetFont(FontName, string(weight), fontSize)

	//lines := helper.PDF.SplitLines([]byte(content), width)
	lines := helper.PDF.SplitText(content, width)

	//log.Println("font:", string(weight), fontSize, lineHeight)
	//for _, v := range lines {
	//	log.Println("v", string(v))
	//}

	totalHeight = float64(len(lines)) * lineHeight
	return
}

// getSignArea 获取签名区域
func (helper *PDFHelper) getSignArea(fillTime, fillIP bool) (leftSignArea, rightSignArea RectArea) {

	//【1】初始化
	leftSignArea, rightSignArea = RectArea{}, RectArea{}

	//【2】获取当前坐标
	tempX, tempY := helper.PDF.GetXY()
	addRow := float64(1) // 1是 签字/盖章 那一行
	if fillTime {
		addRow++
	}
	if fillIP {
		addRow++
	}
	toY := tempY + (SignSpaceRowAmount+addRow)*BlankRowHeight

	//【3】构建初步的区域
	leftSignArea.LeftTop = Point{X: tempX, Y: tempY}
	leftSignArea.LeftBottom = Point{X: tempX, Y: toY} // 填充
	leftSignArea.RightTop = Point{X: tempX + HalfPortraitValidWidth, Y: tempY}
	leftSignArea.RightBottom = Point{X: tempX + HalfPortraitValidWidth, Y: toY} // 填充

	rightSignArea.LeftTop = Point{X: tempX + HalfPortraitValidWidth, Y: tempY}
	rightSignArea.LeftBottom = Point{X: tempX + HalfPortraitValidWidth, Y: toY}
	rightSignArea.RightTop = Point{X: tempX + HalfPortraitValidWidth + HalfPortraitValidWidth, Y: tempY}
	rightSignArea.RightBottom = Point{X: tempX + HalfPortraitValidWidth + HalfPortraitValidWidth, Y: toY}

	//log.Println("leftSignArea", leftSignArea)
	//log.Println("rightSignArea", rightSignArea)
	//helper.PDF.SetFillColor(200, 255, 255) // 设置填充颜色
	//helper.PDF.Rect(leftSignArea.LeftTop.X, leftSignArea.LeftTop.Y, leftSignArea.RightTop.X-leftSignArea.LeftTop.X, leftSignArea.LeftBottom.Y-leftSignArea.LeftTop.Y, "F")
	//helper.PDF.Rect(rightSignArea.LeftTop.X, rightSignArea.LeftTop.Y, rightSignArea.RightTop.X-rightSignArea.LeftTop.X, rightSignArea.LeftBottom.Y-rightSignArea.LeftTop.Y, "F")

	//【4】返回
	return
}

// drawSignArea 绘制签名区域
func (helper *PDFHelper) drawSignArea(leftSignData, rightSignData SignData, fillTime, fillIP bool) (err error) {

	//【1】获取签名区域
	leftSignArea, rightSignArea := helper.getSignArea(fillTime, fillIP) // 一定要在写签字/盖章提示之前调用

	//【2】第一行提示
	helper.PDF.SetFillColor(255, 235, 238) // 设置填充颜色
	helper.PDF.CellFormat(HalfPortraitValidWidth, 8, "签字/盖章：", "", 0, gofpdf.AlignLeft, leftSignData.DoHint, 0, "")
	helper.PDF.CellFormat(HalfPortraitValidWidth, 8, "签字/盖章：", "", 1, gofpdf.AlignLeft, rightSignData.DoHint, 0, "")

	//【3】中间内容（空白行和提示行）
	helper.PDF.SetTextColor(239, 154, 154) // 设置字体颜色
	for k := 0; k < SignSpaceRowAmount; k++ {
		helper.PDF.CellFormat(HalfPortraitValidWidth, BlankRowHeight, leftSignData.SignFormCellStyle[k].Content, "", 0, gofpdf.AlignCenter, leftSignData.SignFormCellStyle[k].Fill, 0, "")
		helper.PDF.CellFormat(HalfPortraitValidWidth, BlankRowHeight, rightSignData.SignFormCellStyle[k].Content, "", 1, gofpdf.AlignCenter, rightSignData.SignFormCellStyle[k].Fill, 0, "")
	}
	helper.PDF.SetTextColor(48, 49, 51) // 把字体颜色改回来

	//【4】填充签名和印章的图片
	err = helper.drawSignImg(&leftSignArea, leftSignData.StampImg, leftSignData.OverflowRate, leftSignData.AutoSign)
	if err != nil {
		return
	}
	err = helper.drawSignImg(&leftSignArea, leftSignData.NameImg, leftSignData.OverflowRate, leftSignData.AutoSign)
	if err != nil {
		return
	}
	err = helper.drawSignImg(&rightSignArea, rightSignData.StampImg, rightSignData.OverflowRate, rightSignData.AutoSign)
	if err != nil {
		return
	}
	err = helper.drawSignImg(&rightSignArea, rightSignData.NameImg, rightSignData.OverflowRate, rightSignData.AutoSign)
	if err != nil {
		return
	}

	//【5】填充IP
	if fillIP {
		helper.PDF.CellFormat(HalfPortraitValidWidth, 8, "IP："+leftSignData.IP, "", 0, gofpdf.AlignLeft, false, 0, "")
		helper.PDF.CellFormat(HalfPortraitValidWidth, 8, "IP："+rightSignData.IP, "", 1, gofpdf.AlignLeft, false, 0, "")
	}

	//【6】填充时间
	var timeStr = [2]string{"签署日期：", "签署日期："}
	if fillTime {
		timeStr[0] += leftSignData.Time
		timeStr[1] += rightSignData.Time
	}

	helper.PDF.CellFormat(HalfPortraitValidWidth, 8, timeStr[0], "", 0, gofpdf.AlignLeft, false, 0, "")
	helper.PDF.CellFormat(HalfPortraitValidWidth, 8, timeStr[1], "", 1, gofpdf.AlignLeft, false, 0, "")
	return
}

// drawSignImg 将签名图片按照位置放好
func (helper *PDFHelper) drawSignImg(signArea *RectArea, img *SignImg, overflowRate float64, autoSign bool) (err error) {
	if img == nil || img.URL == "" {
		return
	}

	//position := img.Position
	if img.Size == nil {
		return errors.New("图片url不为空，但size为空，无法完成签署")
	}
	if img.Position == nil {
		if autoSign == false {
			return
		}
		// 需要自动签署
		img.Position = getRandomImgCenter(signArea, img.Size, overflowRate)
		if err != nil {
			return
		}
	}
	log.Println("position", img.Position)
	log.Println("img.Size", img.Size)
	log.Println("overflowRate", overflowRate)
	log.Println("open")
	log.Println("img.URL", img.URL)
	helper.PDF.Image(img.URL, img.Position.X, img.Position.Y, img.Size.Width, img.Size.Height, false, "", 0, "")
	log.Println("opened")
	return
}

// createSpaceForSignForm 判断当前页面是否可以放得下整个签名表单（包含甲方乙方和其对应信息）
func (helper *PDFHelper) createSpaceForSignForm(fillTime, fillIP bool) {

	//【1】计算表单高度
	formHeight := float64(10) + float64(SignerInfoRowAmount)*BlankRowHeight + float64(SignSpaceRowAmount)*BlankRowHeight
	addRow := float64(1) // 1是 签字/盖章 那一行
	if fillTime {
		addRow++
	}
	if fillIP {
		addRow++
	}
	formHeight += (SignSpaceRowAmount + addRow) * BlankRowHeight

	//【2】获取当前y
	nowY := helper.PDF.GetY()

	//【3】判断
	log.Println("nowY", nowY)
	log.Println("formHeight", formHeight)
	log.Println("nowY+formHeight", nowY+formHeight)
	log.Println("PortraitValidHeight+Margin", PortraitValidHeight+Margin)
	if nowY+formHeight > PortraitValidHeight+Margin {
		helper.PDF.AddPage()
		helper.PDF.SetXY(Margin, Margin)
	}
}
