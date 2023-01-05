package web

import "net/http"

type Server interface {
	http.Handler
	Start(addr string)

	addRoute(method string, path string, handler HandlerFunc)
}

type HandlerFunc func(ctx *Context)

type HttpServer struct {
	router
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		router: newRouter(),
	}
}

func (s *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	s.serve(ctx)
}

func (s *HttpServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("未定义路由"))
		return
	}
	ctx.PathParams = mi.pathParams
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
