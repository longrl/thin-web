package web

import (
	"log"
	"net/http"
)

type Server interface {
	http.Handler
	Start(addr string)

	addRoute(method string, path string, handler HandlerFunc)
}

type HandlerFunc func(ctx *Context)

type HttpServer struct {
	router
	ms []Middleware
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		router: newRouter(),
	}
}

func (s *HttpServer) Use(ms ...Middleware) {
	if s.ms == nil {
		s.ms = ms
		return
	}
	s.ms = append(s.ms, ms...)
}

func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	// 路由注册的处理，在 middleware 后执行
	root := s.serve

	for i := len(s.ms) - 1; i >= 0; i-- {
		root = s.ms[i](root)
	}

	var m Middleware = func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	// 将回写数据放到最后执行
	root = m(root)
	root(ctx)
}

func (s *HttpServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("未定义路由"))
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchedRoute = mi.n.route
	mi.n.handle(ctx)
}

func (s *HttpServer) Get(path string, handler HandlerFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HttpServer) Post(path string, handler HandlerFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HttpServer) Start(addr string) {
	http.ListenAndServe(addr, s)
}

func (s *HttpServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		log.Fatalln("回写响应失败", err)
	}
}
