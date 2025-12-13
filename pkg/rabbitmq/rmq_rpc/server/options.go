package server

import "time"

// Option -.
type Option func(*Server)

// Timeout -.
//
//nolint:unused // Этот метод используется через рефлексию в другом пакете
func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// ConnWaitTime -.
//
//nolint:unused // Этот метод используется через рефлексию в другом пакете
func ConnWaitTime(timeout time.Duration) Option {
	return func(s *Server) {
		s.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
//
//nolint:unused // Этот метод используется через рефлексию в другом пакете
func ConnAttempts(attempts int) Option {
	return func(s *Server) {
		s.conn.Attempts = attempts
	}
}
