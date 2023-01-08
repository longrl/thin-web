package web

import (
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

func (r *router) addRoute(method string, path string, handle HandlerFunc) {
	if path == "" {
		panic("thin-web: 路由是空字符串")
	}
	if path[0] != '/' {
		panic("thin-web: 路由必须以 / 开头")
	}
	if path != "/" && path[len(path)-1] == '/' {
		panic("thin-web: 路由不能以 / 结尾")
	}
	root, ok := r.trees[method]
	if !ok {
		root = &node{path: "/"}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handle != nil {
			panic("thin-web: 路由冲突[/]")
		}
		root.handle = handle
		return
	}
	// 路由树按照 ‘/’ 切割构造
	segments := strings.Split(path[1:], "/")
	for _, segment := range segments {
		if segment == "" {
			panic(fmt.Sprintf("thin-web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.childOrCreate(segment)
	}
	root.handle = handle
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		// method 方法找不到
		return nil, false
	}
	// 根路由匹配
	if path == "/" {
		return &matchInfo{n: root}, true
	}
	segments := strings.Split(path[1:], "/")
	mi := &matchInfo{}
	for _, segment := range segments {
		var matchParam bool
		root, matchParam, ok = root.childOf(segment)
		if !ok {
			return nil, false
		}
		if matchParam {
			mi.addValue(root.path[1:], segment)
		}
	}
	mi.n = root
	return mi, true
}

type node struct {
	path     string
	handle   HandlerFunc
	children map[string]*node

	starChild  *node
	paramChild *node
}

func (n *node) childOrCreate(path string) *node {
	// 处理通配符路由
	if path == "*" {
		if n.starChild != nil {
			// todo 已经注册过该路由
		} else {
			n.starChild = &node{path: path}
		}
		return n.starChild
	}

	// 处理路径参数路由
	if path[0] == ':' {
		if n.paramChild != nil {
			// todo  已经注册过该参数路由
		} else {
			n.paramChild = &node{path: path}
		}
		return n.paramChild
	}

	// 处理静态路由
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	_, ok := n.children[path]
	if !ok {
		n.children[path] = &node{path: path}
	}
	return n.children[path]
}

func (n *node) childOf(path string) (*node, bool, bool) {
	// 匹配优先级 静态 > 参数 > 通配符
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, true
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
