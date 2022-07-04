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
func MergeNetworkImg(bgImgURL string, newFilePath string, coverImgSlice ...*NetworkCoverImg) (err error) {

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
