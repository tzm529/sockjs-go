package sockjs

type pollingBaseSession struct {
	*baseSession
	out *queue
}

func newPollingBaseSession(pool *pool) (s *pollingBaseSession) {
	s = new(pollingBaseSession)
	s.baseSession = newBaseSession(pool)
	s.out = newQueue(true)
	return
}

func (s *pollingBaseSession) Receive() ([]byte, error) {
	return s.in.pull()
}

func (s *pollingBaseSession) Send(m []byte) error {
	s.out.push(m)
	return nil
}

func (s *pollingBaseSession) Close() error {
	s.out.close()
	s.closeBase()
	return nil
}
