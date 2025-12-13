package server

import "time"

// Option -.
type Option func(*Server)

// Timeout -.
//
//nolint:unused
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// ConnWaitTime -.
//
//nolint:unused
func ConnWaitTime(timeout time.Duration) Option {
	return func(s *Server) {
		s.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
//
//nolint:unused
func ConnAttempts(attempts int) Option {
	return func(s *Server) {
		s.conn.Attempts = attempts
	}
}
