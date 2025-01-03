package adif

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

var (
	_ io.WriterTo   = &Document{}
	_ io.ReaderFrom = &Document{}
	_ fmt.Stringer  = &Document{}
)

// Reset clears the document and prepares it for reuse.
func (f *Document) Reset() {
	f.Header = nil
	f.Records = make([]Record, 0, 64)
}

// WriteTo writes the document in ADI format to the given writer.
func (f *Document) WriteTo(w io.Writer) (n int64, err error) {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}

	if f.Header != nil {
		c, err := f.Header.WriteTo(bw)
		n += c
		if err != nil {
			return handleFlush(bw, n, err)
		}
		cc, err := bw.Write([]byte{'\n'})
		n += int64(cc)
		if err != nil {
			return handleFlush(bw, n, err)
		}
	}

	for _, record := range f.Records {
		c, err := record.WriteTo(bw)
		n += c
		if err != nil {
			return handleFlush(bw, n, err)
		}

		cc, err := bw.Write([]byte{'\n'})
		n += int64(cc)
		if err != nil {
			return handleFlush(bw, n, err)
		}
	}

	return handleFlush(bw, n, err)
}

func handleFlush(bw *bufio.Writer, n int64, err error) (int64, error) {
	if bwerr := bw.Flush(); bwerr != nil {
		return n, fmt.Errorf("error flushing writer: %w", bwerr)
	}
	return n, err
}

// ReadFrom reads an ADI document from the given reader.
// You should call Reset() if you do not want add/update the existing document.
func (f *Document) ReadFrom(r io.Reader) (n int64, err error) {
	p := NewADIParser(r, false)

	firstRecord, n, err := p.Parse()
	if err != nil {
		return n, err
	}

	if isHeader, _ := firstRecord.isHeaderRecord(); isHeader {
		f.Header = firstRecord
	} else {
		f.Records = append(f.Records, *firstRecord)
	}

	for {
		record, c, err := p.Parse()
		n += c
		if err != nil {
			if err == io.EOF {
				break
			}
			return n, err
		}
		f.Records = append(f.Records, *record)

		// prevent memory exhaustion attacks
		if n > DocumentMaxSizeInBytes {
			return n, ErrDocumentTooLarge
		}
	}

	return n, nil
}

// String returns the document as an ADI string.
// Returns an empty string if the receiver is nil.
func (f *Document) String() string {
	if f == nil || (len(f.Records) == 0 && f.Header == nil) {
		return ""
	}

	sb := strings.Builder{}

	_, err := f.WriteTo(&sb)
	if err != nil {
		return fmt.Sprintf("error while building adi string: %v", err)
	}

	return sb.String()
}
