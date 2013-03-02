/*
Package sockjs implements a SockJS server.

SockJS is a JavaScript library (for browsers) that provides a WebSocket-like object.
SockJS gives you a coherent, cross-browser, Javascript API which creates a low latency,
full duplex, cross-domain communication channel between the browser and the web server,
with WebSockets or without. 
Under the hood SockJS tries to use native WebSockets first. 
If that fails it can use a variety of browser-specific transport protocols and presents them
through WebSocket-like abstractions. 
SockJS is intended to work for all modern browsers and in 
environments which don't support WebSocket protcol, for example behind restrictive corporate
proxies. 
See http://sockjs.org for more info about SockJS.

Example echo server:

	func echoHandler(s sockjs.Session) {
		for {
			m := s.Receive()
			if m == nil {
				break
			}
			s.Send(m)
		}
	}

	func main() {
		server := sockjs.NewServer(http.DefaultServeMux)
		server.Handle("/echo", echoHandler, sockjs.NewConfig())
		err := http.ListenAndServe(":8081", server)
		if err != nil {
			fmt.Println(err)
		}
	}

*/
package sockjs
