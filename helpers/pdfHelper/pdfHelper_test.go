package pdfHelper

import (
	"go.uber.org/zap"
	"log"
	"os"
	"testing"
)

var fontOption = &FontOption{
	//LightTTFURL:   "/static/fonts/Alibaba-PuHuiTi-Light.ttf",
	//RegularTTFURL: "/static/fonts/Alibaba-PuHuiTi-Regular.ttf",
	//BoldTTFURL:    "/static/fonts/Alibaba-PuHuiTi-Medium.ttf",
	//HeavyTTFURL:   "/static/fonts/Alibaba-PuHuiTi-Heavy.ttf",
}
var footerOption = &FooterOption{
	LeftText:  "", // 写平台名称（心宿乐清）
	RightText: "ANTARES (YQ) INFO-TECH. CO.LTD",
}
var headerOption = &HeaderOption{
	//LeftImgURL: "./static/images/logo.png",
	RightText: "", // 写合同编号
}
var waterMarkOption = &WaterMarkOption{
	TextCn:   "联合电气",
	TextEn:   "21B.cn",
	FontSize: 16,
	Color: &RGBColor{
		R: 189,
		G: 205,
		B: 215,
	},
}

func TestAdd(t *testing.T) {

	logger, _ := zap.NewProduction()
	logger = logger.WithOptions(zap.AddCaller()) // 添加 caller 选项
	defer logger.Sync()

	mainDir, _ := os.Getwd()
	fontOption = &FontOption{
		LightTTFURL:   "./fonts/Alibaba-PuHuiTi-Light.ttf",
		RegularTTFURL: "./fonts/Alibaba-PuHuiTi-Regular.ttf",
		BoldTTFURL:    "./fonts/Alibaba-PuHuiTi-Medium.ttf",
		HeavyTTFURL:   "./fonts/Alibaba-PuHuiTi-Heavy.ttf",
	}

	log.Println("main")
	pdfH := NewPDFHelper(fontOption, headerOption, footerOption, waterMarkOption)
	pdfH.MainTitle("购销合同")

	headerSlice := []HeaderSlice{
		{"序号", 8},
		{"品名", 44},
		{"规格型号", 52},
	} // 表头数据

	pdfH.AddTableBody(headerSlice[0].Width, ToTheRight, "123")
	pdfH.AddTableBody(headerSlice[1].Width, ToTheRight, "456")
	pdfH.AddTableBody(headerSlice[2].Width, ToTheRight, "789")

	logger.Info("test")

	//【6】输出
	fileName := "contact-letter-test"

	imgNames, err := pdfH.SaveAsImgs(mainDir+"/temp", fileName)
	logger.Info("imgNames", zap.Reflect("imgNames", imgNames))
	logger.Error("Error occurred", zap.Error(err))

	return
}
