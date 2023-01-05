package web

import (
	"net/http"
	"testing"
)

func Test_router_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/hello/user",
		},
		{
			method: http.MethodGet,
			path:   "/hello/name",
		},
		{
			method: http.MethodGet,
			path:   "/user/info",
		},
		{
			method: http.MethodPost,
			path:   "/user/login",
		},
	}

	mockHandler := func(ctx *Context) {}
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRoute(tr.method, tr.path, mockHandler)
	}

	//wantRouters := &router{
	//	trees: map[string]*node{
	//		http.MethodGet: &node{
	//			path: "/",
	//			children: map[string]*node{
	//				"hello": &node{
	//					path: "hello",
	//					children: map[string]*node{
	//						"user": &node{
	//							path:   "user",
	//							handle: mockHandler,
	//						},
	//						"name": &node{
	//							path:   "name",
	//							handle: mockHandler,
	//						},
	//					},
	//				},
	//				"user": &node{
	//					path: "user",
	//					children: map[string]*node{
	//						"info": &node{
	//							path:   "info",
	//							handle: mockHandler,
	//						},
	//					},
	//				},
	//			},
	//		},
	//		http.MethodPost: &node{
	//			path: "/",
	//			children: map[string]*node{
	//				"user": &node{
	//					path: "user",
	//					children: map[string]*node{
	//						"login": &node{
	//							path:   "user",
	//							handle: mockHandler,
	//						},
	//					},
	//				},
	//			},
	//		},
	//	},
	//}

}
