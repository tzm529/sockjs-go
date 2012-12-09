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
	c := sockjs.NewConfig()
	c.Prefix = "/echo"
	echoHandler := sockjs.NewHandler(c)
	echoHandler.OnOpen = echoOpen
	echoHandler.OnMessage = echoMessage
	echoHandler.OnClose = echoClose

	c = sockjs.NewConfig()
	c.Prefix = "/disabled_websocket_echo"
	c.Websocket = false
	dwsechoHandler := sockjs.NewHandler(c)
	dwsechoHandler.OnOpen = echoOpen
	dwsechoHandler.OnMessage = echoMessage
	dwsechoHandler.OnClose = echoClose

	c = sockjs.NewConfig()
	c.Prefix = "/close"
	closeHandler := sockjs.NewHandler(c)
	closeHandler.OnOpen = func(c *sockjs.Conn) { c.Close() }

	mux.Handle(echoHandler)
	mux.Handle(dwsechoHandler)
	mux.Handle(closeHandler)

	err := http.ListenAndServe(":8081", mux)
	if err != nil { fmt.Println(err) }
}
