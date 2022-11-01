package imgHelper

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"os"
)

//字体相关
type TextBrush struct {
	FontType  *truetype.Font
	FontSize  float64
	FontColor *image.Uniform
	TextWidth int
}

func NewTextBrush(FontFilePath string, FontSize float64, FontColor *image.Uniform, textWidth int) (*TextBrush, error) {
	fontFile, err := ioutil.ReadFile(FontFilePath)
	if err != nil {
		return nil, err
	}
	fontType, err := truetype.Parse(fontFile)
	if err != nil {
		return nil, err
	}
	if textWidth <= 0 {
		textWidth = 20
	}
	return &TextBrush{FontType: fontType, FontSize: FontSize, FontColor: FontColor, TextWidth: textWidth}, nil
}

// DrawFontOnRGBA 图片插入文字
func (fb *TextBrush) DrawFontOnRGBA(rgba *image.RGBA, pt image.Point, content string) {
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fb.FontType)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fb.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fb.FontColor)
	c.DrawString(content, freetype.Pt(pt.X, pt.Y))

}

type FontStyle struct {
	Color       color.RGBA
	Width       int
	Size        float64
	SpaceAmount int // 前面空几格
}

// DrawFontOnRGBAWithStyle 图片插入文字
func (fb *TextBrush) DrawFontOnRGBAWithStyle(rgba *image.RGBA, pt image.Point, content string, fontStyle FontStyle) {

	fb.FontColor = image.NewUniform(fontStyle.Color)
	fb.FontSize = fontStyle.Size
	fb.TextWidth = fontStyle.Width
	for k := 0; k < fontStyle.SpaceAmount; k++ {
		content += " "
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(fb.FontType)
	c.SetHinting(font.HintingFull)
	c.SetFontSize(fb.FontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fb.FontColor)
	c.DrawString(content, freetype.Pt(pt.X, pt.Y))

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
