package ipMng

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const URL = "https://c2ba.api.huachen.cn/ip"

// 以下方法均为高德地图在阿里云市场中提供的服务

type IPMng struct {
	Config *configStruct.AliApiConfig
}

// NewIPMng : 返回地理位置管理器
func NewIPMng(config *configStruct.AliApiConfig) *IPMng {
	return &IPMng{
		Config: config,
	}
}

// GetRegionInfo 根据IP获取区域信息
func (mng *IPMng) GetRegionInfo(ip string) (data *RespData, err error) {

	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, URL, map[string]interface{}{
		"ip": ip,
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	}, &Resp{}, false)

	if err != nil {
		return nil, err
	}

	return resStr.(*Resp).Data, nil
}
