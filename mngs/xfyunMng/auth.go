package xfyunMng

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const rfc1123GMT = "Mon, 02 Jan 2006 15:04:05 GMT"

// buildAuthURL 构建带鉴权参数的websocket地址
func (mng *XFYunMng) buildAuthURL(ts time.Time) (string, error) {
	if mng == nil {
		return "", errors.New("xfyun manager is nil")
	}
	if mng.Config == nil {
		return "", errors.New("xfyun config is nil")
	}
	if mng.Config.ApiKey == "" || mng.Config.ApiSecret == "" || mng.Config.AppID == "" {
		return "", errors.New("xfyun config missing credentials")
	}

	scheme, host, path := mng.resolveEndpoint()
	if host == "" {
		return "", errors.New("xfyun host is empty")
	}

	date := ts.UTC().Format(rfc1123GMT)
	requestLine := fmt.Sprintf("GET %s HTTP/1.1", path)
	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\n%s", host, date, requestLine)

	mac := hmac.New(sha256.New, []byte(mng.Config.ApiSecret))
	_, err := mac.Write([]byte(signatureOrigin))
	if err != nil {
		return "", fmt.Errorf("sign signature origin failed: %w", err)
	}
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authorizationOrigin := fmt.Sprintf(`api_key="%s", algorithm="hmac-sha256", headers="host date request-line", signature="%s"`, mng.Config.ApiKey, signature)
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))

	values := url.Values{}
	values.Set("authorization", authorization)
	values.Set("date", date)
	values.Set("host", host)

	u := url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     path,
		RawQuery: values.Encode(),
	}
	return u.String(), nil
}

// resolveEndpoint 解析最终访问地址（含默认值）
func (mng *XFYunMng) resolveEndpoint() (scheme, host, path string) {
	scheme = defaultScheme
	host = defaultHost
	path = defaultPath

	if mng != nil && mng.Config != nil {
		if temp := strings.TrimSpace(mng.Config.Scheme); temp != "" {
			scheme = strings.ToLower(temp)
		}
		if temp := strings.TrimSpace(mng.Config.Host); temp != "" {
			host = temp
		}
		if temp := strings.TrimSpace(mng.Config.Path); temp != "" {
			path = temp
		}
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return
}
