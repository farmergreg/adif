package adif

import (
	"errors"
	"io"
)

var errMockWrite = errors.New("mock write error")

var _ = (io.Writer)(&mockAlwaysErrorWriter{})

// mockAlwaysErrorWriter is a writer that always returns an error
type mockAlwaysErrorWriter struct{}

func (mw *mockAlwaysErrorWriter) Write(p []byte) (n int, err error) {
	return 0, errMockWrite
}
