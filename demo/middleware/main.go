package main

import (
	"../../middleware"
	"log"
	"main/proxy"
	"net/http"
	"net/url"
)

func main() {
	// 实现一个反向代理
	reverseProxy := func(c *middleware.SliceRouterContext) http.Handler {
		rs1 := "http://127.0.0.1:2003/base"
		ur11, err1 := url.Parse(rs1)
		if err1 {
			log.Println(err1)
		}

		rs2 := "http://127.0.0.1:2004/base"
		url2, err2 := url.Parse(rs2)
		if err2 {
			log.Println(err2)
		}

		urls := []*url.URL{ur11, url2}
		return proxy.NewMultipleHostsReverseProxy()
	}
}
