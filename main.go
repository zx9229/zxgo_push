package main

import (
	businessservice "github.com/zx9229/zxgo_push/BusinessService"
	simplehttpserver "github.com/zx9229/zxgo_push/SimpleHttpServer"
	"golang.org/x/net/websocket"
)

func main() {

	shs := simplehttpserver.New_SimpleHttpServer("localhost:8080")
	bs := businessservice.New_BusinessService()
	shs.GetHttpServeMux().Handle("/websocket", websocket.Handler(bs.GetConnectionManager().HandleWebsocket))
	shs.Run()
}
