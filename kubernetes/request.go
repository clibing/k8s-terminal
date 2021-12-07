package kubernetes

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/k8s-terminal/config"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

const (
	defaultPageNumber int = 1
	defaultPageSize   int = 100
)

/**
 * 对外暴露的请求入口
 */
var (
	Client         *http.Client
	CsrfTokenModel = CsrfToken{}
	JweTokenModel  = JweToken{}
	refreshTokenLock      *sync.RWMutex = new(sync.RWMutex)
)

type BaseModel struct {
	namespace  string   // namespace
	filter     string   // 需要过滤的名字
	name       []string // deployment/service/pod name
	selected   string   // 选择的deployment/service/pod 名字
	pageNumber int      // 请求的当前页
	pageSize   int      // 每页显示记录
	total      int      // 总记录数
	cursor     int      // 上下位置记录
	req        *Request // 实现类
}

// 初始化 http client
func init() {
	if !config.CheckInstalled() {
		//fmt.Println("请先初始化配置,使用`k8s-terminal init --ip $1 --port $2 --token $3")
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	Client = &http.Client{Timeout: 15 * time.Second, Transport: tr}
}

func GetCsrtToken() (value string){
	refreshTokenLock.RLock()
	value = CsrfTokenModel.Value
	refreshTokenLock.RUnlock()
	return
}

func GetJweToken() (value string) {
	refreshTokenLock.RLock()
	value = JweTokenModel.JweToken
	refreshTokenLock.RUnlock()
	return
}

func SetCsrtToken(token string) (){
	refreshTokenLock.Lock()
	CsrfTokenModel.Value = token
	refreshTokenLock.Unlock()
}

func SetJweToken(jwe string) (){
	refreshTokenLock.Lock()
	JweTokenModel.JweToken = jwe
	refreshTokenLock.Unlock()
}

func addHeader(request *http.Request, post bool, ext map[string]string) {
	request.Header.Add("sec-fetch-site", "same-origin")
	request.Header.Add("sec-fetch-mode", "cors")
	request.Header.Add("sec-fetch-dest", "empty")
	if post {
		request.Header.Add("content-type", "application/json;charset=UTF-8")
	}
	if ext != nil {
		for key, value := range ext {
			request.Header.Add(key, value)
		}
	}
}

func noWithTokenRequest(url string, post bool, parameter interface{}, ext map[string]string) (data string, err error) {
	return commonRequest(url, post, parameter, false, false, ext)
}

func withCsrtTokenRequest(url string, post bool, parameter interface{}, ext map[string]string) (data string, err error) {
	return commonRequest(url, post, parameter, true, false, ext)
}

func withJweTokenRequest(url string, post bool, parameter interface{}, ext map[string]string) (data string, err error) {
	return commonRequest(url, post, parameter, false, true, ext)
}

func commonRequest(url string, post bool, parameter interface{}, appendCsrfToken bool, appendJweToken bool, ext map[string]string) (data string, err error) {
	data, err, _ = baseCommonRequest(url, post, parameter, appendCsrfToken, appendJweToken, ext)
	return
}

/**
 * 发送请求
 */
func baseCommonRequest(url string, post bool, parameter interface{}, appendCsrfToken bool, appendJweToken bool, ext map[string]string) (data string, err error, code int) {

	var req *http.Request
	if post {
		var result []byte
		result, err = json.Marshal(parameter)
		if err != nil {
			err = fmt.Errorf("JSON序列化parameter异常, %s", err)
			return
		}
		reader := bytes.NewReader(result)
		req, err = http.NewRequest("POST", url, reader)
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}

	if err != nil {
		err = fmt.Errorf("创建Reqeuest异常, %s", err)
		return
	}
	addHeader(req, post, ext)
	if appendCsrfToken {
		req.Header.Add("x-csrf-token", GetCsrtToken())
	}
	if appendJweToken {
		req.Header.Add("jwetoken", GetJweToken())
	}
	response, err := Client.Do(req)
	if err != nil {
		err = fmt.Errorf("发起请求异常, %s", err)
		return
	}
	defer response.Body.Close()
	code = response.StatusCode
	if response.StatusCode != 200 {
		err = fmt.Errorf("响应异常, http status: %d, err: %s", response.StatusCode, err)
		return
	}
	value, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("读取响应流程异常, %s", err)
		return
	}
	data = string(value)
	return
}
