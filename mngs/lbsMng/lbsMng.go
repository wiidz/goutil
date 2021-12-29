package lbsMng

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const ReGeoURL = "https://regeo.market.alicloudapi.com/v3/geocode/regeo"
const GeoURL = "https://geo.market.alicloudapi.com/v3/geocode/geo"
const DriveRouteURL = "http://direction.market.alicloudapi.com/v3/direction/driving"
const WalkRouteURL = "http://direction.market.alicloudapi.com/v3/direction/walking"
const BusRouteURL = "http://direction.market.alicloudapi.com/v3/direction/transit/integrated"

// 以下方法均为高德地图在阿里云市场中提供的服务

type LbsMng struct {
	Config *configStruct.AliApiConfig
}

// GetLbsMng : 返回地理位置管理器
func GetLbsMng(config *configStruct.AliApiConfig) *LbsMng {
	return &LbsMng{
		Config: config,
	}
}

// ReGeo : 逆地理编码(将经纬度转换为详细结构化的地址，且返回附近周边的POI、AOI信息)
func (mng *LbsMng) ReGeo(longitude, latitude string) (data *ReGeoData,err error) {

	tempStr, _, _, err := networkHelper.RequestRaw(networkStruct.Get, ReGeoURL, map[string]interface{}{
		"location": longitude + "," + latitude,
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	})

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
func (mng *LbsMng) Geo(address string) (*ReGeoData, error) {

	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, GeoURL, map[string]interface{}{
		"address": address,
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	}, &ReGeoRes{})

	if err != nil {
		return nil, err
	}

	return resStr.(*ReGeoRes).ReGeoCode, nil
}

// GetDriveRoute 驾车路径规划
// Docs: https://market.aliyun.com/products/56928004/cmapi020537.html?spm=5176.2020520132.101.1.4ed572180w4m2J#sku=yuncode1453700000
func (mng *LbsMng) GetDriveRoute(originLongitude, originLatitude, targetLongitude, targetLatitude string) (*RouteRes, error) {
	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, DriveRouteURL, map[string]interface{}{
		"destination": targetLongitude + "," + targetLatitude,
		"origin":      originLongitude + "," + originLatitude,
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	}, &RouteRes{})

	if err != nil {
		return nil, err
	}

	return resStr.(*RouteRes), nil
}
