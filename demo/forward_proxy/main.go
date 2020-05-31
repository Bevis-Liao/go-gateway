package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type Proxy struct {
	
}

func (p *Proxy)ServeHTTP(rw http.ResponseWriter, req *http.Request)  {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)
	//
	// step 1 重新新建一个对象，然后修改 ip 属性
	outReq := new(http.Request)
	*outReq = *req
	// 将 ip 和端口分开
	if clientIp, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		// fmt.Printf("Client ip : %s, X-Forwarded-For : %s \n", clientIp, outReq.Header["X-Forwarded-For"])
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIp = strings.Join(prior, ", ") + ", " + clientIp
		}
		outReq.Header.Set("X-Forwarded-For", clientIp)
	}

	// step 2. 请求下游
	transport := http.DefaultTransport
	res, err := transport.RoundTrip(outReq)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("transport error : %s \n", err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	// step 3. 把下游的内容返回给上游
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)
}

func main() {
	// 目前还不支持 https 正向代理，代理只成功一半
	fmt.Println("Http Forward Proxy, Listening 7890 port.")
	http.Handle("/", &Proxy{})
	http.ListenAndServe("0.0.0.0:7890", nil)
}
