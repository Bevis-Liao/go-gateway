package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

var (
	// 需要代理的地址，真实服务地址
	proxyAddr = "http://127.0.0.1:2003"
)

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start proxy server at port 2001")
	err := http.ListenAndServe(":2001", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request)  {
	// step 1 解析代理地址，并更改请求体的协议和主机
	proxy, err := url.Parse(proxyAddr)
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host

	// step 2 请求下游
	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r)
	defer resp.Body.Close()
	if err != nil {
		log.Print(err)
	}

	// step 3 把下游请求的内容，返回给上游，为什么这里需要设置 header, step 2 的时候不用
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	bufio.NewReader(resp.Body).WriteTo(w)
}
