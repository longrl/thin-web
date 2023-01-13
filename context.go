package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Resp http.ResponseWriter
	Req  *http.Request

	// 缓存的响应部分
	// 这部分数据会在最后刷新
	RespStatusCode int
	RespData       []byte

	PathParams map[string]string
	// 缓存查询数据
	cacheQueryValues url.Values

	MatchedRoute string
}

func (ctx *Context) BindJSON(val any) error {
	if ctx.Req.Body == nil {
		return errors.New("thin-web: body 为 nil")
	}
	decoder := json.NewDecoder(ctx.Req.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(val)
}

func (ctx *Context) FormValue(key string) StringValue {
	if err := ctx.Req.ParseForm(); err != nil {
		return StringValue{err: err}
	}
	return StringValue{val: ctx.Req.FormValue(key)}
}

func (ctx *Context) QueryValue(key string) StringValue {
	if ctx.cacheQueryValues == nil {
		ctx.cacheQueryValues = ctx.Req.URL.Query()
	}
	// get 方法会判空不方便排错
	//v := ctx.cacheQueryValues.Get(key)
	v, ok := ctx.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("thin-web: 找不到这个 key")}
	}
	return StringValue{val: v[0]}
}

func (ctx *Context) PathValue(key string) StringValue {
	v, ok := ctx.PathParams[key]
	if !ok {
		return StringValue{err: errors.New("thin-web: 找不到这个 key")}
	}
	return StringValue{val: v}
}

func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.Resp, cookie)
}

func (ctx *Context) RespJSON(status int, val any) error {
	bs, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx.Resp.WriteHeader(status)
	_, err = ctx.Resp.Write(bs)
	return err
}

func (ctx *Context) RespJSONOK(val any) error {
	return ctx.RespJSON(http.StatusOK, val)
}

type StringValue struct {
	val string
	err error
}

func (str StringValue) String() (string, error) {
	return str.val, str.err
}

func (str StringValue) ToInt64() (int64, error) {
	if str.err != nil {
		return 0, str.err
	}
	return strconv.ParseInt(str.val, 10, 64)
}
