package proxy

import (
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var transport = &http.Transport{
	DialContext: (&net.Dialer{
		Timeout:       30 * time.Second, // 连接超时
		KeepAlive:     30 * time.Second, // 长连接超时时间
	}).DialContext,
	MaxIdleConns: 100, // 最大空闲链接
	IdleConnTimeout: 90 * time.Second, // 空闲超时时间
	TLSHandshakeTimeout: 10 * time.Second, // tls 握手超时时间
	ExpectContinueTimeout: 1 * time.Second, // 100-continue 超时时间
}

func NewMultipleHostsReverseProxy(targets []*url.URL) *httputil.ReverseProxy {
	// 请求协调者
	director := func(req *http.Request) {
		// url 重写 这个目前不清楚
		// re, _ := regexp.Compile("^dir(.*)")
		// req.URL.Path = re.ReplaceAllString(req.URL.Path, "$1")

		// 随机的负债均衡
		targetIndex := rand.Intn(len(targets))
		target := targets[targetIndex]
		// 获取目标地址的参数
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

		// url 地址重写
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		// 拼接 get 参数
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		// 这是第一层代理，只在第一层代理设置真实的 header 数据
		req.Header.Set("X-Real-Ip", req.RemoteAddr)
	}

	// 更改内容
	modifyFunc := func(r *http.Response) error {
		// 暂时不修改内容
		return nil
	}

	// 错误回调
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "ErrorHandler error:"+err.Error(), 500)
	}

	return &httputil.ReverseProxy{
		Director:       director,
		Transport:      transport,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}