package lbsMng

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
)

const Domain = "https://regeo.market.alicloudapi.com/v3"

// 以下方法均为高德地图在阿里云市场中提供的服务

type LbsMng struct {
	Config *configStruct.AliApiConfig
}

// GetLbsMng : 返回地理位置管理器
func GetLbsMng(config *configStruct.AliApiConfig)*LbsMng{
	return &LbsMng{
		Config:config,
	}
}


// ReGeo : 逆地理编码(将经纬度转换为详细结构化的地址，且返回附近周边的POI、AOI信息)
func (mng *LbsMng)ReGeo(longitude,latitude string)(*ReGeoData,error){

	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkHelper.Get, Domain+"/geocode/regeo", map[string]interface{}{
		"location": longitude + "," + latitude,
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	},&ReGeoRes{})

	if err != nil {
		return nil,err
	}

	return resStr.(*ReGeoRes).ReGeoCode,nil
}


// Geo : 地理编码(将详细的结构化地址转换为高德经纬度坐标。且支持对地标性名胜景区、建筑物名称解析为高德经纬度坐标)
// Tips：举例，北京市朝阳区阜通东大街6号转换后经纬度：116.480881,39.989410 地
