package msgSend

import (
	"errors"
	"strings"

	"github.com/wiidz/goutil/mngs/eyunMng/base"
)

type ReturnCode string

const Success ReturnCode = "1000"
const Failed ReturnCode = "1001"

type Api struct {
	Config *base.Config
}

type MsgKind string

const (
	Text        MsgKind = "text"
	Image       MsgKind = "image"
	Audio       MsgKind = "audio"
	Video       MsgKind = "video"
	Link        MsgKind = "link"
	NameCard    MsgKind = "name_card"
	Emoji       MsgKind = "emoji"
	Applet      MsgKind = "applet"
	MiniProgram MsgKind = "mini_program"
)

type SendMiniParam struct {
	DisplayName string `json:"displayName"` // 小程序的名称，例如：京东
	IconUrl     string `json:"iconUrl"`     // 小程序卡片图标的url(50KB以内的png/jpg)
	AppId       string `json:"appId"`       // 小程序的appID,例如：wx7c544xxxxxx
	PagePath    string `json:"pagePath"`    // 点击小程序卡片跳转的url
	ThumbUrl    string `json:"thumbUrl"`    // 小程序卡片缩略图的url (50KB以内的png/jpg)
	Title       string `json:"title"`       // 标题
	UserName    string `json:"userName"`    // 小程序所有人的ID,例如：gh_1c0daexxxx@app
}

type SendMsgResp struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Data    *SendMsgData `json:"data"`
}

type SendMsgData struct {
	Type       int    `json:"type"`       // 类型
	MsgID      int64  `json:"msgId"`      // 消息msgId
	NewMsgID   int64  `json:"newMsgId"`   // 消息newMsgId
	CreateTime int    `json:"createTime"` // 消息发送时间戳
	WcID       string `json:"wcId"`       // 消息接收方id
}

type Mention struct {
	IDs []string `json:"ids"`
	All bool     `json:"all"`
}

type MessagePayload interface {
	Kind() MsgKind
	Endpoint() string
	BuildBody(wID string) (map[string]interface{}, error)
}

func buildBaseBody(wID, wcID string) (map[string]interface{}, error) {
	if strings.TrimSpace(wID) == "" {
		return nil, errors.New("msgSend: wechat wId is empty")
	}
	if strings.TrimSpace(wcID) == "" {
		return nil, errors.New("发送失败，对象主体不能为空")
	}
	return map[string]interface{}{
		"wId":  wID,
		"wcId": wcID,
	}, nil
}

type TextMessage struct {
	WcID    string   `json:"wc_id"`
	Content string   `json:"content"`
	Mention *Mention `json:"mention,omitempty"`
}

func (m *TextMessage) Kind() MsgKind { return Text }

func (m *TextMessage) Endpoint() string { return "/sendText" }

