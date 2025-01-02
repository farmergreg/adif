package adif

import (
	"bufio"
	"bytes"
	"io"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

var _ ADIFParser = (*adiParser)(nil) // Implements ADIFParser

const (
	// 1MB - this is the maximum size of a field value that we will accept.
	// This is intended to be a generous limit for most applications while providing protection against malformed and/or malicious input.
	//
	// The data is part of the ADIF "Data-Specifier."
	// Per the ADIF spec:
	//   ADI-exporting applications can place as much data in a Data-Specifier as they choose.
	//   ADI-importing applications can import as much data from a Data-Specifier as they choose.
	maxADIParserDataSize = 1024 * 1024 * 1
)

// adiParser is a high-performance ADIF parser that can parse ADIF *.adi formatted records.
type adiParser struct {
	r        *bufio.Reader
	bufValue []byte

	// we assume most records will have similar field counts and allocate accordingly.
	preAllocateFields int

	// skipHeader is true if the header record should be skipped.
	skipHeader bool
}

// NewADIParser returns an ADIFParser that can parse ADIF *.adi formatted records.
// If skipHeader is true, Parse() will not return the header record if it exists.
func NewADIParser(r io.Reader, skipHeader bool) ADIFParser {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	return &adiParser{
		r:                 br,
		skipHeader:        skipHeader,
		preAllocateFields: 8,
	}
}

// Parse reads and returns the next Record.
// It returns io.EOF when no more records are available.
func (p *adiParser) Parse() (*Record, int64, error) {
	result := NewRecord(p.preAllocateFields)
	var n int64
	for {
		// Find the start of the next field for parsing
		c, err := p.discardUntilLessThan()
		n += c

		if err != nil {
			if err == io.EOF && len(result.Fields) > 0 {
				// We have a valid record, return it without the EOF error
				// The next call to Parse() will return io.EOF
				return result, n, nil
			}
			return result, n, err
		}

		field, value, c, err := p.parseOneField()
		n += c
		if err != nil {
			return result, n, err
		}

		switch field {
		case adifield.EOH:
			if !p.skipHeader && len(result.Fields) > 0 {
				return result, n, nil
			}
			// the record is not empty, reset it to prepare for the next record
			result = NewRecord(p.preAllocateFields)
			continue
		case adifield.EOR:
			if len(result.Fields) > 0 {
				if len(result.Fields) > p.preAllocateFields {
					p.preAllocateFields = len(result.Fields)
				}
				return result, n, nil
			}
			// we know record is empty... no need to reset it
			continue
		}

		// n.b. if a duplicate field is found, it will replace the previous value
		result.Set(field, value)
	}
}

// parseOneField reads the next field definition and returns the field name, value, and the number of bytes read.
//
// It is heavily optimized for speed and memory use.
// Currently, It can double the speed of go's stdlib JSON marshaling for similar data.
//
// Future Plans: I would like to take a look at using simd directly.
// However, the current implementation IS attempting to take advantage of the standard library's existing simd capabilities.
func (p *adiParser) parseOneField() (field adifield.Field, value string, n int64, err error) {
	// Step 1: Read in the entire data specifier "<fieldname:length:...>" and remove the trailing '>'
	volatileSpecifier, n, err := p.readDataSpecifierVolatile()
	if err != nil {
		return "", "", n, err
	}

	// Step 2: Parse Field Name
	volatileField, volatileLength, foundFirstColon := bytes.Cut(volatileSpecifier, []byte(":"))
	if len(volatileField) == 0 {
		return "", "", n, ErrMalformedADI // field name is empty
	}
	fastToUpper(volatileField)
	field = adifield.Field(string(volatileField))

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
			return field, "", n, err
		}

		// Step 4: Read the field value (if any)
		// inlining v.s. a function call gains a tiny amount of performance...
		if length > 0 {
			if cap(p.bufValue) < length {
				p.bufValue = make([]byte, length)
			}
			p.bufValue = p.bufValue[:length]

			var c int
			c, err = io.ReadFull(p.r, p.bufValue) // this will overwrite volatileSpecifier (see above)
			n += int64(c)
			if err != nil {
				if err == io.EOF {
					return field, string(p.bufValue[:c]), n, ErrMalformedADI
				}
				return field, string(p.bufValue[:c]), n, err
			}
			return field, string(p.bufValue), n, nil
		}
	}

	return field, "", n, nil
}

// readDataSpecifierVolatile reads and returns the next data specifier as a byte slice, the number of bytes read, and any error encountered.
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
func (p *adiParser) readDataSpecifierVolatile() (volatileSpecifier []byte, n int64, err error) {
	// If ReadSlice returns bufio.ErrBufferFull this accumulates ALL of the bytes read.
	// In most cases, this will be null because we won't hit the bufio.ErrBufferFull condition.
	var accumulator []byte
	for {
		volatileSpecifier, err = p.r.ReadSlice('>')
		n += int64(len(volatileSpecifier))
		if err == nil {
			if accumulator != nil {
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
func (p *adiParser) discardUntilLessThan() (n int64, err error) {
	for {
		var b []byte
		b, err = p.r.ReadSlice('<')
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
		if newVal < value || newVal > maxADIParserDataSize {
			return 0, ErrInvalidFieldLength
		}

		value = newVal
	}

	return value, nil
}
