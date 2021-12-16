package lbsMng

import (
	"github.com/wiidz/goutil/mngs/lbsMng"
	"github.com/wiidz/goutil/structs/configStruct"
	"testing"
)

func TestReGeo(t *testing.T){
	t.Log("start")
	lbsM := lbsMng.GetLbsMng(&configStruct.AliApiConfig{
		AppCode:   "7191c7623e4043388c6a7d06f6589997",
		AppID:     "",
		AppSecret: "7191c7623e4043388c6a7d06f6589997",
	})
	//latitude := "30.319352"
	//longitude := "120.388651"
	latitude := "39.990056"
	longitude := "116.482005"
	res,err :=  lbsM.ReGeo(longitude,latitude)

	if err != nil {
		t.Log("err", err)
	}

	t.Log("res", res)
}



//func TestGeo(t *testing.T){
//	t.Log("start")
//	lbsM := lbsMng.GetLbsMng(&configStruct.AliApiConfig{
//		AppCode:   "7191c7623e4043388c6a7d06f6589997",
//		AppID:     "",
//		AppSecret: "7191c7623e4043388c6a7d06f6589997",
//	})
//	res,err :=  lbsM.Geo("北京市朝阳区阜通东大街")
//
//	if err != nil {
//		t.Log("err", err)
//	}
//
//	t.Log("res", res)
//}



//func TestExpress(t *testing.T){
//
//	resStr, _, _, err := networkHelper.RequestJson(networkHelper.Get, "https://wuliu.market.alicloudapi.com/kdi", map[string]interface{}{
//		"no": "75831566556911",
//	}, map[string]string{
//		"Authorization": "APPCODE 7191c7623e4043388c6a7d06f6589997",
//	})
//	t.Log("resStr", resStr)
//	t.Log("err", err)
//}

func TestDriveRoute(t *testing.T){
	t.Log("start")
	lbsM := lbsMng.GetLbsMng(&configStruct.AliApiConfig{
		AppCode:   "7191c7623e4043388c6a7d06f6589997",
		AppID:     "",
		AppSecret: "7191c7623e4043388c6a7d06f6589997",
	})
	//latitude := "30.319352"
	//longitude := "120.388651"
	originLatitude := "39.990056"
	originLongitude := "116.482005"

	targetLatitude := "39.994356"
	targetLongitude := "116.442005"

	res,err :=  lbsM.GetDriveRoute(originLongitude,originLatitude,targetLongitude,targetLatitude)

	if err != nil {
		t.Log("err", err)
	}

	t.Log("res", res.Route.Paths[0].Distance)
}