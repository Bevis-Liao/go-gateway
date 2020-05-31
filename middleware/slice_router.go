package middleware

import (
	"context"
	"math"
	"net/http"
	"strings"
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

func newSliceRouteContext(rw http.ResponseWriter, req *http.Request, r *SliceRouter) *SliceRouterContext {
	newSliceGroup := &SliceGroup{}

	// 最长 url 前缀匹配,最长的等于它
	matchUrlLen := 0
	for _, group := range r.groups {
		// 如果请求的 url 中符合已经设置的 uri
		if strings.HasPrefix(req.RequestURI, group.path) {
			pathLen := len(group.path)
			if pathLen > matchUrlLen {
				matchUrlLen = pathLen
				// path 和 handler
				*newSliceGroup = *group
			}
		}
	}

	c := &SliceRouterContext{Rw: rw, Req:req, Ctx:req.Context(), SliceGroup: newSliceGroup}
	c.Reset()
	return c
}

func (c *SliceRouterContext) Reset()  {
	c.index = -1
}

