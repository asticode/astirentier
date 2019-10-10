package main

type session struct{}

func newSession() *session {
	return &session{}
}

func (s *session) close() (err error) {
	return
}
