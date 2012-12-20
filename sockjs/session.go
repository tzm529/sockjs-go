package sockjs

import (
	"errors"
)

var ErrSessionClosed error = errors.New("session closed")
var ErrSessionTimeout error = errors.New("session timeout")

type Session interface {
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
}
