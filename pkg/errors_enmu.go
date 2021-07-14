package pkg

import "errors"

var (
	ErrNonPublic        = errors.New("Registered non-public service")
	ErrNoAvailable      = errors.New("No service is available, or provide service is not open")
	ErrCrc32            = errors.New("checksumIEEE error")
	ErrSerialization404 = errors.New("serialization 404")
	ErrCompressor404    = errors.New("compressor 404")
	ErrTimeout          = errors.New("time out")
	ErrCircuitBreaker   = errors.New("circuit breaker")
)
