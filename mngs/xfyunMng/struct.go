package xfyunMng

import (
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wiidz/goutil/structs/configStruct"
)

const (
	defaultScheme        = "wss"
	defaultHost          = "iat.xf-yun.com"
	defaultPath          = "/v1"
	defaultFrameSize     = 1280
	defaultFrameInterval = 40 * time.Millisecond
	defaultSeqStart      = 1
)

// FrameStatus 音频帧状态
// 0: 首帧 1: 中间帧 2: 结束帧
// see https://www.xfyun.cn/doc/asr/quickstart.html
// (根据接口文档约定)
type FrameStatus int

const (
	FrameStatusStart    FrameStatus = 0
	FrameStatusContinue FrameStatus = 1
	FrameStatusEnd      FrameStatus = 2
)

// XFYunMng 科大讯飞语音识别管理器
// 负责签名、连接、音频帧发送与结果汇总
// 需提供配置 (configStruct.XFYunConfig)
type XFYunMng struct {
	Config *configStruct.XFYunConfig
	Dialer *websocket.Dialer
}

// RecognizeOptions 识别可选参数
// 可覆盖默认参数（domain、language 等）
// 以及帧大小、帧间隔等音频发送策略
type RecognizeOptions struct {
	ResID string // 热词ID，可选

	Parameter *IATParameter // 识别模式参数
	Audio     *AudioConfig  // 音频参数

	FrameSize     int           // 单帧字节长度，默认1280
	FrameInterval time.Duration // 帧发送间隔，默认40ms
	SeqStart      int           // 初始帧序号，默认1

	Debug bool // 是否打印调试日志
}

// Request 上行请求结构
// header / parameter / payload 结构
// 参考讯飞文档
// https://www.xfyun.cn/doc/asr/dictation/API.html
type Request struct {
	Header    RequestHeader     `json:"header"`
	Parameter *RequestParameter `json:"parameter,omitempty"`
	Payload   *RequestPayload   `json:"payload"`
}

// RequestHeader 请求头
// status 对应音频帧状态
// app_id 必传
// res_id 可选
// status: 0-首帧,1-中间,2-结束
type RequestHeader struct {
	AppID  string      `json:"app_id"`
	ResID  string      `json:"res_id,omitempty"`
	Status FrameStatus `json:"status"`
}

// RequestParameter 请求参数（服务配置）
type RequestParameter struct {
	IAT *IATParameter `json:"iat"`
}

// IATParameter 语音听写参数
// domain / language / accent 必传
// 其他根据业务按需设置
type IATParameter struct {
	Domain   string           `json:"domain"`
	Language string           `json:"language"`
	Accent   string           `json:"accent"`
	EOS      int              `json:"eos,omitempty"`
	LTC      int              `json:"ltc,omitempty"`
	VInfo    int              `json:"vinfo,omitempty"`
	Dwa      string           `json:"dwa,omitempty"`
	Dhw      string           `json:"dhw,omitempty"`
	Result   *ResultParameter `json:"result,omitempty"`
}

// ResultParameter 响应结果参数设置
type ResultParameter struct {
	Encoding string `json:"encoding"`
	Compress string `json:"compress"`
	Format   string `json:"format"`
}

// RequestPayload 请求载荷
type RequestPayload struct {
	Audio *AudioPayload `json:"audio"`
}

// AudioConfig 音频基础配置
// encoding raw/lame
type AudioConfig struct {
	Encoding   string `json:"encoding"`
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
	BitDepth   int    `json:"bit_depth"`
}

// AudioPayload 音频帧主体
type AudioPayload struct {
	AudioConfig
	Seq    int         `json:"seq"`
	Status FrameStatus `json:"status"`
	Audio  string      `json:"audio"`
}

// Response 服务端返回结构
type Response struct {
	Header  ResponseHeader   `json:"header"`
	Payload *ResponsePayload `json:"payload,omitempty"`
}

// ResponseHeader 返回头部信息
type ResponseHeader struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	SID     string      `json:"sid"`
	Status  FrameStatus `json:"status"`
}

// ResponsePayload 返回载荷
type ResponsePayload struct {
	Result *ResponseResult `json:"result,omitempty"`
}

// ResponseResult 识别结果结构
type ResponseResult struct {
	Compress string      `json:"compress"`
	Encoding string      `json:"encoding"`
	Format   string      `json:"format"`
	Seq      int         `json:"seq"`
	Status   FrameStatus `json:"status"`
	Text     string      `json:"text"`
}

// RecognitionResult 汇总结果
// Text 为汇总后的完整文本
// Sentences 记录每个分片信息
// Frames 返回所有原始响应
type RecognitionResult struct {
	SID       string
	Text      string
	Sentences []SentenceResult
	Frames    []*Response
}

// SentenceResult 单个分片结果
type SentenceResult struct {
	Sn     int
	Text   string
	IsLast bool
}

// APIError 讯飞接口错误
type APIError struct {
	Code    int
	Message string
	SID     string
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("xfyun api error: code=%d message=%s sid=%s", e.Code, e.Message, e.SID)
}

// TextResult 音频片段解析结果（payload.result.text 解码后）
type TextResult struct {
	Sn  int      `json:"sn"`
	Ls  bool     `json:"ls"`
	Pgs string   `json:"pgs,omitempty"`
	Rg  []int    `json:"rg,omitempty"`
	Ws  []textWS `json:"ws"`
}

// PlainText 提取纯文本
func (t *TextResult) PlainText() string {
	if t == nil {
		return ""
	}
	var builder strings.Builder
	for _, ws := range t.Ws {
		for _, cw := range ws.Cw {
			builder.WriteString(cw.W)
		}
	}
	return builder.String()
}

type textWS struct {
	Cw []textWord `json:"cw"`
}

type textWord struct {
	W  string `json:"w"`
	Lg string `json:"lg,omitempty"`
}
