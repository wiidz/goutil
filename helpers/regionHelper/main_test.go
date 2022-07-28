package regionSpiderMine

import (
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	//log.Print("开始采集")
	//defer spiderUtil.CostTime()()
	//spider := spider.NewSpider()
	//spider.Run()
	//fmt.Println("采集完成...")
	fileContent, err := GetJSFile()
	log.Println("err", err)
	ParsedJsFile(fileContent)
}
