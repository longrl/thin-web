package web

import "testing"

func TestServer(t *testing.T) {
	server := &HttpServer{}
	server.Start(":8080")
}
