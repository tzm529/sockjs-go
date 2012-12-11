package sockjs

type Session interface{
	Receive() ([]byte, error)
	Send([]byte) error
	Close() error
}
