package pdfHelper

import "github.com/jung-kurt/gofpdf"

type PDFHelper struct {
	PDF             *gofpdf.Fpdf
	FontOption      *FontOption      // 字体设置
	HeaderOption    *HeaderOption    // 页眉设置
	FooterOption    *FooterOption    // 页脚设置
	WaterMarkOption *WaterMarkOption // 水印设置
}

type HeaderSlice struct {
	Label string
	Width float64
}

type FontOption struct {
	LightTTFURL   string // 细体字体文件
	RegularTTFURL string // 常规体字体文件
	BoldTTFURL    string // 粗体字体文件
}
type FooterOption struct {
	LeftText  string // 左侧文字
	RightText string // 右侧文字
}

type HeaderOption struct {
	LeftImgURL string // 左侧Logo地址
	RightText  string // 右侧文字（一般是文件编号）
}

// WaterMarkOption 水印设置
type WaterMarkOption struct {
	TextCn   string    // 水印文字（中文）
	TextEn   string    // 水印文字（英文）
	FontSize float64   // 字体大小
	Color    *RGBColor // 颜色
}

// TextAlign 文字水平对齐方式
type TextAlign string

const TextAlignLeft = "L"
const TextAlignRight = "R"
const TextAlignCenter = "C"

type FontWeight string

const FontBold = "B"
const FontRegular = ""
const FontLight = "L"

// RGBColor 文字等颜色
type RGBColor struct {
	R int
	G int
	B int
}

type ContentStyle struct {
	DoIntent   bool       // 是否进行首行缩进两格
	TextAlign  TextAlign  // 水平对齐方式
	FontWeight FontWeight // 文字粗细
	FontSize   float64    // 字体大小
	Color      *RGBColor
}
