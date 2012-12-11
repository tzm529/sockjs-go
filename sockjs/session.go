package sockjs

import (
	"errors"
)

var ErrSessionClosed error = errors.New("session is closed")

type Session interface{
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
}