func (m *TextMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: text message payload is nil")
	}
	if strings.TrimSpace(m.Content) == "" {
		return nil, errors.New("发送失败，内容不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["content"] = m.Content
	if m.Mention != nil {
		if m.Mention.All {
			body["at"] = "notify@all"
		} else if len(m.Mention.IDs) > 0 {
			body["at"] = strings.Join(m.Mention.IDs, ",")
		}
	}
	return body, nil
}

type ImageMessage struct {
	WcID string `json:"wc_id"`
	URL  string `json:"url"`
}

func (m *ImageMessage) Kind() MsgKind { return Image }

func (m *ImageMessage) Endpoint() string { return "/sendImage" }

func (m *ImageMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: image message payload is nil")
	}
	if strings.TrimSpace(m.URL) == "" {
		return nil, errors.New("发送失败，图片地址不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["content"] = m.URL
	return body, nil
}

type VideoMessage struct {
	WcID     string `json:"wc_id"`
	URL      string `json:"url"`
	ThumbURL string `json:"thumb_url,omitempty"`
}

func (m *VideoMessage) Kind() MsgKind { return Video }

func (m *VideoMessage) Endpoint() string { return "/sendVideo" }

func (m *VideoMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: video message payload is nil")
	}
	if strings.TrimSpace(m.URL) == "" {
		return nil, errors.New("发送失败，视频地址不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["path"] = m.URL
	if m.ThumbURL != "" {
		body["thumbPath"] = m.ThumbURL
	} else {
		body["thumbPath"] = ""
	}
	return body, nil
}

type LinkMessage struct {
	WcID        string `json:"wc_id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ThumbURL    string `json:"thumb_url"`
}

func (m *LinkMessage) Kind() MsgKind { return Link }

func (m *LinkMessage) Endpoint() string { return "/sendVideo" }

func (m *LinkMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: link message payload is nil")
	}
	if strings.TrimSpace(m.URL) == "" {
		return nil, errors.New("发送失败，链接地址不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["url"] = m.URL
	body["title"] = m.Title
	body["description"] = m.Description
	body["thumbUrl"] = m.ThumbURL
	return body, nil
}

type NameCardMessage struct {
	WcID       string `json:"wc_id"`
	NameCardID string `json:"name_card_id"`
}

func (m *NameCardMessage) Kind() MsgKind { return NameCard }

func (m *NameCardMessage) Endpoint() string { return "/sendNameCard" }

func (m *NameCardMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: name card message payload is nil")
	}
	if strings.TrimSpace(m.NameCardID) == "" {
		return nil, errors.New("发送失败，名片ID不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["nameCardId"] = m.NameCardID
	return body, nil
}

type EmojiMessage struct {
	WcID      string `json:"wc_id"`
	ImageMD5  string `json:"image_md5"`
	ImageSize string `json:"image_size"`
}

func (m *EmojiMessage) Kind() MsgKind { return Emoji }

func (m *EmojiMessage) Endpoint() string { return "/sendEmoji" }

func (m *EmojiMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: emoji message payload is nil")
	}
	if strings.TrimSpace(m.ImageMD5) == "" {
		return nil, errors.New("发送失败，表情信息不完整：缺少imageMd5")
	}
	if strings.TrimSpace(m.ImageSize) == "" {
		return nil, errors.New("发送失败，表情信息不完整：缺少imgSize")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["imageMd5"] = m.ImageMD5
	body["imgSize"] = m.ImageSize
	return body, nil
}

type AppletContentMessage struct {
	WcID    string `json:"wc_id"`
	Content string `json:"content"`
}

func (m *AppletContentMessage) Kind() MsgKind { return Applet }

func (m *AppletContentMessage) Endpoint() string { return "/sendApplet" }

func (m *AppletContentMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: applet content message payload is nil")
	}
	if strings.TrimSpace(m.Content) == "" {
		return nil, errors.New("发送失败，小程序内容不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["content"] = m.Content
	return body, nil
}

type MiniProgramMessage struct {
	WcID  string         `json:"wc_id"`
	Param *SendMiniParam `json:"param"`
}

func (m *MiniProgramMessage) Kind() MsgKind { return MiniProgram }

func (m *MiniProgramMessage) Endpoint() string { return "/sendApplet" }

func (m *MiniProgramMessage) BuildBody(wID string) (map[string]interface{}, error) {
	if m == nil {
		return nil, errors.New("msgSend: mini program message payload is nil")
	}
	if m.Param == nil {
		return nil, errors.New("发送失败，小程序参数不能为空")
	}
	body, err := buildBaseBody(wID, m.WcID)
	if err != nil {
		return nil, err
	}
	body["displayName"] = m.Param.DisplayName
	body["iconUrl"] = m.Param.IconUrl
	body["appId"] = m.Param.AppId
	body["pagePath"] = m.Param.PagePath
	body["thumbUrl"] = m.Param.ThumbUrl
	body["title"] = m.Param.Title
	body["userName"] = m.Param.UserName
	return body, nil
}
