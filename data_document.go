package adif

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

type adifDocument struct {
	headerFields  []adifield.Field
	headerRecords []string

	fields  []adifield.Field
	records [][]string
}

type AdifDocument interface {
	GetHeaderPreamble() string
	SetHeaderPreamble(preamble string)

	GetHeader(key string) string
	SetHeader(key, value string)

	FieldCount() int
	GetField(index int, key string) string
	SetField(index int, key string) string
}

type AdifParser interface {
	ReadFrom(r io.Reader) (n int64, err error)
	GetDocument() *adifDocument
}

type adiParser struct {
	doc *adifDocument

	// bufValue is a reusable buffer used to temporarily store the VALUE of the current field.
	bufValue     []byte
	stringIntern map[string]string
}

func NewAdiParser() AdifParser {
	doc := &adifDocument{
		fields:  make([]adifield.Field, 0, 16),
		records: make([][]string, 0, 8),
	}

	return &adiParser{
		doc:          doc,
		stringIntern: make(map[string]string, 16),
	}
}

var (
	_ AdifParser    = (*adiParser)(nil) // Implements DataDocumentReader
	_ io.ReaderFrom = (*adiParser)(nil) // Implements io.ReaderFrom

	// _ io.WriterTo        = (*dataDocumentReader)(nil) // Implements io.WriterTo
	_ fmt.Stringer = (*adifDocument)(nil) // Implements fmt.Stringer
)

func (d *adiParser) GetDocument() *adifDocument {
	return d.doc
}

func (d *adiParser) ReadFrom(r io.Reader) (n int64, err error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	endOfRecord := true
	for {
		// Find the start of the next adi field
		c, err := discardUntilLessThan(br)
		n += c
		if err != nil {
			return n, err
		}

		if endOfRecord {
			alloc := 8
			if len(d.doc.records) > 1 {
				alloc = len(d.doc.records[len(d.doc.records)-1]) // use the current record's length as a hint for the next record's allocation
			}
			d.doc.records = append(d.doc.records, make([]string, 0, alloc))
		}

		endOfRecord, n, err = d.parseOneField(br)
		n += c
		if err != nil {
			return n, err
		}
	}
}

func (d *adiParser) parseOneField(br *bufio.Reader) (endOfRecord bool, n int64, err error) {
	// Step 1: Read in the entire data specifier "<fieldname:length:...>" and remove the trailing '>'
	volatileSpecifier, n, err := readDataSpecifierVolatile(br)
	if err != nil {
		return false, n, err
	}

	// Step 2: Parse Field Name
	volatileField, volatileLength, foundFirstColon := bytes.Cut(volatileSpecifier, []byte(":"))
	if len(volatileField) == 0 {
		return false, n, ErrMalformedADI // field name is empty
	}
	fastToUpper(volatileField)

	// Step 2.1: field name string interning - reduce memory allocations
	field := adifield.Field(unsafe.String(&volatileField[0], len(volatileField)))
	if fieldDef, ok := adifield.FieldMap[field]; ok {
		field = fieldDef.ID
	} else {
		field = adifield.Field(string(volatileField))
	}

	if field == adifield.EOR || field == adifield.EOH {
		if field == adifield.EOH {
			d.doc.headerFields = d.doc.fields
			d.doc.headerRecords = d.doc.records[0]

			d.doc.fields = make([]adifield.Field, 0, 16)
			d.doc.records = make([][]string, 0, 8)
		}
		return true, n, nil
	}

	// TODO: due to string interning, the below code could do pointer comparisons for more performance.

	// we assume the fields are in order, so we can use the last record to determine the current field index.
	fieldIndex := len(d.doc.records[len(d.doc.records)-1])
	if fieldIndex < len(d.doc.fields) && d.doc.fields[fieldIndex] == field {
		// the field is already in the list, so we can use the index.
	} else {
		// if the field at the index is not the same as the field we are parsing, then we need to find the index of the field.
		fieldIndex = -1
		for i, f := range d.doc.fields {
			if f == field {
				fieldIndex = i
				break
			}
		}

		if fieldIndex == -1 {
			d.doc.fields = append(d.doc.fields, field)
		}
	}

	// Step 3: Parse Field Length
	var length int
	if foundFirstColon {
		// look for the second colon which if present, indicates an adi optional type
		endIdx := bytes.IndexByte(volatileLength, ':')
		if endIdx == -1 {
			length, err = parseDataLength(volatileLength)
		} else {
			// we have a data type indicator.
			// this parser doesn't support it; ignore it
			length, err = parseDataLength(volatileLength[:endIdx])
		}
		if err != nil {
			// handle data length parsing errors
			return false, n, err
		}

		// Step 4: Read the field value (if any)
		// parseDataLength ensures that length is a reasonable value for us.
		// inlining v.s. a function call gains a tiny, but measurable amount of performance...
		if length > 0 {
			if cap(d.bufValue) < length {
				d.bufValue = make([]byte, length)
			}
			d.bufValue = d.bufValue[:length]

			var c int
			c, err = io.ReadFull(br, d.bufValue) // this will overwrite all of the 'volatile' variables (see above)
			n += int64(c)
			recordIndex := len(d.doc.records) - 1
			if err != nil {
				if err == io.EOF {
					return false, n, ErrMalformedADI
				}
				d.doc.records[recordIndex] = append(d.doc.records[recordIndex], string(d.bufValue[:c]))
				return false, n, err
			}
			d.doc.records[recordIndex] = append(d.doc.records[recordIndex], string(d.bufValue))
			return false, n, nil
		}
	}

	return false, n, nil
}

