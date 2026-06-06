package adif

import (
	"io"
	"strings"

	"github.com/farmergreg/spec/v6/adifield"
)

// Record is a map of ADIF fields to their values, representing either a header or QSO record.
// Field keys are normalized to uppercase by the Scanner.
// Use adifield constants or adifield.New() to construct field names.
//
// Example:
//
//	r := adif.NewRecord()
//	r[adifield.CALL] = "K9CTS"
//	r[adifield.BAND] = "20m"
type Record map[adifield.Field]string

// NewRecord returns a new empty Record with a sensible default capacity.
func NewRecord() Record {
	return make(Record, 7)
}

// String returns the record's fields serialized in ADI format, without an EOR or EOH tag.
// Priority fields are written first in a fixed order; remaining fields follow in map iteration order.
// Implements fmt.Stringer.
func (r Record) String() string {
	var sb strings.Builder
	r.WriteTo(&sb) //nolint:errcheck — strings.Builder.Write never returns an error
	return sb.String()
}

// WriteTo writes the record's fields in ADI format to w, without an EOR or EOH tag.
// Priority fields are written first in a fixed order; remaining fields follow in map iteration order.
// Implements io.WriterTo.
func (r Record) WriteTo(w io.Writer) (int64, error) {
	// Use the shared writer buffer pool to avoid a per-call allocation and the cost of pre-sizing the buffer.
	bufPtr := writerBufPool.Get().(*[]byte)
	buf := appendFieldsADI(r, (*bufPtr)[:0])
	n, err := w.Write(buf)
	*bufPtr = buf
	writerBufPool.Put(bufPtr)
	return int64(n), err
}
