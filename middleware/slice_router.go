package middleware

import (
	"context"
	"math"
	"net/http"
)

const abortIndex int8 = math.MaxInt8 / 2

type HandlerFunc func(*SliceRouterContext)

// router 结构体
type SliceRouter struct {
	groups []*SliceGroup
}

// 中间件组
type SliceGroup struct {
	*SliceRouter
	path string
	handlers []HandlerFunc
}

// router 上下文
type SliceRouterContext struct {
	Rw http.ResponseWriter
	Req *http.Request
	Ctx context.Context
	*SliceGroup
	index int8
}

func newSliceRouteContext(rw http.ResponseWriter, req *http.Request)  {
	
}