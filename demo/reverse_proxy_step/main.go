package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var addr = "127.0.0.1:2002"

// 1, 先使用门户代理租户
// 2, 修改代理请求
// 3, 修改门户的返回值

func main() {
	rs1 := "http://ycz-test.myyscm.com/hscz?source=ycg"
	// rs1 = "http://www.zhelinmg.com"
	urlInfo, err := url.Parse(rs1)
	if err != nil {
		log.Println(err)
	}

	// fmt.Println(urlInfo.Path, urlInfo.RawQuery)
	// url 参数
	targetQuery := urlInfo.RawQuery
	director := func(r *http.Request) {
		r.URL.Scheme = urlInfo.Scheme
		r.URL.Host = urlInfo.Host
		r.URL.Path = urlInfo.Path
		// r.URL.Path = singleJoiningSlash(urlInfo.Path, r.URL.Path)
		// 拼接参数
		if targetQuery == "" || r.URL.RawQuery == "" {
			r.URL.RawQuery = targetQuery + r.URL.RawQuery
		} else {
			r.URL.RawQuery = targetQuery + "&" + r.URL.RawQuery
		}
		// 如果没有设置请求UserAgent
		r.Header.Set("User-Agent", "Go-client" )
		if _, ok := r.Header["User-Agent"]; !ok {
			r.Header.Set("User-Agent", "Go-client" )
		}
		log.Println(r.URL.RawQuery)
		log.Println(r.URL.Path)
	}

	// 处理返回 302 response 头
	modifyFunc := func(res *http.Response) error {
		//oldPayload, err := ioutil.ReadAll(res.Body)
		//if err != nil {
		//	log.Println(oldPayload)
		//}
		//if res.StatusCode {
		//
		//}

		log.Println(res.StatusCode, res.Request.URL)

		return nil
	}

	// httputil.NewSingleHostReverseProxy()
	// 使用 go 的自带代理 reverse proxy
	// proxy := NewSingleHostReverseProxy(urlInfo)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, &httputil.ReverseProxy{Director: director, ModifyResponse:modifyFunc}))
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		// a/
		// /b
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

