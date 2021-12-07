package kubernetes

import (
	"encoding/json"
	"fmt"
)

/**
 * auth 认证
 * csrf token: 跨域token (动态请求kubernetes dashboard 获取)
 * auth token: 认证的token (由运维人员提供)
 * jwe  token: 登录成功后返回的token (每次发送http请求都需要使用)
 * 使用token进行登录
 * 具体的登录步骤为
 * 1. 先获取csrf的凭证，即跨域请求的凭证(csrfToken)
 * 2. 使用获取到的跨域凭证(csrfToken)和登录的auth token进行登录获取请求的jwe token
 * 3. jwe token时有过期时间的，默认很短大概60秒,如果需要长时间程序，例如查看日志，需要异步刷新csrf token 和 jwe token
 */

type LoginApi interface {
	CsrfToken() string                // 获取跨域请求
	JweToken(authToken string) string // 使用跨域token和认证token进行认证获取jwe token
	RefreshCsrtToken() string         // 刷新csrf token
	RefreshJweToken() string          // 刷新jwe token
}

func (req *Request) CsrfToken() string {
	url := fmt.Sprintf("https://%s:%d/api/v1/csrftoken/login", req.Ip, req.Port)
	data, err := noWithTokenRequest(url, false, nil, nil)
	if err != nil {
		panic(err)
	}
	var tmp CsrfToken
	json.Unmarshal([]byte(data), &tmp)
	v := tmp.Value
	SetCsrtToken(v)
	return v
}

func (req *Request) JweToken() string {
	login := LoginRequest{
		Username:   "",
		Password:   "",
		KubeConfig: "",
		Token:      req.Token,
	}
	url := fmt.Sprintf("https://%s:%d/api/v1/login", req.Ip, req.Port)
	data, err := withCsrtTokenRequest(url, true, login, nil)
	if err != nil {
		panic(err)
	}
	var tmp JweToken
	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		panic(err)
	}
	v := tmp.JweToken
	SetJweToken(v)
	return v
}

func (req *Request) RefreshCsrtToken() string {
	if JweTokenModel.JweToken == "" {
		err := fmt.Errorf("先获取jwe token后在刷新token")
		panic(err)
	}

	url := fmt.Sprintf("https://%s:%d/api/v1/csrftoken/token", req.Ip, req.Port)
	data, err := withJweTokenRequest(url, false, nil, nil)
	if err != nil {
		panic(err)
	}

	var tmp CsrfToken
	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		panic(err)
	}
	v := tmp.Value
	SetCsrtToken(v)
	return v
}

func (req *Request) RefreshJweToken() string {
	if JweTokenModel.JweToken == "" {
		err := fmt.Errorf("先获取jwe token后在刷新token")
		panic(err)
	}
	oldToken := RJweToken{
		JweToken: JweTokenModel.JweToken,
	}
	url := fmt.Sprintf("https://%s:%d/api/v1/token/refresh", req.Ip, req.Port)
	data, err := commonRequest(url, true, oldToken, true, true, nil)
	if err != nil {
		panic(err)
	}
	var tmp JweToken
	err = json.Unmarshal([]byte(data), &tmp)
	if err != nil {
		panic(err)
	}
	v := tmp.JweToken
	SetJweToken(v)
	return v
}
