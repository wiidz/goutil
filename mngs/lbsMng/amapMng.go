package lbsMng

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
	"log"
)

const ReGeoURL = "https://restapi.amap.com/v3/geocode/regeo"
const GeoURL = "https://restapi.amap.com/v3/geocode/geo"
const DriveRouteURL = "https://restapi.amap.com/v5/direction/driving"
const WalkRouteURL = "http://direction.market.alicloudapi.com/v3/direction/walking"
const BusRouteURL = "http://direction.market.alicloudapi.com/v3/direction/transit/integrated"

// 以下方法均为高德地图在阿里云市场中提供的服务

type AmapMng struct {
	//Config *configStruct.AliApiConfig
	Config *configStruct.AmapConfig
}

// NewAmapMng : 返回地理位置管理器
func NewAmapMng(config *configStruct.AmapConfig) *AmapMng {
	return &AmapMng{
		Config: config,
	}
}

// ReGeo : 逆地理编码(将经纬度转换为详细结构化的地址，且返回附近周边的POI、AOI信息)
func (mng *AmapMng) ReGeo(longitude, latitude string) (data *ReGeoData, err error) {

	tempStr, _, _, err := networkHelper.RequestRaw(networkStruct.Get, ReGeoURL, map[string]interface{}{
		"key":      mng.Config.Key,
		"location": longitude + "," + latitude,
	}, nil)

	log.Println("【tempStr】", tempStr)

	temp := ReGeoRes{}
	err = typeHelper.JsonDecodeWithStruct(tempStr, &temp)
	if err != nil {
		return
	}

	data = temp.ReGeoCode
	return
}

// Geo : 地理编码(将详细的结构化地址转换为高德经纬度坐标。且支持对地标性名胜景区、建筑物名称解析为高德经纬度坐标)
// Tips：举例，北京市朝阳区阜通东大街6号转换后经纬度：116.480881,39.989410 地
func (mng *AmapMng) Geo(address string) (*ReGeoData, error) {

	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, GeoURL, map[string]interface{}{
		"key":     mng.Config.Key,
		"address": address,
	}, nil, &ReGeoRes{})

	if err != nil {
		return nil, err
	}

	return resStr.(*ReGeoRes).ReGeoCode, nil
}

// GetDriveRoute 驾车路径规划
// Docs: https://lbs.amap.com/api/webservice/guide/api/newroute
func (mng *AmapMng) GetDriveRoute(originLongitude, originLatitude, targetLongitude, targetLatitude string) (*RouteRes, error) {
	resStr, _, _, err := networkHelper.RequestJsonWithStructTest(networkStruct.Get, DriveRouteURL, map[string]interface{}{
		"key":         mng.Config.Key,
		"destination": targetLongitude + "," + targetLatitude,
		"origin":      originLongitude + "," + originLatitude,
		"show_fields": "cost,tmcs",
	}, nil, &RouteRes{})

	log.Println("resStr", resStr)

	if err != nil {
		return nil, err
	}

	return resStr.(*RouteRes), nil
}
