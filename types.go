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
	Header *Record

	// Records is a slice of Record.
	Records []Record
}

// Record represents one ADIF record which may be a Header or a QSO.
type Record struct {
	// Fields is a slice of Field.
	Fields []Field
}

// Field represents an ADIF field and its data.
// It is designed to ensure cpu cache locality during field lookup and value retrieval.
type Field struct {

	// Name is the field name.
	// Unlike the ADIF specification, the field name MUST be in UPPERCASE for use in this library.
	// The UPPERCASE only rule allows for faster lookup and retrieval of field values.
	Name adifield.Field

	// Data is the field value
	Data string
}
