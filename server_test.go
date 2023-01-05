package web

import (
	"testing"
)

func TestServer(t *testing.T) {
	server := NewHttpServer()
	server.Get("/hello/name", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	server.Get("/hello", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello!!!"))
	})
	server.Get("/hello/name/*", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, star"))
	})
	server.Get("/hello/:id", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, path id = " + ctx.PathParams["id"]))
	})
	server.Start(":8080")
}
