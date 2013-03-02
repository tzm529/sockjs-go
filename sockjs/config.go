package sockjs

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	// URL for SockJS client library.
	// Default: "https://cdn.sockjs.org/sockjs-0.3.4.min.js"
	SockjsURL string

	// Enables websocket transport.
	// Default: true
	Websocket bool

	// Byte limit that can be sent over streaming session before it's closed.
	// Default: 128 * 1024
	ResponseLimit int

	// Adds JSESSIONID cookie to requests to enable sticky sessions.
	// Dummy value is used unless JsessionFunc is set.
	// Default: false.
	Jsessionid bool

	// Function that is called to set the JSESSIONID cookie, if Jsessionid setting is set.
	JsessionidFunc func(http.ResponseWriter, *http.Request)

	// Enables IP-address checks for legacy transports. 
	// If enabled, all subsequent calls must be from the same IP-address.
	// Default: true
	VerifyAddr bool

	// Heartbeat delay.
	// Default: 25 seconds
	HeartbeatDelay time.Duration

	// Disconnect delay.
	// Default: 5 seconds
	DisconnectDelay time.Duration

	// List of headers that are copied from incoming requests to SessionInfo.
	// Default: []string{"referer", "x-client-ip", "x-forwarded-for",
	//                   "x-cluster-client-ip", "via", "x-real-ip", "host"}
	Headers []string

	// Logger used for logging various information such as errors.
	// Default: log.New(os.Stdout, "sockjs", 0)
	Logger *log.Logger

	iframePage []byte
	iframeHash string
}

// NewConfig returns a new Config with the default settings.
func NewConfig() (c Config) {
	c.SockjsURL = "https://cdn.sockjs.org/sockjs-0.3.4.min.js"
	c.Websocket = true
	c.ResponseLimit = 128 * 1024
	c.VerifyAddr = true
	c.HeartbeatDelay = time.Duration(25) * time.Second
	c.DisconnectDelay = time.Duration(5) * time.Second
	c.Headers = []string{"referer", "x-client-ip", "x-forwarded-for",
		"x-cluster-client-ip", "via", "x-real-ip", "host"}
	c.Logger = log.New(os.Stdout, "sockjs", 0)
	c.iframePage = []byte(fmt.Sprintf(iframePageFormat, c.SockjsURL))
	hash := md5.New()
	hash.Write(c.iframePage)
	c.iframeHash = fmt.Sprintf("%x", hash.Sum(nil))
	return
}
