package xfyunMng

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wiidz/goutil/structs/configStruct"
)

// New 创建讯飞语音识别管理器
func New(config *configStruct.XFYunConfig) (*XFYunMng, error) {
	if config == nil {
		return nil, errors.New("xfyun config is nil")
	}
	if config.AppID == "" || config.ApiKey == "" || config.ApiSecret == "" {
		return nil, errors.New("xfyun config missing credentials")
	}

	return &XFYunMng{
		Config: config,
		Dialer: &websocket.Dialer{
			HandshakeTimeout: 10 * time.Second,
		},
	}, nil
}

// Recognize 发送音频流并返回识别结果
func (mng *XFYunMng) Recognize(ctx context.Context, reader io.Reader, opts *RecognizeOptions) (*RecognitionResult, error) {
	if mng == nil {
		return nil, errors.New("xfyun manager is nil")
	}
	if reader == nil {
		return nil, errors.New("audio reader is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	options := mergeRecognizeOptions(opts)
	if mng.Config != nil && mng.Config.Debug {
		options.Debug = true
	}

	wsURL, err := mng.buildAuthURL(time.Now())
	if err != nil {
		return nil, err
	}

	if options.Debug {
		log.Printf("xfyun connecting to %s", wsURL)
	}

	dialer := mng.Dialer
	if dialer == nil {
		dialer = &websocket.Dialer{HandshakeTimeout: 10 * time.Second}
	}

	conn, resp, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		if resp != nil && resp.Body != nil {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			if len(body) > 0 {
				err = fmt.Errorf("dial failed: %w: %s", err, strings.TrimSpace(string(body)))
			}
		}
		return nil, err
	}
	defer conn.Close()
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	responses := make(chan responseMessage)
	go streamResponses(ctx, conn, responses)

	if options.Debug {
		log.Printf("xfyun start streaming audio, frame_size=%d, interval=%s", options.FrameSize, options.FrameInterval)
	}

	if _, err := mng.sendAudioFrames(ctx, conn, reader, options); err != nil {
		return nil, err
	}

	aggregator := newSentenceAggregator()
	var (
		sid    string
		frames []*Response
	)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case msg, ok := <-responses:
			if !ok {
				text, sentences := aggregator.summary()
				if options.Debug {
					log.Printf("xfyun recognition completed sid=%s text=%s", sid, text)
				}
				return &RecognitionResult{
					SID:       sid,
					Text:      text,
					Sentences: sentences,
					Frames:    frames,
				}, nil
			}
			if msg.Err != nil {
				return nil, msg.Err
			}

			frame := msg.Frame
			frames = append(frames, frame)

			if frame.Header.Code != 0 {
				return nil, &APIError{
					Code:    frame.Header.Code,
					Message: frame.Header.Message,
					SID:     frame.Header.SID,
				}
			}

			if sid == "" {
				sid = frame.Header.SID
			}

			if frame.Payload != nil && frame.Payload.Result != nil && frame.Payload.Result.Text != "" {
				textResult, err := parseTextResult(frame.Payload.Result)
				if err != nil {
					return nil, err
				}
				aggregator.apply(textResult)
			}

			if frame.Header.Status == FrameStatusEnd {
				text, sentences := aggregator.summary()
				if options.Debug {
					log.Printf("xfyun recognition end sid=%s text=%s", sid, text)
				}
				return &RecognitionResult{
					SID:       sid,
					Text:      text,
					Sentences: sentences,
					Frames:    frames,
				}, nil
			}
		}
	}
}

