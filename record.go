package adif

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

// Interface implementations
var (
	_ io.WriterTo   = &Record{}
	_ io.ReaderFrom = &Record{}
	_ fmt.Stringer  = &Record{}
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 1024)
		return &b
	},
}

// NewRecord creates a new Record with the default initial capacity.
func NewRecord() Record {
	return NewRecordWithCapacity(-1)
}

// NewRecordWithCapacity creates a new Record with a specific initial capacity.
func NewRecordWithCapacity(initialCapacity int) Record {
	if initialCapacity < 1 {
		initialCapacity = 7
	}
	return make(Record, initialCapacity)
}

// Reset clears the record of all fields.
func (r Record) Reset() {
	clear(r)
}

// ReadFrom reads an ADIF formatted record from the provided io.Reader.
// It returns the number of bytes read and any error encountered.
// io.EOF is returned when no more records are available.
// Existing fields will be updated and new fields added as necessary.
// Header records are SKIPPED.
//
// n.b. This method is best used to read a single record.
// When reading multiple records, create an ADIFReader using NewADIReader().
// Use its Next() method to obtain maximum speed and memory efficiency while iterating over records.
func (r *Record) ReadFrom(src io.Reader) (int64, error) {
	p := NewADIReader(src, true)
	var n int64

	record, _, c, err := p.Next()
	n += c
	if err != nil {
		return n, err
	}

	*r = record
	return n, nil
}

// WriteTo writes ADI formatted record data to the provided io.Writer.
// It WILL NOT write the <EOR> / <EOH> tag.
// It returns the number of bytes written and any error encountered.
func (r *Record) WriteTo(dest io.Writer) (int64, error) {
	adiLength := r.appendAsADIPreCalculate()
	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr

	if cap(buf) < adiLength {
		buf = make([]byte, 0, adiLength)
		*bufPtr = buf
	}
	buf = buf[:0]

	buf = r.appendAsADI(buf)
	n, err := dest.Write(buf)

	bufferPool.Put(bufPtr)
	return int64(n), err
}

// appendAsADI writes the ADI formatted QSO record to the provided buffer.
// The buffer should have sufficient capacity to avoid reallocations.
// You should use appendAsADIPreCalculate() to determine the required buffer capacity.
// Field order is NOT guaranteed to be stable.
func (r *Record) appendAsADI(buf []byte) []byte {
	if len(*r) == 0 {
		return buf
	}

	// Priority fields first.
	// This may change.
	// These fields and their specific order are NOT guaranteed to be in the same position in future versions of this library.
	buf = appendField(buf, adifield.CALL, (*r)[adifield.CALL])
	buf = appendField(buf, adifield.BAND, (*r)[adifield.BAND])
	buf = appendField(buf, adifield.MODE, (*r)[adifield.MODE])
	buf = appendField(buf, adifield.QSO_DATE, (*r)[adifield.QSO_DATE])
	buf = appendField(buf, adifield.TIME_ON, (*r)[adifield.TIME_ON])
	buf = appendField(buf, adifield.QSO_DATE_OFF, (*r)[adifield.QSO_DATE_OFF])
	buf = appendField(buf, adifield.TIME_OFF, (*r)[adifield.TIME_OFF])
	buf = appendField(buf, adifield.PROP_MODE, (*r)[adifield.PROP_MODE])
	buf = appendField(buf, adifield.SAT_NAME, (*r)[adifield.SAT_NAME])
	buf = appendField(buf, adifield.OPERATOR, (*r)[adifield.OPERATOR])
	buf = appendField(buf, adifield.STATION_CALLSIGN, (*r)[adifield.STATION_CALLSIGN])
	buf = appendField(buf, adifield.QTH, (*r)[adifield.QTH])
	buf = appendField(buf, adifield.GRIDSQUARE, (*r)[adifield.GRIDSQUARE])

	// Remaining fields
	for field, value := range *r {

		// Skip fields we've already handled
		switch field {
		case adifield.CALL, adifield.BAND, adifield.MODE,
			adifield.QSO_DATE, adifield.TIME_ON,
			adifield.QSO_DATE_OFF, adifield.TIME_OFF,
			adifield.PROP_MODE, adifield.SAT_NAME:
			continue
		}
		buf = appendField(buf, field, value)
	}

	return buf
}

// appendField adds a single ADIF field to the buffer
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

// appendAsADIPreCalculate returns the length of the record in bytes when exported to ADI format by the AppendAsADI method.
func (r *Record) appendAsADIPreCalculate() (adiLength int) {
	if len(*r) == 0 {
		return 0
	}

	for field, value := range *r {
		valueLength := len(value)
		if valueLength == 0 {
			continue
		}
		adiLength += 3 + valueLength + len(field) // 3 for '<', ':', '>'

		// Avoid strconv.Itoa string allocation by calculating number of base 10 digits mathematically
		switch {
		case valueLength < 10:
			adiLength += 1
		case valueLength < 100:
			adiLength += 2
		case valueLength < 1_000:
			adiLength += 3
		default:
			adiLength += int(math.Log10(float64(valueLength))) + 1
		}
	}

	return adiLength
}

// Clean
// trims whitespace in the field values
func (r *Record) Clean() {
	for field, value := range *r {
		trimmed := strings.TrimSpace(value)
		(*r)[field] = trimmed
	}
}

// String returns the ADIF record as a string.
// The field order is NOT guaranteed to be stable.
func (r *Record) String() string {
	adiLength := r.appendAsADIPreCalculate()
	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr

	if cap(buf) < adiLength {
		buf = make([]byte, 0, adiLength)
		*bufPtr = buf
	}
	buf = buf[:0]

	buf = r.appendAsADI(buf)
	s := string(buf)

	bufferPool.Put(bufPtr)
	return s
}