func (d *adiParser) internString(s string) string {
	if interned, ok := d.stringIntern[s]; ok {
		return interned
	}
	d.stringIntern[s] = s
	return s
}

// readDataSpecifierVolatile2 reads and returns the next data specifier as a byte slice, the number of bytes read, and any error encountered.
// The trailing '>' is removed from the returned byte slice.
//
// IMPORTANT:
// The returned byte slice is VOLATILE and will be invalidated by the next read from the underlying bufio.Reader.
//
// Per the Spec:
//
//	Data-Specifiers used to convey data in an ADI file are composed of a case-independent
//	field name F, a data length L, and an optional data type indicator T separated by colons and enclosed in angle brackets,
//	followed by data D of length L:
//
//	<F:L:T>
func readDataSpecifierVolatile(br *bufio.Reader) (volatileSpecifier []byte, n int64, err error) {
	// If ReadSlice returns bufio.ErrBufferFull, accumulator will contain ALL of the bytes read.
	// In most cases, accumulator will be null because we won't hit the bufio.ErrBufferFull condition.
	var accumulator []byte
	for {
		volatileSpecifier, err = br.ReadSlice('>')
		n += int64(len(volatileSpecifier))
		if err == nil {
			if accumulator != nil {
				// We've found '>' and have accumulated some bytes.
				// Update volatileSpecifier to point at the accumulated bytes before breaking out of the loop.
				volatileSpecifier = append(accumulator, volatileSpecifier...)
			}
			break
		}

		if err == bufio.ErrBufferFull {
			accumulator = append(accumulator, volatileSpecifier...)
			continue
		}

		if err == io.EOF {
			return volatileSpecifier, n, ErrMalformedADI
		}

		return volatileSpecifier, n, err
	}
	volatileSpecifier = volatileSpecifier[:len(volatileSpecifier)-1] // remove the trailing '>'
	return volatileSpecifier, n, nil
}

// discardUntilLessThan reads until it finds the '<' character, returning the number of bytes read
func discardUntilLessThan(br *bufio.Reader) (n int64, err error) {
	for {
		var b []byte
		b, err = br.ReadSlice('<')
		n += int64(len(b))

		switch err {
		case nil:
			return n, nil
		case bufio.ErrBufferFull:
			continue
		default:
			return n, err
		}
	}
}

/*
// fastToUpper is a faster version of bytes.ToUpper (we assume ASCII and that most characters are already UPPERCASE)
func fastToUpper(data []byte) {
	for i, b := range data {
		if b&0b01000000 > 0 && 'a' <= b && b <= 'z' {
			data[i] = b - 'a' + 'A'
		}
	}
}

// parseDataLength is an optimized replacement for strconv.Atoi.
func parseDataLength(data []byte) (value int, err error) {
	if len(data) == 0 {
		return 0, ErrInvalidFieldLength
	}

	for _, b := range data {
		if b < '0' || b > '9' {
			return 0, ErrInvalidFieldLength
		}

		// Parse digit, avoiding string allocations
		newVal := value*10 + int(b-'0')

		// Check for overflow or too big
		if newVal < value || newVal > maxADIReaderDataSize {
			return 0, ErrInvalidFieldLength
		}

		value = newVal
	}

	return value, nil
}

*/

func (d *adifDocument) String() string {
	var b strings.Builder

	// Preamble
	b.WriteString(adiHeaderPreamble)

	// Header (if exists)
	if d.headerFields != nil {
		for fieldIndex := range d.headerFields {
			b.WriteString("<")
			b.WriteString((string)(d.headerFields[fieldIndex]))
			b.WriteString(":")
			b.WriteString(strconv.Itoa(len(d.headerRecords[fieldIndex])))
			b.WriteString(">")
			b.WriteString(d.headerRecords[fieldIndex])
		}
		b.WriteString("<EOH>\n")
	}

	// QSOs
	for _, record := range d.records {
		for fieldIndex := range d.fields {
			b.WriteString("<")
			b.WriteString((string)(d.fields[fieldIndex]))
			b.WriteString(":")
			b.WriteString(strconv.Itoa(len(record[fieldIndex])))
			b.WriteString(">")
			b.WriteString(record[fieldIndex])
		}
		b.WriteString("<EOR>\n")
	}
	return b.String()
}
