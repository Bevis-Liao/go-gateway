package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var addr = "127.0.0.1:2002"

func main() {
	rs1 := "http://127.0.0.1:2003/base"
	urlInfo, err := url.Parse(rs1)
	if err != nil {
		log.Println(err)
	}

	// 使用 go 的自带代理 reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(urlInfo)
	log.Println("Starting httpserver at " + addr)
	log.Fatal(http.ListenAndServe(addr, proxy))
}
