package web

import "net/http"

type Server interface {
	http.Handler
	Start(addr string)
}

type HttpServer struct {
}

func (s *HttpServer) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	rsp.Write([]byte{'h', 'e', 'l', 'l', 'o'})
}

func (s *HttpServer) Start(addr string) {
	http.ListenAndServe(addr, s)
}
