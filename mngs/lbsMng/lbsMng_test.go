package lbsMng

import (
	"github.com/wiidz/goutil/mngs/lbsMng"
	"github.com/wiidz/goutil/structs/configStruct"
	"testing"
)

func TestGetLbsMng(t *testing.T){
	t.Log("start")
	lbsM := lbsMng.GetLbsMng(&configStruct.AliApiConfig{
		AppCode:   "",
		AppID:     "",
		AppSecret: "",
	})
	latitude := "30.319352"
	longitude := "120.388651"
	res,err :=  lbsM.ReGeo(longitude,latitude)

	if err != nil {
		t.Log("err", err)
	}

	t.Log("res", res)
}
