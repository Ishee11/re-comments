package client

import "time"

// Option -.
type Option func(*Client)

// Timeout -.
//
//nolint:unused
func Timeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// ConnWaitTime -.
//
//nolint:unused
func ConnWaitTime(timeout time.Duration) Option {
	return func(c *Client) {
		c.conn.WaitTime = timeout
	}
}

// ConnAttempts -.
//
//nolint:unused
func ConnAttempts(attempts int) Option {
	return func(c *Client) {
		c.conn.Attempts = attempts
	}
}
