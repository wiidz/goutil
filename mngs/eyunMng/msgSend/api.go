package msgSend

import (
	"errors"
	"fmt"

	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/mngs/eyunMng/base"
	"github.com/wiidz/goutil/structs/networkStruct"
)

func NewApi(config *base.Config) (api *Api) {
	return &Api{
		Config: config,
	}
}

func (api *Api) SendMsg(message MessagePayload) (*SendMsgData, error) {
	if message == nil {
		return nil, errors.New("msgSend: message payload cannot be nil")
	}
	switch msg := message.(type) {
	case *TextMessage:
		if msg.Mention != nil && (msg.Mention.All || len(msg.Mention.IDs) > 0) {
			return api.SendTextAt(msg.WcID, msg.Content, msg.Mention.IDs, msg.Mention.All)
		}
		return api.SendText(msg.WcID, msg.Content)
	case *ImageMessage:
		return api.SendImage(msg.WcID, msg.URL)
	case *VideoMessage:
		if msg.ThumbURL != "" {
			return api.SendVideo(msg.WcID, msg.URL, msg.ThumbURL)
		}
		return api.SendVideo(msg.WcID, msg.URL)
	case *LinkMessage:
		return api.SendURL(msg.WcID, msg.URL, msg.Title, msg.Description, msg.ThumbURL)
	case *NameCardMessage:
		return api.SendNameCard(msg.WcID, msg.NameCardID)
	case *EmojiMessage:
		return api.SendEmoji(msg.WcID, msg.ImageMD5, msg.ImageSize)
	case *AppletContentMessage:
		return api.SendApplet(msg.WcID, msg.Content)
	case *MiniProgramMessage:
		return api.SendApplets(msg.WcID, msg.Param)
	default:
		return nil, fmt.Errorf("msgSend: unsupported message type %T", msg)
	}
}

func (api *Api) sendWithPayload(message MessagePayload) (*SendMsgData, error) {

	body, err := message.BuildBody(api.Config.WechatData.WID)
	if err != nil {
		return nil, err
	}

	endpoint := message.Endpoint()
	if endpoint == "" {
		return nil, errors.New("msgSend: endpoint is empty")
	}

	url := api.Config.Host + endpoint
	headers := map[string]string{
		"Authorization": api.Config.Authorization,
	}

	res, _, _, err := networkHelper.RequestWithStruct(
		networkStruct.Post,
		networkStruct.BodyJson,
		url,
		body,
		headers,
		&SendMsgResp{},
	)
	if err != nil {
		return nil, err
	}

	resp := res.(*SendMsgResp)
	if resp.Code != string(Success) {
		if resp.Message != "" {
			return nil, errors.New(resp.Message)
		}
		return nil, fmt.Errorf("msgSend: request failed with code %s", resp.Code)
	}

	return resp.Data, nil
}

// SendText 发送文本消息
func (api *Api) SendText(wcID, content string) (data *SendMsgData, err error) {
	return api.sendWithPayload(&TextMessage{
		WcID:    wcID,
		Content: content,
	})
}

// SendTextAt 发送文本消息带@
func (api *Api) SendTextAt(wcID, content string, atIDs []string, atAll bool) (data *SendMsgData, err error) {
	return api.sendWithPayload(&TextMessage{
		WcID:    wcID,
		Content: content,
		Mention: &Mention{
			IDs: atIDs,
			All: atAll,
		},
	})
}

// SendImage 发送图片消息
func (api *Api) SendImage(wcID, url string) (data *SendMsgData, err error) {
	return api.sendWithPayload(&ImageMessage{
		WcID: wcID,
		URL:  url,
	})
}

// SendVideo 发送视频消息
func (api *Api) SendVideo(wcID, url string, thumb ...string) (data *SendMsgData, err error) {
	var thumbURL string
	if len(thumb) > 0 {
		thumbURL = thumb[0]
	}
	return api.sendWithPayload(&VideoMessage{
		WcID:     wcID,
		URL:      url,
		ThumbURL: thumbURL,
	})
}

// SendURL 发送链接消息
func (api *Api) SendURL(wcID, url, title, desc, thumbURL string) (data *SendMsgData, err error) {
	return api.sendWithPayload(&LinkMessage{
		WcID:        wcID,
		URL:         url,
		Title:       title,
		Description: desc,
		ThumbURL:    thumbURL,
	})
}

// SendNameCard 发送名片消息
func (api *Api) SendNameCard(wcID, nameCardID string) (data *SendMsgData, err error) {
	return api.sendWithPayload(&NameCardMessage{
		WcID:       wcID,
		NameCardID: nameCardID,
	})
}

// SendEmoji 发送发送emoji表情
func (api *Api) SendEmoji(wcID, imageMd5, imgSize string) (data *SendMsgData, err error) {
	return api.sendWithPayload(&EmojiMessage{
		WcID:      wcID,
		ImageMD5:  imageMd5,
		ImageSize: imgSize,
	})
}

// SendApplet 发送APP类消息
func (api *Api) SendApplet(wcID, content string) (data *SendMsgData, err error) {
	return api.sendWithPayload(&AppletContentMessage{
		WcID:    wcID,
		Content: content,
	})
}

// SendApplets 发送小程序
func (api *Api) SendApplets(wcID string, param *SendMiniParam) (data *SendMsgData, err error) {
	return api.sendWithPayload(&MiniProgramMessage{
		WcID:  wcID,
		Param: param,
	})
}

// RevokeMsg 撤回消息
func (api *Api) RevokeMsg(wcID, msgID, newMsgID, createTime string) (data *SendMsgData, err error) {
	//【1】组合URL
	var URL = api.Config.Host + "/sendApplet"

	//【2】请求数据

	res, _, _, err := networkHelper.RequestWithStruct(networkStruct.Post, networkStruct.BodyJson, URL, map[string]interface{}{
		"wId":        api.Config.WechatData.WID, // 登录实例标识
		"wcId":       wcID,                      // 接收人微信id/群id
		"msgID":      msgID,                     // 消息msgId(发送类接口返回的msgId)
		"newMsgId":   newMsgID,                  // 消息newMsgId(发送类接口返回的msgId)
		"createTime": createTime,                // 发送时间（选填）
	}, map[string]string{
		"Authorization": api.Config.Authorization,
	}, &SendMsgResp{})
	if err != nil {
		return
	}

	//【3】判断
	resp := res.(*SendMsgResp)
	if resp.Code == string(Success) {

	} else {
		err = errors.New(resp.Message)
	}

	data = resp.Data

	return
}
