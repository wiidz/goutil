package pdfHelper

import (
	"fmt"
	"github.com/gen2brain/go-fitz"
	"github.com/jung-kurt/gofpdf"
	"image"
	"image/jpeg"
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
		helper.PDF.SetTextColor(144, 147, 153)             //设置字体
		helper.PDF.SetFont(FontName, gofpdf.AlignLeft, 12) //设置字体
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
	helper.PDF.SetFont(FontName, FontBold, 24) // 设置字体

	totalWidth := PortraitValidWidth
	if !helper.isPortraitHeader() {
		totalWidth = LandscapeValidWidth
	}

	helper.PDF.CellFormat(totalWidth, 24, text, "0", 2, gofpdf.AlignCenter, false, 0, "")
}

// FirstTitle 一级标题
func (helper *PDFHelper) FirstTitle(text string) {
	helper.PDF.SetFont(FontName, "B", 14)
	helper.PDF.SetTextColor(0, 0, 0)
	helper.PDF.MultiCell(190, 14, text, "", gofpdf.AlignLeft, false)
}

// SecondTitle 二级标题
func (helper *PDFHelper) SecondTitle(text string) {
	helper.PDF.SetFont(FontName, "B", 16) // 设置字体
	helper.PDF.CellFormat(helper.getValidWidth(), 16, text, "1", 2, gofpdf.AlignCenter, false, 0, "")
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

	//【2】判断有无设置的样式
	if len(opt) != 0 {
		doIndent = opt[0].DoIntent
		if opt[0].FontSize != 0 {
			fontSize = opt[0].FontSize
		}
		if opt[0].TextAlign != "" {
			textAlign = string(opt[0].TextAlign)
		}
		if opt[0].FontWeight != "" {
			fontWeight = string(opt[0].FontWeight)
		}
		if opt[0].Color != nil {
			color = opt[0].Color
		}
	}

	//【3】处理缩进
	if doIndent {
		text = "        " + text
	}

	//【3】写入
	helper.PDF.SetFont(FontName, fontWeight, fontSize)
	helper.PDF.SetTextColor(color.R, color.G, color.B)
	helper.PDF.MultiCell(190, 8, text, "", textAlign, false)
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
func (helper *PDFHelper) AddSignForm(firstParty, secondParty SignerInterface, fillTime, fillIP bool) {

	//【1】获取两边的数据
	leftData := getPartyInfo(firstParty)
	rightData := getPartyInfo(secondParty)

	helper.PDF.Ln(-1)
	helper.PDF.SetFont(FontName, FontBold, 10)
	helper.PDF.CellFormat(95, 10, "甲方", "1", 0, gofpdf.AlignCenter, false, 0, "")
	helper.PDF.CellFormat(95, 10, "乙方", "1", 1, gofpdf.AlignCenter, false, 0, "")

	helper.PDF.SetFont(FontName, FontRegular, 9)

	//【3】循环填充数据
	for k := range leftData {
		helper.PDF.CellFormat(95, 8, leftData[k], "LR", 0, gofpdf.AlignLeft, false, 0, "")
		helper.PDF.CellFormat(95, 8, rightData[k], "LR", 1, gofpdf.AlignLeft, false, 0, "")
	}

	//【4】下方签字盖章区域
	helper.PDF.SetFillColor(255, 235, 238) // 设置填充颜色

	helper.PDF.CellFormat(95, 8, "签字/盖章：", "LTR", 0, gofpdf.AlignLeft, firstParty.GetDoHint(), 0, "")
	helper.PDF.CellFormat(95, 8, "签字/盖章：", "LTR", 1, gofpdf.AlignLeft, secondParty.GetDoHint(), 0, "")

	//【5】获取两边的数据
	leftStyle := getSignData(firstParty)
	rightStyle := getSignData(secondParty)

	helper.PDF.SetTextColor(239, 154, 154) // 设置字体颜色
	for k := range leftStyle {
		helper.PDF.CellFormat(95, 8, leftStyle[k].Content, "LR", 0, gofpdf.AlignCenter, leftStyle[k].Fill, 0, "")
		helper.PDF.CellFormat(95, 8, rightStyle[k].Content, "LR", 1, gofpdf.AlignCenter, rightStyle[k].Fill, 0, "")
	}
	helper.PDF.SetTextColor(48, 49, 51) // 把字体颜色改回来

	//【6】填充IP、时间
	if fillIP {
		helper.PDF.CellFormat(95, 8, "IP："+firstParty.GetIP(), "LR", 0, gofpdf.AlignLeft, false, 0, "")
		helper.PDF.CellFormat(95, 8, "IP："+secondParty.GetIP(), "LR", 1, gofpdf.AlignLeft, false, 0, "")
	}

	var timeStr = [2]string{"签署日期：", "签署日期："}
	if fillTime {
		timeStr[0] += firstParty.GetTime()
		timeStr[1] += secondParty.GetTime()
	}
	helper.PDF.CellFormat(95, 8, timeStr[0], "LBR", 0, gofpdf.AlignLeft, false, 0, "")
	helper.PDF.CellFormat(95, 8, timeStr[1], "LBR", 1, gofpdf.AlignLeft, false, 0, "")

}

// getPartyInfo 获取甲方/乙方信息数据
func getPartyInfo(party SignerInterface) (fillData [7]string) {

	if party.GetKind() == Company {
		temp, _ := party.(CompanySigner)

		ogBankName := temp.OgBankName
		if temp.OgBankNo != "" {
			ogBankName += "（行号" + temp.OgBankNo + "）"
		}

		fillData = [7]string{
			"单位名称：" + temp.OgName,
			"税        号：" + temp.OgLicenseNo,
			"单位地址：" + temp.OgAddress,
			"开户银行：" + ogBankName,
			"法定代表：" + temp.LawPersonName,
			"电        话：" + temp.OgTel,
			"传        真：" + temp.OgFax,
		}
	} else {
		temp, _ := party.(PersonSigner)
		fillData = [7]string{
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

// getSignData 获取签署区域的数据
func getSignData(party SignerInterface) (fillData [4]*SignFormCellStyle) {

	fillData = [4]*SignFormCellStyle{{
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
	if party.GetDoHint() {

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
