package audioToTextMng

import (
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

func New(config *configStruct.AliApiConfig) *AudioToTextMng {
	return &AudioToTextMng{
		Config: config,
	}
}

const URL = "https://smyuyin.market.alicloudapi.com/v2/voice_to_text/generate"

// Generate 语音转文字
func (mng *AudioToTextMng) Generate(voice, voiceUrl, format string) (data *Resp, err error) {

	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Post, URL, map[string]interface{}{
		"voice":    voice,    // 语音文件，不超过1MB，和voiceUrl二选一
		"voiceUrl": voiceUrl, // 音频文件url，下载音频不超过1MB，和voice二选一
		"format":   format,   // 语音文件的格式，pcm/wav/amr/m4a。不区分大小写。推荐pcm文件
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	}, &Resp{}, true)

	if err != nil {
		return nil, err
	}

	return resStr.(*Resp), nil
}

// Generate 语音转文字
func (mng *AudioToTextMng) GenerateRaw(voice, voiceUrl, format string) (data string, err error) {

	resStr, _, _, err := networkHelper.RequestRaw(networkStruct.Post, URL, map[string]interface{}{
		"voice":    voice,    // 语音文件，不超过1MB，和voiceUrl二选一
		"voiceUrl": voiceUrl, // 音频文件url，下载音频不超过1MB，和voice二选一
		"format":   format,   // 语音文件的格式，pcm/wav/amr/m4a。不区分大小写。推荐pcm文件
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	})

	if err != nil {
		return "", err
	}

	return resStr, nil
}
