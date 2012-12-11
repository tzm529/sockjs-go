package sockjs

import (
	"crypto/md5"
	"fmt"
	"time"
)

type Config struct {
	SockjsURL       string
	ResponseLimit   int
	Websocket       bool
	HeartbeatDelay  time.Duration
	DisconnectDelay time.Duration

	// Iframe page
	iframePage []byte
	iframeHash string
}

func NewConfig() (c Config) {
	c.SockjsURL = "http://cdn.sockjs.org/sockjs-0.3.4.min.js"
	c.ResponseLimit = 131072 // 128K
	c.Websocket = true
	c.HeartbeatDelay = time.Duration(25)*time.Second
	c.DisconnectDelay = time.Duration(5)*time.Second

	c.iframePage = []byte(fmt.Sprintf(iframePageFormat, c.SockjsURL))
	hash := md5.New()
	hash.Write(c.iframePage)
	c.iframeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return
}
