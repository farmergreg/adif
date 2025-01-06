package adif

import (
	"fmt"
	"io"
	"math"
	"slices"
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

var fieldPool = sync.Pool{
	New: func() interface{} {
		s := make([]adifield.Field, 0, 64)
		return &s
	},
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 1024)
		return &b
	},
}

// NewRecord creates a new Record with the default initial capacity.
func NewRecord() Record {
	return make(Record)
}

// NewRecordWithCapacity creates a new Record with a specific initial capacity.
func NewRecordWithCapacity(initialCapacity int) Record {
	if initialCapacity < 0 {
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
// Existing fields will be updated and add new fields added when necessary.
// Header records are SKIPPED.
//
// n.b. This method is best used when reading a single record.
// When reading multiple records, create an ADIFReader using NewadiReader().
// Use its Next() method to obtain maximum speed and memory efficiency.
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
// It WILL NOT write the <EOR>/<EOH> tag.
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
// You should use appendAsADIPreCalculate() to determine the required capacity
func (r *Record) appendAsADI(buf []byte) []byte {
	if len(*r) == 0 {
		return buf
	}

	fieldsPtr := fieldPool.Get().(*[]adifield.Field)
	fields := (*fieldsPtr)[:0]

	for field := range *r {
		fields = append(fields, field)
	}
	slices.Sort(fields)

	for _, field := range fields {
		value := (*r)[field]
		if value == "" {
			continue
		}

		buf = append(buf, '<')
		buf = append(buf, field...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(len(value)), 10)
		buf = append(buf, '>')
		buf = append(buf, value...)
	}

	// Return slice to pool
	fieldPool.Put(fieldsPtr)
	return buf
}

// appendAsADIPreCalculate returns:
// 1) the length of the record in bytes when exported to ADI format by the AppendAsADI method.
// 2) a boolean indicating if the record is a header record.
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
// 1) trims whitespace in the field values
// 2) deletes fields with empty values
func (r *Record) Clean() {
	for field, value := range *r {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			delete(*r, field)
		} else {
			(*r)[field] = trimmed
		}
	}
}

// String returns the ADIF record as a string
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
