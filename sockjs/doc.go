/*
Package sockjs implements server side counterpart for the SockJS-client browser library.

See http://sockjs.org for more information about SockJS.

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
