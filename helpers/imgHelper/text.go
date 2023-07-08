package imgHelper

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
)

//字体相关
//type TextBrush struct {
//	FontType  *truetype.Font
//	FontSize  float64
//	FontColor *image.Uniform
//	TextWidth int
//}

type HorizonAlign int8

const Left HorizonAlign = 1
const Center HorizonAlign = 2
const Right HorizonAlign = 3

type FontStyle struct {
	DPI         float64
	Family      *truetype.Font
	Color       color.RGBA
	Width       int
	Size        float64
	SpaceAmount int          // 前面空几格
	Align       HorizonAlign // 水平对齐方式
}

// GetFontType 根据文件地址构建字库
func GetFontType(fontFilePath string) (fontType *truetype.Font, err error) {
	fontFile, err := ioutil.ReadFile(fontFilePath)
	if err != nil {
		return nil, err
	}
	fontType, err = truetype.Parse(fontFile)
	return
}

//func CalculateTextWidth(fnt font.Face, text string) fixed.Int26_6 {
//	bounds := font.MeasureString(fnt, text)
//	return bounds.Ceil()
//}

// DrawFontToImage 图片插入文字
func DrawFontToImage(rgba *image.RGBA, pt image.Point, content string, fontStyle FontStyle) (err error) {
	c := freetype.NewContext()

	for k := 0; k < fontStyle.SpaceAmount; k++ {
		content = "　" + content // 全角空格
	}

	// 创建绘图上下文
	c.SetDPI(fontStyle.DPI)
	c.SetFont(fontStyle.Family)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fontStyle.Size)
	c.SetClip(rgba.Bounds())
	c.SetSrc(image.NewUniform(fontStyle.Color))
	c.SetDst(rgba)

	// 获取文本宽度
	face := truetype.NewFace(fontStyle.Family, &truetype.Options{
		Size:    fontStyle.Size,
		DPI:     fontStyle.DPI,
		Hinting: font.HintingFull,
	})
	textWidth := font.MeasureString(face, content)

	// 计算对齐位置
	var textX fixed.Int26_6
	switch fontStyle.Align {
	case Left:
		textX = fixed.I(pt.X)
	case Center:
		textX = fixed.I(pt.X) - textWidth/2
	case Right:
		textX = fixed.I(pt.X) - textWidth
	default:
		textX = fixed.I(pt.X) // 默认居左
	}

	newPt := freetype.Pt(textX.Ceil(), pt.Y)
	//log.Println("fixed.I(pt.X)", fixed.I(pt.X))
	//log.Println("textWidth", textWidth)
	//log.Println("content", content)
	//log.Println("int(textX)", textX.Ceil())
	//log.Println("pt.Y", pt.Y)
	_, err = c.DrawString(content, newPt)
	//c.DrawString(content, freetype.Pt(pt.X, pt.Y))
	return
}

// Image2RGBA 图片转换成image.RGBA
func Image2RGBA(img image.Image) *image.RGBA {

	baseSrcBounds := img.Bounds().Max
	newWidth := baseSrcBounds.X
	newHeight := baseSrcBounds.Y
	des := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight)) // 底板
	//首先将一个图片信息存入jpg
	draw.Draw(des, des.Bounds(), img, img.Bounds().Min, draw.Over)
	return des
}

func SaveImage(targetPath string, m image.Image) error {
	fSave, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer fSave.Close()

	err = jpeg.Encode(fSave, m, nil)

	if err != nil {
		return err
	}

	return nil
}
