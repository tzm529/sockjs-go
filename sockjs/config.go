package sockjs

import (
	"crypto/md5"
	"fmt"
)

type Config struct {
	SockjsURL       string
	ResponseLimit   int
	Websocket       bool
	HeartbeatDelay  int
	DisconnectDelay int

	// Iframe page
	iframePage []byte
	iframeHash string
}

func NewConfig() (c Config) {
	c.SockjsURL = "http://cdn.sockjs.org/sockjs-0.3.4.min.js"
	c.Websocket = true
	c.ResponseLimit = 4096

	c.iframePage = []byte(fmt.Sprintf(iframePageFormat, c.SockjsURL))
	hash := md5.New()
	hash.Write(c.iframePage)
	c.iframeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return
}
