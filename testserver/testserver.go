package main

import (
	"github.com/fzzbt/sockjs-go/sockjs"
	"fmt"
	"net/http"
)

func echoOpen(c *sockjs.Conn) {
	fmt.Println("connection created")	
}

func echoMessage(c *sockjs.Conn, message string) {
	println(message)
	c.Send(message)
}

func echoClose(c *sockjs.Conn) {
	fmt.Println("connection closed")	
}

func main() {
	http.Handle("/static", http.FileServer(http.Dir("./static")))
	mux := sockjs.NewServeMux(http.DefaultServeMux)
	echoHandler := sockjs.NewHandler(sockjs.Config{
		SockjsUrl:     "http://cdn.sockjs.org/sockjs-0.3.4.min.js",
		Prefix: "/echo",
		Websocket:     true,
		ResponseLimit: 4096,
	})
	echoHandler.OnOpen = echoOpen
	echoHandler.OnMessage = echoMessage
	echoHandler.OnClose = echoClose

	dwsechoHandler := sockjs.NewHandler(sockjs.Config{
		SockjsUrl:     "http://cdn.sockjs.org/sockjs-0.3.4.min.js",
		Prefix: "/disabled_websocket_echo",
		Websocket:     false,
		ResponseLimit: 4096,
	})
	dwsechoHandler.OnOpen = echoOpen
	dwsechoHandler.OnMessage = echoMessage
	dwsechoHandler.OnClose = echoClose

	closeHandler := sockjs.NewHandler(sockjs.Config{
		SockjsUrl:     "http://cdn.sockjs.org/sockjs-0.3.4.min.js",
		Prefix: "/close",
		Websocket:     true,
		ResponseLimit: 4096,
	})
	closeHandler.OnOpen = func(c *sockjs.Conn) { c.Close() }

	mux.Handle(echoHandler)
	mux.Handle(dwsechoHandler)
	mux.Handle(closeHandler)

	err := http.ListenAndServe(":8081", mux)
	if err != nil { fmt.Println(err) }
}
