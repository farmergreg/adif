package adif

import (
	"errors"
	"io"
)

var _ = (io.Reader)(&mockFailReader{})

type mockFailReader struct {
	backingData []byte
	pos         int
	maxBytes    int
}

func (fr *mockFailReader) Read(p []byte) (n int, err error) {
	if fr.pos >= fr.maxBytes {
		return 0, errors.New("read failed: max bytes exceeded")
	}

	bytesRead := 0
	for i := range p {
		p[i] = fr.backingData[fr.pos]
		fr.pos++
		bytesRead++
		if fr.pos >= fr.maxBytes {
			return bytesRead, errors.New("read failed: max bytes exceeded")
		}
	}

	return bytesRead, nil
}
