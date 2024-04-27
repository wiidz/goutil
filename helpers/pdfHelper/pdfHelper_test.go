package pdfHelper

import (
	"github.com/wiidz/goutil/helpers/loggerHelper"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func TestAdd(t *testing.T) {
	logH, _ := loggerHelper.NewLoggerHelper(&loggerHelper.Config{
		IsFullPath:      true,
		ShowFileAndLine: true,
		Json:            false,
		Level:           zapcore.DebugLevel,
	})
	fontOption := &FontOption{
		FontName:      "MyFont",
		LightTTFURL:   "./fonts/Alibaba-PuHuiTi-Light.ttf",
		RegularTTFURL: "./fonts/Alibaba-PuHuiTi-Regular.ttf",
		BoldTTFURL:    "./fonts/Alibaba-PuHuiTi-Medium.ttf",
		HeavyTTFURL:   "./fonts/Alibaba-PuHuiTi-Heavy.ttf",
	}
	pdfH := NewPDFHelper(fontOption, nil, nil, nil)
	pdfH.MainTitle("购销合同2")

	headerSlice := []HeaderSlice{
		{"序号", 40},
		{"品名", 100},
		{"规格型号", 52},
	} // 表头数据

	var dataSlice = [][]string{
		[]string{
			"test",
			"一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789",
			"23321",
		},
		[]string{
			"test",
			"一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789",
			"23321",
		},
		[]string{
			"test",
			"一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789一二三四五六七八九123456789",
			"23321",
		},
	}
	for k := range dataSlice {
		pdfH.AddTableBodyRow([]float64{headerSlice[0].Width, headerSlice[1].Width, headerSlice[2].Width}, dataSlice[k], &ContentStyle{
			FontSize:   10,
			LineHeight: 10 * 0.6,
		})
	}

	//【6】输出
	mainDir, _ := os.Getwd()
	imgNames, err := pdfH.SaveAsPDF(mainDir+"/", "test")
	logH.Info("imgNames ", imgNames)
	logH.Error("Error occurred", err)
	return
}
