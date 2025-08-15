package kookMng

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/wiidz/goutil/helpers/networkHelper"
	"github.com/wiidz/goutil/helpers/typeHelper"
	"github.com/wiidz/goutil/structs/configStruct"
	"github.com/wiidz/goutil/structs/networkStruct"
)

const Domain = "https://www.kookapp.cn/api/v3"

// NewKookMng 创建一个kook管理器
func NewKookMng(config *configStruct.KookConfig) (mng *KookMng, err error) {
	if config.Token == "" || config.EncryptKey == "" || config.VerifyToken == "" || config.CallbackURL == "" {
		err = errors.New("kook配置参数有误")
		return
	}
	mng = &KookMng{
		Config: config,
	}
	return
}

// DecryptData 解压数据
func (m *KookMng) DecryptData(sourceStr string) (data *ReceiveData, err error) {
	//【1】原始密文base64 解码
	raw, err := base64.StdEncoding.DecodeString(sourceStr)
	if err != nil {
		err = fmt.Errorf("密文base64解码失败: %v", err)
		return
	}
	if len(raw) < 16 {
		err = fmt.Errorf("解码后数据长度不足16字节")
		return
	}

	//【2】截取前16字节为iv，剩下的是新的密文（还需要 base64 解码）
	iv := raw[:16]
	enc2Base64 := raw[16:]
	// 3. 再 base64 解码新的密文
	realCiphertext, err := base64.StdEncoding.DecodeString(string(enc2Base64))
	if err != nil {
		err = fmt.Errorf("内部密文base64解码失败: %v", err)
		return
	}

	//【3】key补0到32字节
	key := []byte(m.Config.EncryptKey)
	if len(key) > 32 {
		key = key[:32]
	} else if len(key) < 32 {
		key = append(key, bytes.Repeat([]byte{0x00}, 32-len(key))...)
	}

	//【4】AES-256-CBC 解密
	block, err := aes.NewCipher(key)
	if err != nil {
		err = fmt.Errorf("创建AES实例失败: %v", err)
		return
	}
	if len(realCiphertext)%aes.BlockSize != 0 {
		err = fmt.Errorf("密文长度不是块大小的整数倍")
		return
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(realCiphertext))
	mode.CryptBlocks(decrypted, realCiphertext)

	//【5】去除PKCS7填充
	var result []byte
	result, err = pkcs7Unpad(decrypted)
	if err != nil {
		err = fmt.Errorf("去除填充失败: %v", err)
		return
	}

	//【6】解析json字符串
	err = typeHelper.JsonDecodeWithStruct(string(result), &data)
	return
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据长度为0")
	}
	padLen := int(data[length-1])
	if padLen == 0 || padLen > length {
		return nil, fmt.Errorf("填充长度解析异常")
	}
	for _, v := range data[length-padLen:] {
		if int(v) != padLen {
			return nil, fmt.Errorf("填充数据格式错误")
		}
	}
	return data[:length-padLen], nil
}

// SendRequest 发送请求
func (m *KookMng) sendRequest(urlSuffix string, params map[string]interface{}, method networkStruct.Method) (resp *ApiResponse, err error) {

	res, _, _, err := networkHelper.RequestJsonWithStruct(method, Domain+urlSuffix, params, map[string]string{
		"Authorization": "Bot " + m.Config.Token,
	}, &ApiResponse{}, m.Config.Debug)
	if err != nil {
		return nil, err
	}

	resp = res.(*ApiResponse)
	if resp.Code != 0 {
		err = errors.New(resp.Message)
	}

	return
}

// EditUserNickname 修改服务器中玩家的昵称
func (m *KookMng) EditUserNickname(guildID, targetID, nickname string) (resp *ApiResponse, err error) {
	data, err := typeHelper.JsonEncodeDecodeMap(&EditUserNicknameParam{
		GuildID:  guildID,
		Nickname: nickname,
		UserID:   targetID,
	})
	if err != nil {
		return
	}
	return m.sendRequest(suffixs.EditUserNickname, data, networkStruct.Post)
}
