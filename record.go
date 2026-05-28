package adif

import (
	"io"
	"math"
	"strconv"
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
	size := recordSizeADI(r) - 6 // recordSizeADI includes the EOR/EOH tag; we omit it here
	buf := make([]byte, 0, size)
	buf = appendFieldsADI(r, buf)
	n, err := w.Write(buf)
	return int64(n), err
}

// appendField appends a single ADIF field in ADI format to buf.
// Returns buf unchanged when value is empty.
func appendField(buf []byte, field adifield.Field, value string) []byte {
	if value == "" {
		return buf
	}
	buf = append(buf, '<')
	buf = append(buf, field...)
	buf = append(buf, ':')
	buf = strconv.AppendInt(buf, int64(len(value)), 10)
	buf = append(buf, '>')
	buf = append(buf, value...)
	return buf
}

// digitCount returns the number of base-10 digits needed to represent n.
func digitCount(n int) int {
	switch {
	case n < 10:
		return 1
	case n < 100:
		return 2
	case n < 1_000:
		return 3
	case n < 10_000:
		return 4
	case n < 100_000:
		return 5
	default:
		return int(math.Log10(float64(n))) + 1
	}
}