// mergeRecognizeOptions 合并自定义与默认参数
func mergeRecognizeOptions(opts *RecognizeOptions) RecognizeOptions {
	base := defaultRecognizeOptions()
	if opts == nil {
		return base
	}

	if opts.ResID != "" {
		base.ResID = opts.ResID
	}
	if opts.FrameSize > 0 {
		base.FrameSize = opts.FrameSize
	}
	if opts.FrameInterval > 0 {
		base.FrameInterval = opts.FrameInterval
	}
	if opts.SeqStart > 0 {
		base.SeqStart = opts.SeqStart
	}
	if opts.Parameter != nil {
		p := opts.Parameter
		if p.Domain != "" {
			base.Parameter.Domain = p.Domain
		}
		if p.Language != "" {
			base.Parameter.Language = p.Language
		}
		if p.Accent != "" {
			base.Parameter.Accent = p.Accent
		}
		if p.EOS > 0 {
			base.Parameter.EOS = p.EOS
		}
		if p.LTC > 0 {
			base.Parameter.LTC = p.LTC
		}
		if p.VInfo > 0 {
			base.Parameter.VInfo = p.VInfo
		}
		if p.Dwa != "" {
			base.Parameter.Dwa = p.Dwa
		}
		if p.Dhw != "" {
			base.Parameter.Dhw = p.Dhw
		}
		if p.Result != nil {
			if base.Parameter.Result == nil {
				base.Parameter.Result = &ResultParameter{}
			}
			r := p.Result
			if r.Encoding != "" {
				base.Parameter.Result.Encoding = r.Encoding
			}
			if r.Compress != "" {
				base.Parameter.Result.Compress = r.Compress
			}
			if r.Format != "" {
				base.Parameter.Result.Format = r.Format
			}
		}
	}
	if opts.Audio != nil {
		a := opts.Audio
		if a.Encoding != "" {
			base.Audio.Encoding = a.Encoding
		}
		if a.SampleRate > 0 {
			base.Audio.SampleRate = a.SampleRate
		}
		if a.Channels > 0 {
			base.Audio.Channels = a.Channels
		}
		if a.BitDepth > 0 {
			base.Audio.BitDepth = a.BitDepth
		}
	}
	if opts.Debug {
		base.Debug = true
	}

	return base
}

// defaultRecognizeOptions 默认识别配置
func defaultRecognizeOptions() RecognizeOptions {
	return RecognizeOptions{
		Parameter: &IATParameter{
			Domain:   "slm",
			Language: "zh_cn",
			Accent:   "mandarin",
			EOS:      6000,
			VInfo:    1,
			Result: &ResultParameter{
				Encoding: "utf8",
				Compress: "raw",
				Format:   "json",
			},
		},
		Audio: &AudioConfig{
			Encoding:   "raw",
			SampleRate: 16000,
			Channels:   1,
			BitDepth:   16,
		},
		FrameSize:     defaultFrameSize,
		FrameInterval: defaultFrameInterval,
		SeqStart:      defaultSeqStart,
	}
}

// sendAudioFrames 逐帧发送音频数据
func (mng *XFYunMng) sendAudioFrames(ctx context.Context, conn *websocket.Conn, reader io.Reader, opts RecognizeOptions) (int, error) {
	if opts.Audio == nil {
		return opts.SeqStart, errors.New("audio config is nil")
	}

	frameSize := opts.FrameSize
	if frameSize <= 0 {
		frameSize = defaultFrameSize
	}
	buffer := make([]byte, frameSize)
	seq := opts.SeqStart
	firstFrame := true

	for {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return seq, ctx.Err()
			default:
			}
		}

		n, readErr := io.ReadFull(reader, buffer)

		if readErr != nil {
			if readErr == io.EOF && n == 0 {
				break
			}
			if !errors.Is(readErr, io.EOF) && !errors.Is(readErr, io.ErrUnexpectedEOF) {
				return seq, readErr
			}
		}

		if n > 0 {
			status := FrameStatusContinue
			includeParam := false
			if firstFrame {
				status = FrameStatusStart
				includeParam = true
				firstFrame = false
			}

			audioBase64 := base64.StdEncoding.EncodeToString(buffer[:n])
			if err := mng.writeFrame(ctx, conn, seq, status, audioBase64, includeParam, opts); err != nil {
				return seq, err
			}
			seq++

			if opts.FrameInterval > 0 {
				select {
				case <-ctx.Done():
					return seq, ctx.Err()
				case <-time.After(opts.FrameInterval):
				}
			}
		}

		if errors.Is(readErr, io.EOF) {
			break
		}
		if errors.Is(readErr, io.ErrUnexpectedEOF) {
			break
		}
	}

	if firstFrame {
		// 没有音频数据，也需发送首帧头部
		if err := mng.writeFrame(ctx, conn, seq, FrameStatusStart, "", true, opts); err != nil {
			return seq, err
		}
		seq++
	}

	if err := mng.writeFrame(ctx, conn, seq, FrameStatusEnd, "", false, opts); err != nil {
		return seq, err
	}
	seq++

	return seq, nil
}

