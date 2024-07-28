//go:generate mockgen -package readers -destination readers_mock.go -source readers.go

// Package readers test utilities to simulate read and close errors
package readers

import "io"

type ReadCloserError interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type ReadCloserErrorWrapper struct {
	reader func(p []byte) (n int, err error)
	closer func() error
}

func NewReadCloserErrorWrapper(reader func(p []byte) (n int, err error), closer func() error) io.ReadCloser {
	return &ReadCloserErrorWrapper{
		reader: reader,
		closer: closer,
	}
}

func (r *ReadCloserErrorWrapper) Read(p []byte) (n int, err error) {
	return r.reader(p)
}

func (r *ReadCloserErrorWrapper) Close() error {
	return r.closer()
}
