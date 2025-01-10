package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

// Document represents a complete ADIF document.
//
// Future Work:
// This type intentionally resembles the ADX XML structure even though XML is not currently supported by this library.
type Document struct {
	// Header is nil if there is no header.
	// Otherwise it will be a Record with header fields inside.
	Header *Record `json:"header,omitempty"`

	// Records is a slice of Record.
	// It contains the QSO records.
	Records []*Record
}

// Record represents either a header or QSO record.
type Record struct {
	fields map[adifield.Field]string
}
