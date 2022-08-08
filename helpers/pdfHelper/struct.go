package pdfHelper

import (
	"github.com/jung-kurt/gofpdf"
	"github.com/wiidz/goutil/helpers/imgHelper"
)

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
	HeavyTTFURL   string // 超粗体字体文件
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

type FontWeight string

const FontLight = "L"
const FontRegular = ""
const FontBold = "B"
const FontHeavy = "H"

// RGBColor 文字等颜色
type RGBColor struct {
	R int
	G int
	B int
}

type ContentStyle struct {
	LineHeight float64    // 行高
	DoIntent   bool       // 是否进行首行缩进两格
	TextAlign  string     // 水平对齐方式  gofpdf.AlignCenter 等
	FontWeight FontWeight // 文字粗细
	FontSize   float64    // 字体大小
	Color      *RGBColor  // 文字颜色
	BgColor    *RGBColor  // 背景色
}

type SignKind int8

const Person SignKind = 1  // 个人
const Company SignKind = 2 // 单位

type SignerInterface interface {
	GetKind() SignKind
	GetSignData() SignData
	GetHintName() string // 提示签字的姓名
}

// CompanySigner 公司签署
type CompanySigner struct {
	OgName        string //【单位】单位名称
	OgLicenseNo   string //【单位】营业执照编号
	OgBankName    string //【单位】开户行名称
	OgBankNo      string //【单位】行号
	OgBankAccount string //【单位】银行账号
	OgTel         string //【单位】电话
	OgFax         string //【单位】传真
	OgAddress     string //【单位】地址
	LawPersonName string //【单位】法人姓名

	SignerName  string // 签署人真实姓名
	SignerPhone string // 签署人手机号

	SignData SignData // 签名信息
}

func (signer CompanySigner) GetKind() SignKind {
	return Company
}
func (signer CompanySigner) GetSignData() SignData {
	return signer.SignData
}

func (signer CompanySigner) GetHintName() string {
	return signer.SignerName
}

// PersonSigner 个人签署
type PersonSigner struct {
	TrueName string // 个人真实姓名
	Address  string // 地址
	Phone    string // 手机号
	IDCardNo string // 身份证号

	SignData SignData // 签名信息
}

type SignData struct {
	DoHint            bool // 是否提示签名区域高亮
	AutoSign          bool // 自动签名
	StampImg          *SignImg
	NameImg           *SignImg
	OverflowRate      float64 // 签名浮动区域（仅自动签名有效 0 - 1）
	PageNo            int     // 签署在合同的第几页
	Time              string  // 签署日期
	IP                string  // 签署IP
	SignFormCellStyle [SignSpaceRowAmount]*SignFormCellStyle

	FaceImg  *SignImg // 人脸识别照片（目前不做相关处理）
	SignerID uint64   // 签署人主体ID（目前不做相关处理）
}

func (signer PersonSigner) GetKind() SignKind {
	return Person
}

func (signer PersonSigner) GetSignData() SignData {
	return signer.SignData
}

func (signer PersonSigner) GetHintName() string {
	return signer.TrueName
}

type SignFormCellStyle struct {
	Content string
	Fill    bool
}

// ln indicates where the current position should go after the call. Possible
// values are 0 (to the right), 1 (to the beginning of the next line), and 2
// (below). Putting 1 is equivalent to putting 0 and calling Ln() just after.

type Ln int

const ToTheRight Ln = 0
const Wrap Ln = 1
const Below Ln = 2

type RectArea struct {
	LeftTop     imgHelper.Position
	RightTop    imgHelper.Position
	LeftBottom  imgHelper.Position
	RightBottom imgHelper.Position
}

type SignImg struct {
	URL      string              // 图片地址
	Position *imgHelper.Position // 中心点
	Size     *imgHelper.Size     // 尺寸
}
