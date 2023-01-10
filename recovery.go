package web

import (
	"log"
	"net/http"
)

func Recovery() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = http.StatusInternalServerError
					ctx.RespData = []byte("thin-web: 发生 panic 了")
					log.Println(ctx.Req.URL.Path)
				}
			}()
			next(ctx)
		}
	}
}
