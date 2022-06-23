package imgHelper

import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/wiidz/goutil/helpers/osHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"image"
	"io/ioutil"
	"os"
)

// OpenImageFile 打开图像文件
//func OpenImageFile(localUri string) (image.Image, error) {
//	var m image.Image
//	ff, _ := ioutil.ReadFile(localUri) //读取文件 要先下载
//	bbb := bytes.NewBuffer(ff)
//
//	datatype, err := imgtype.Get(localUri)
//
//	if err != nil {
//		fmt.Println(err)
//		return m, err
//	}
//	fmt.Println("【datatype】", datatype)
//
//	switch datatype {
//	case "image/jpeg":
//		m, err = jpeg.Decode(bbb)
//	case "image/png":
//		m, err = png.Decode(bbb)
//	case "image/gif":
//		m, err = gif.Decode(bbb)
//	default:
//		fmt.Println("不支持的格式", reflect.TypeOf(datatype).String())
//	}
//	return m, nil
//}

func OpenImageFile(localUri string) (image.Image, error) {
	return gg.LoadImage(localUri)
}

// Buff2Image 字节转图片
func Buff2Image(bytes []byte, filePath string) (err error) {
	//filePath = "/tmp/test.jpg"
	err = ioutil.WriteFile(filePath, bytes, 0666)
	return
}

// MergeLocalImg 拼合本地图片
func MergeLocalImg(bgImgFilePath, coverImgFilePath string, position *Position, coverSize *Size, newFilePath string) (err error) {

	bgImg, err := OpenImageFile(bgImgFilePath)
	if err != nil {
		return
	}

	coverImg, err := OpenImageFile(coverImgFilePath)
	if err != nil {
		return
	}

	if coverSize != nil {
		coverImg = imaging.Resize(coverImg, coverSize.Width, coverSize.Width, imaging.Lanczos)
	}

	context := gg.NewContextForImage(bgImg)
	context.DrawImage(coverImg, position.X, position.Y)
	err = context.SavePNG(newFilePath)

	return
}

// MergeNetworkImg 拼合网络图片
// 这个函数会将图片先下载到项目根目录下的/temp，记得给权限
func MergeNetworkImg(bgImgURL, coverImgURL string, position *Position, coverSize *Size, newFilePath string) (err error) {

	//【1】下载文件到本地
	dir, _ := os.Getwd()
	bgImgFilePath := dir + "/temp/" + strHelper.GetRandomString(8)
	coverImgFilePath := dir + "/temp/" + strHelper.GetRandomString(8)
	err = osHelper.DownloadFile(bgImgURL, bgImgFilePath, nil)
	if err != nil {
		return
	}
	defer osHelper.Delete(bgImgFilePath)

	err = osHelper.DownloadFile(coverImgURL, coverImgFilePath, nil)
	if err != nil {
		return
	}
	defer osHelper.Delete(coverImgFilePath)

	err = MergeLocalImg(bgImgFilePath, coverImgFilePath, position, coverSize, newFilePath)
	return
}
