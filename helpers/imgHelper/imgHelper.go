package imgHelper

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	imgtype "github.com/shamsher31/goimgtype"
	"github.com/wiidz/goutil/helpers/osHelper"
	"github.com/wiidz/goutil/helpers/strHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"reflect"
)

func OpenImageFile(localUri string) (image.Image, error) {
	return gg.LoadImage(localUri)
}

// Buff2Image 字节转图片
func Buff2Image(bytes []byte, filePath string) (err error) {
	err = ioutil.WriteFile(filePath, bytes, 0666)
	return
}

// LocalImg2Buff 本地图片转字节
func LocalImg2Buff(filePath string) (sendS3 []byte, err error) {

	//【1】打开image
	var img image.Image
	img, err = OpenImageFile(filePath)
	if err != nil {
		return
	}

	//【2】获取图片格式
	datatype, err := imgtype.Get(filePath)
	if err != nil {
		return
	}
	fmt.Println("【datatype】", datatype)

	//【3】转换成buff
	buf := new(bytes.Buffer)

	switch datatype {
	case "image/jpeg":
		err = jpeg.Encode(buf, img, nil)
	case "image/png":
		err = png.Encode(buf, img)
	case "image/gif":
		err = gif.Encode(buf, img, nil)
	default:
		fmt.Println()
		err = errors.New("不支持的格式" + reflect.TypeOf(datatype).String())
	}

	sendS3 = buf.Bytes()
	return
}

// MergeLocalImg 拼合本地图片
func MergeLocalImg(bgImgFilePath string, newFilePath string, coverImgSlice ...CoverImgInterface) (err error) {

	//【1】打开背景图
	var bgImg image.Image
	bgImg, err = OpenImageFile(bgImgFilePath)
	if err != nil {
		return
	}

	//【2】循环插入图片
	context := gg.NewContextForImage(bgImg)
	for k := range coverImgSlice {

		//【2-1】提取所需数据
		var size = coverImgSlice[k].GetSize()
		var position = coverImgSlice[k].GetPosition()
		var localFilePath = coverImgSlice[k].GetLocalFilePath()

		//【2-2】打开cover图片
		var temp image.Image
		temp, err = OpenImageFile(localFilePath)
		if err != nil {
			return
		}

		//【2-3】判断是否需要缩放
		if size != nil {
			temp = imaging.Resize(temp, int(size.Width), int(size.Height), imaging.Lanczos)
		}

		//【2-4】插入图片
		context.DrawImage(temp, int(position.X), int(position.Y))
	}

	//【3】输出新图片
	err = context.SavePNG(newFilePath)
	return
}

// MergeNetworkImg 拼合网络图片
// 这个函数会将图片先下载到项目根目录下的/temp，记得给权限
func MergeNetworkImg(bgImgURL string, newFilePath string, coverImgSlice []*NetworkCoverImg) (err error) {

	//【1】下载文件到本地
	dir, _ := os.Getwd()
	bgImgFilePath := dir + "/temp/" + strHelper.GetRandomString(8)

	err = osHelper.DownloadFile(bgImgURL, bgImgFilePath, nil)
	if err != nil {
		return
	}
	defer osHelper.Delete(bgImgFilePath)

	localFilePaths := []string{bgImgFilePath}

	interfaceSlice := []CoverImgInterface{}

	for k := range coverImgSlice {
		networkURL := coverImgSlice[k].NetworkURL
		localFilePath := dir + "/temp/" + strHelper.GetRandomString(8)
		err = osHelper.DownloadFile(networkURL, localFilePath, nil)
		if err != nil {
			return
		}

		coverImgSlice[k].LocalFilePath = localFilePath
		interfaceSlice = append(interfaceSlice, coverImgSlice[k])

		localFilePaths = append(localFilePaths, networkURL)
	}

	err = MergeLocalImg(bgImgFilePath, newFilePath, interfaceSlice...)

	go func() {
		for _, v := range localFilePaths {
			osHelper.Delete(v)
		}
	}()
	return
}

// DownloadNetworkImg 根据网络图片地址，转化为本地文件及路径
func DownloadNetworkImg(bgImgURL string) (bgImgFilePath string, err error) {
	//【1】下载文件到本地
	dir, _ := os.Getwd()
	bgImgFilePath = dir + "/temp/" + strHelper.GetRandomString(8)

	err = osHelper.DownloadFile(bgImgURL, bgImgFilePath, nil) // 注意这里是没有后缀名的
	return
}

// DownloadNetworkImgToDir 根据网络图片地址，转化为本地文件及路径
func DownloadNetworkImgToDir(bgImgURL string, localDirPath string) (bgImgFilePath string, err error) {
	//【1】下载文件到本地
	bgImgFilePath = localDirPath + "/temp/" + strHelper.GetRandomString(8)

	err = osHelper.DownloadFile(bgImgURL, bgImgFilePath, nil) // 注意这里是没有后缀名的
	return
}

// GetSizeFromStr 从字符串中转换成size
// 200,400 这种格式
func GetSizeFromStr(str string) (size *Size) {
	temp := typeHelper.ExplodeFloat64(str, ",")
	if len(temp) == 2 {
		size = &Size{
			Width:  temp[0],
			Height: temp[1],
		}
	}
	return
}

// GetPositionFromStr 从字符串中转换成point
// 0,120,234 这种格式
func GetPositionFromStr(str string) (imgNo int, position *Position) {

	temp := typeHelper.ExplodeFloat64(str, ",")
	if len(temp) == 3 {
		imgNo = int(temp[0])
		position = &Position{
			X: temp[1],
			Y: temp[2],
		}
	}
	return

}
