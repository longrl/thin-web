package web

import (
	"encoding/json"
	"log"
)

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}

func AccessLog() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx *Context) {
			next(ctx)
			l := accessLog{
				Host:       ctx.Req.Host,
				Route:      ctx.MatchedRoute,
				Path:       ctx.Req.URL.Path,
				HTTPMethod: ctx.Req.Method,
			}
			val, _ := json.Marshal(l)
			log.Println(string(val))
		}
	}
}
