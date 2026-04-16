package main

import (
	"net/http"
	"time"
)

const (
	ServerReadTimeout = 100 * time.Minute
	ServerWriteTimeout = 100 * time.Minute
	ServerDefaultAddress = ":8080"
)

type Option func (*http.Server) error

func WithAddr(addr string) Option {
	return func(s *http.Server) error {
		s.Addr = addr
		return nil
	}
}

func WithReadTimeout(timeout time.Duration) Option {
	return func(s *http.Server) error {
		s.ReadHeaderTimeout = timeout
		return nil
	}
}

func NewServer(opts ...Option) (*http.Server, error) {
	// Default values
	s := &http.Server{
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		Addr:         ServerDefaultAddress,
	}
	// Apply options
	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}