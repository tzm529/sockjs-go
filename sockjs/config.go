package sockjs

import (
	"crypto/md5"
	"fmt"
	"time"
)

type Config struct {
	// URL for SockJS client library.
	// Default: "http://cdn.sockjs.org/sockjs-0.3.4.min.js"
	SockjsURL       string

	// Enables websocket transport.
	// Default: true
	Websocket       bool

	// Byte limit that can be sent over streaming session before it's closed.
	// Default: 131072
	ResponseLimit   int

	// Enables sticky sessions. 
	// Default: false.
	Jsessionid bool

	// Enables IP-address checks for polling transports. 
	// If enabled, all subsequent polling calls must be from the same IP-address.
	// Default: true
	VerifyAddr bool

	// Heartbeat delay.
	// Default: 25 seconds
	HeartbeatDelay  time.Duration

	// Disconnection delay.
	// Default: 5 seconds
	DisconnectDelay time.Duration

	// List of headers that are copied from incoming requests to SessionInfo.
	// Default: []string{"referer", "x-client-ip", "x-forwarded-for",
	//                   "x-cluster-client-ip", "via", "x-real-ip", "host"}
	Headers []string

	iframePage []byte
	iframeHash string
}

func NewConfig() (c Config) {
	c.SockjsURL = "http://cdn.sockjs.org/sockjs-0.3.4.min.js"
	c.Websocket = true
	c.ResponseLimit = 131072 // 128 KiB
	c.Jsessionid = false
	c.VerifyAddr = true
	c.HeartbeatDelay = time.Duration(25) * time.Second
	c.DisconnectDelay = time.Duration(5) * time.Second
	c.Headers = []string{"referer", "x-client-ip", "x-forwarded-for",
	"x-cluster-client-ip", "via", "x-real-ip", "host"}

	c.iframePage = []byte(fmt.Sprintf(iframePageFormat, c.SockjsURL))
	hash := md5.New()
	hash.Write(c.iframePage)
	c.iframeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return
}
