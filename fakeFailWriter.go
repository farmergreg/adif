package adif

import "errors"

// fakeFailWriter is a writer that fails after writing a certain number of bytes
type fakeFailWriter struct {
	maxBytes int
	written  int
}

func (fw *fakeFailWriter) Write(p []byte) (n int, err error) {
	if fw.written >= fw.maxBytes {
		return 0, errors.New("write failed: max bytes exceeded")
	}

	remaining := fw.maxBytes - fw.written
	if len(p) <= remaining {
		fw.written += len(p)
		return len(p), nil
	}

	// Partial write before failure
	fw.written = fw.maxBytes
	return remaining, errors.New("write failed: max bytes exceeded")
}
