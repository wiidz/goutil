package imgHelper

import (
	"github.com/fogleman/gg"
	"image"
)

func OpenImageFile(localUri string) (image.Image, error) {
	return gg.LoadImage(localUri)
}

// OpenImageFileContext 获取
func OpenImageFileContext(localUri string) (*gg.Context, error) {
	temp, err := gg.LoadImage(localUri)
	if err != nil {
		return nil, err
	}

	dc := gg.NewContextForImage(temp)
	return dc, nil
}

// CropCircleCenter 将图片从中心裁切成最大圆
func CropCircleCenter(target image.Image) *gg.Context {

	dc := gg.NewContextForImage(target)

	// 获取图片尺寸
	width := float64(dc.Width())
	height := float64(dc.Height())

	// 创建新的绘图上下文，大小为正方形
	size := int(width)
	if width > height {
		size = int(height)
	}
	dc = gg.NewContext(size, size)

	// 绘制圆形路径
	radius := float64(0)
	if width > height {
		radius = height / 2
	} else {
		radius = width / 2
	}
	dc.DrawCircle(width/2, height/2, radius)
	dc.Clip()

	// 填充背景色
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// 绘制裁剪后的圆形图片
	dc.DrawImage(target, 0, 0)
	return dc
}

// CropCornerRadius 裁圆角
func CropCornerRadius(target *gg.Context, radius float64) *gg.Context {

	// 获取图片尺寸
	width := float64(target.Width())
	height := float64(target.Height())

	target.DrawRoundedRectangle(0, 0, width, height, radius)
	target.SetRGB(0, 0, 0)
	target.Fill()

	return target
}
