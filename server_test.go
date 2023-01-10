package web

import (
	"net/http"
	"testing"
)

type User struct {
	Name string
	Id   int
}

func TestServer(t *testing.T) {
	request := &User{}
	server := NewHttpServer()
	server.Use(Recovery(), AccessLog())
	server.Post("/hello/name", func(ctx *Context) {
		ctx.BindJSON(request)
		ctx.RespJSON(http.StatusNotFound, request)
	})
	server.Get("/hello", func(ctx *Context) {
		s := ctx.QueryValue("name")
		v, err := s.String()
		if err != nil {
			ctx.RespJSON(http.StatusNotFound, err.Error())
		} else {
			ctx.RespJSONOK(v)
		}
	})
	server.Get("/hello/name/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, star"))
	})
	server.Get("/hello/:id", func(ctx *Context) {
		ctx.RespJSONOK(ctx.PathValue("id").val)
	})
	server.Start(":8080")
}
