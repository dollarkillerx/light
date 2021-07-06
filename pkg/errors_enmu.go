package pkg

import "errors"

var (
	ErrNonPublic   = errors.New("Registered non-public service")
	ErrNoAvailable = errors.New("No service is available, or provide service is not open")
	ErrCrc32       = errors.New("checksumIEEE error")
)
