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

// GetRegionInfo 根据IP获取区域信息
func (mng *AudioToTextMng) GetRegionInfo(voice, voiceUrl, format string) (data *Resp, err error) {

	resStr, _, _, err := networkHelper.RequestJsonWithStruct(networkStruct.Get, URL, map[string]interface{}{
		"voice":    voice,    // 语音文件，不超过1MB，和voiceUrl二选一
		"voiceUrl": voiceUrl, // 音频文件url，下载音频不超过1MB，和voice二选一
		"format":   format,   // 语音文件的格式，pcm/wav/amr/m4a。不区分大小写。推荐pcm文件
	}, map[string]string{
		"Authorization": "APPCODE " + mng.Config.AppCode,
	}, &Resp{}, false)

	if err != nil {
		return nil, err
	}

	return resStr.(*Resp), nil
}
