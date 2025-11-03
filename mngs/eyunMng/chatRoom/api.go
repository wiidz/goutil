package chatRoom

import (
	"errors"
	"log"

	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/mngs/redisMng"
	"github.com/wiidz/goutil/structs/networkStruct"
)

func NewApi(config *base.Config, redisM *redisMng.RedisMng) (api *Api) {
	return &Api{
		Config:   config,
		RedisMng: redisM,
	}
}

// GetChatRoomMember 获取群成员
func (api *Api) GetChatRoomMember(chatRoomID string) (data []*MemberData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/getChatRoomMember"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"chatRoomId": chatRoomID,                // 群号
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &GetChatRoomMemberResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*GetChatRoomMemberResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

// GetChatRoomMemberInfo 获取群成员详情
func (api *Api) GetChatRoomMemberInfo(chatRoomID string, userID string) (data *MemberInfoData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/getChatRoomMemberInfo"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"chatRoomId": chatRoomID,                // 群号
		"userList":   userID,                    // 群成员标识  PS: 暂不支持多个群成员查询，可间隔调用获取
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &GetChatRoomMemberInfoResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*GetChatRoomMemberInfoResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	if len(resp.Data) == 0 {
		err = errors.New("获取用户信息失败")
		return
	}

	data = resp.Data[0]

	return
}

// GetChatRoomInfo 获取群信息
func (api *Api) GetChatRoomInfo(chatRoomID string) (data []*ChatRoomInfoData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/getChatRoomMemberInfo"

	//【2】请求数据
	res, _, _, err := networkHelper.RequestWithStructTest(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"chatRoomId": chatRoomID,                // 群号
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &GetChatRoomInfoResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*GetChatRoomInfoResp)
	if resp.Code == string(base.Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}

func (api *Api) Test() {
	log.Println("test", api.Config.WechatData.WID)
}