// writeFrame 写入单帧数据
func (mng *XFYunMng) writeFrame(ctx context.Context, conn *websocket.Conn, seq int, status FrameStatus, audio string, includeParam bool, opts RecognizeOptions) error {
	frame := Request{
		Header: RequestHeader{
			AppID:  mng.Config.AppID,
			ResID:  opts.ResID,
			Status: status,
		},
		Payload: &RequestPayload{
			Audio: &AudioPayload{
				AudioConfig: *opts.Audio,
				Seq:         seq,
				Status:      status,
				Audio:       audio,
			},
		},
	}

	if includeParam {
		frame.Parameter = &RequestParameter{IAT: opts.Parameter}
	}

	if ctx != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	if opts.Debug {
		log.Printf("xfyun send frame seq=%d status=%d audio_bytes=%d", seq, status, len(audio))
	}

	if err := conn.WriteJSON(frame); err != nil {
		return fmt.Errorf("send frame failed: %w", err)
	}
	return nil
}

// parseTextResult 解析返回的文本分片
func parseTextResult(res *ResponseResult) (*TextResult, error) {
	if res == nil || res.Text == "" {
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(res.Text)
	if err != nil {
		return nil, fmt.Errorf("decode text failed: %w", err)
	}

	var result TextResult
	if err := json.Unmarshal(decoded, &result); err != nil {
		return nil, fmt.Errorf("unmarshal text failed: %w", err)
	}

	return &result, nil
}

// responseMessage 用于异步读取websocket的消息
type responseMessage struct {
	Frame *Response
	Err   error
}

// streamResponses 异步读取响应
func streamResponses(ctx context.Context, conn *websocket.Conn, ch chan<- responseMessage) {
	defer close(ch)

	for {
		var resp Response
		if err := conn.ReadJSON(&resp); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			ch <- responseMessage{Err: err}
			return
		}

		frame := resp

		select {
		case ch <- responseMessage{Frame: &frame}:
		case <-ctx.Done():
			return
		}

		if resp.Header.Status == FrameStatusEnd {
			return
		}
	}
}

// sentenceAggregator 对识别结果进行累计
type sentenceAggregator struct {
	sentences map[int]SentenceResult
}

func newSentenceAggregator() *sentenceAggregator {
	return &sentenceAggregator{
		sentences: make(map[int]SentenceResult),
	}
}

func (agg *sentenceAggregator) apply(result *TextResult) {
	if result == nil {
		return
	}

	if result.Pgs == "rpl" && len(result.Rg) == 2 {
		for sn := result.Rg[0]; sn <= result.Rg[1]; sn++ {
			delete(agg.sentences, sn)
		}
	}

	agg.sentences[result.Sn] = SentenceResult{
		Sn:     result.Sn,
		Text:   result.PlainText(),
		IsLast: result.Ls,
	}
}

func (agg *sentenceAggregator) summary() (string, []SentenceResult) {
	keys := make([]int, 0, len(agg.sentences))
	for sn := range agg.sentences {
		keys = append(keys, sn)
	}
	sort.Ints(keys)

	var builder strings.Builder
	results := make([]SentenceResult, 0, len(keys))
	for _, sn := range keys {
		sentence := agg.sentences[sn]
		builder.WriteString(sentence.Text)
		results = append(results, sentence)
	}

	return builder.String(), results
}
