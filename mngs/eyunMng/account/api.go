package account

import (
	"errors"

	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/structs/networkStruct"
)

func NewApi(config *base.Config) (api *Api) {
	return &Api{
		Config: config,
	}
}

// QueryLoginWx 查询账号中在线的微信列表
func (api *Api) QueryLoginWx() (data []*LoginWxData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/queryLoginWx"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, nil, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &QueryLoginWxResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*QueryLoginWxResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// IsOnline 查询微信是否在线
func (api *Api) IsOnline(wcID string) (isOnline bool, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/isOnline"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wcId": wcID,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &IsOnlineResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*IsOnlineResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	isOnline = resp.Data.IsOnline

	return
}

// Offline 离线
func (api *Api) Offline(wcIDs []string) (err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/member/offline"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"account": api.Config.Account,
		"wcIds":   wcIDs,
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &base.BaseResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*base.BaseResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	return
}
