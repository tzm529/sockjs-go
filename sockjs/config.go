package sockjs

type Config struct {
	SockjsUrl       string
	Prefix          string
	ResponseLimit   int
	Websocket       bool
	HeartbeatDelay  int
	DisconnectDelay int
}
