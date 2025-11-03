package chatRoom

import (
	"context"
	"errors"
	"log"

	"github.com/wiidz/goutil/helpers/sliceHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
)

const REDIS_KEY = "wechat_robot_user_info"

const CHAT_TEAM_CACHE_PREFIX = "wechat_team_info"

// GetMemberInfo 根据wcID获取用户的信息
// 1.从redis中读取
// 2.从api获取
func (api *Api) GetMemberInfo(ctx context.Context, chatRoomID, userID string) (memberInfo *MemberInfoData, err error) {

	//【1】尝试从redis中获取
	memberInfo, err = api.GetMemberInfoFromRedis(ctx, userID)
	if err != nil || memberInfo != nil {
		return
	}

	//【2】从api获取
	if chatRoomID == "" {
		err = errors.New("获取用户信息，必须在群聊中")
		return
	}
	memberInfo, err = api.GetChatRoomMemberInfo(chatRoomID, userID)
	if err != nil {
		return
	}
	go api.SetMemberInfoToRedis(ctx, memberInfo)
	return
}

// GetMemberInfoFromRedis 从redis中读取用户信息
func (api *Api) GetMemberInfoFromRedis(ctx context.Context, wcID string) (userInfo *MemberInfoData, err error) {

	//【1】读取
	userInfo = nil
	tempStr, err := api.RedisMng.HGetString(ctx, REDIS_KEY, wcID)
	if err != nil {
		return
	}

	if tempStr != "" {
		err = typeHelper.JsonDecodeWithStruct(tempStr, &userInfo)
	}
	return
}

// SetMemberInfoToRedis 写入redis中用户信息
func (api *Api) SetMemberInfoToRedis(ctx context.Context, memberInfo *MemberInfoData) (err error) {

	valueStr, err := typeHelper.JsonEncode(memberInfo)

	log.Println("valueStr", valueStr)
	log.Println("err", err)

	_, err = api.RedisMng.HDel(ctx, REDIS_KEY, []string{
		memberInfo.UserName,
	})
	_, err = api.RedisMng.HSetNX(ctx, REDIS_KEY, memberInfo.UserName, valueStr)

	return
}

func (api *Api) IsMembersChanged(ctx context.Context, chatRoomID string, nowLength int64) (changeAmount int64, err error) {
	cacheLen, err := api.RedisMng.HLen(ctx, CHAT_TEAM_CACHE_PREFIX+"_"+chatRoomID)
	if err != nil {
		return
	}
	changeAmount = nowLength - cacheLen

	return
}

// GetChatRoomInfoFromRedis 从redis中读取群聊信息
func (api *Api) GetChatRoomInfoFromRedis(ctx context.Context, chatRoomID string) (memberInfos []*MemberInfoData, err error) {

	//【1】读取
	memberInfos = []*MemberInfoData{}
	tempMap, err := api.RedisMng.HGetAll(ctx, CHAT_TEAM_CACHE_PREFIX+"_"+chatRoomID)
	if err != nil {
		return
	}

	if len(tempMap) != 0 {
		for k := range tempMap {
			var memberInfo *MemberInfoData
			err = typeHelper.JsonDecodeWithStruct(tempMap[k], &memberInfo)
			if err != nil {
				break
			}
			memberInfos = append(memberInfos, memberInfo)
		}
	}
	return
}

// SetChatRoomInfoToRedis 写入redis中群聊信息
func (api *Api) SetChatRoomInfoToRedis(ctx context.Context, chatRoomID string, memberList []*MemberData) (err error) {

	_, err = api.RedisMng.HDelAll(ctx, CHAT_TEAM_CACHE_PREFIX+"_"+chatRoomID)
	if err != nil {
		return
	}

	for k := range memberList {
		var valueStr string
		valueStr, err = typeHelper.JsonEncode(memberList[k])
		if err != nil {
			return
		}
		_, err = api.RedisMng.HSetNX(ctx, CHAT_TEAM_CACHE_PREFIX+"_"+chatRoomID, memberList[k].UserName, valueStr)
	}

	return
}

// GetMemberQuit 判断哪个用户退出了群聊
func (api *Api) GetMemberQuit(ctx context.Context, chatRoomID string, memberInfo []*MemberData) (quiteMembers []*MemberData, err error) {

	var REDIS_KEY_NAME = CHAT_TEAM_CACHE_PREFIX + "_" + chatRoomID

	//【1】提取缓存中的所有key
	userNames, err := api.RedisMng.HKeys(ctx, REDIS_KEY_NAME)
	if err != nil {
		return
	}

	//【2】循环判断
	for k := range memberInfo {
		var flag = sliceHelper.IndexOfStrSlice(memberInfo[k].UserName, userNames)
		if flag != -1 {
			userNames = append(userNames[:flag], userNames[flag+1:]...)
		}

	}

	//【3】将结果查出并删除
	if len(userNames) == 0 {
		return
	}
	quiteMembers = []*MemberData{}
	for k := range userNames {

		//【3-1】提取数据
		var tempStr string
		tempStr, err = api.RedisMng.HGetString(ctx, REDIS_KEY_NAME, userNames[k])
		if err != nil {
			return
		}

		//【3-2】解码
		var tempInfo *MemberData
		err = typeHelper.JsonDecodeWithStruct(tempStr, &tempInfo)
		if err != nil {
			return
		}
		quiteMembers = append(quiteMembers, tempInfo)
		//【3-3】从cache中删除
	}

	_, err = api.RedisMng.HDel(ctx, REDIS_KEY_NAME, userNames)

	return
}

// GetMemberJoin 判断哪个用户进入了群聊
func (api *Api) GetMemberJoin(ctx context.Context, chatRoomID string, memberInfo []*MemberData) (joinMembers []*MemberData, err error) {

	var REDIS_KEY_NAME = CHAT_TEAM_CACHE_PREFIX + "_" + chatRoomID

	//【1】提取缓存中的所有key
	userNames, err := api.RedisMng.HKeys(ctx, REDIS_KEY_NAME)
	if err != nil {
		return
	}
	if len(userNames) == 0 {
		return memberInfo, nil
	}

	//【2】循环判断
	joinMembers = []*MemberData{}
	for k := range memberInfo {
		var flag = sliceHelper.IndexOfStrSlice(memberInfo[k].UserName, userNames)
		if flag == -1 {
			//memberInfo = append(memberInfo[:flag], memberInfo[flag+1:]...)
			joinMembers = append(joinMembers, memberInfo[k])
		}
	}

	//【3】将新的数据写入数据库
	if len(memberInfo) == 0 {
		return
	}

	api.SetChatRoomInfoToRedis(ctx, chatRoomID, memberInfo)
	//_, err = api.RedisMng.HDel(REDIS_KEY_NAME, userNames)

	return
}
