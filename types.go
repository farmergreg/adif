package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

// Document represents a complete ADIF document.
type Document struct {
	Header  *Record
	Records []Record
}

// Record represents one ADIF record which may be a Header or a QSO.
type Record struct {
	Fields []FieldEntry
}

// FieldEntry represents an ADIF field and its data.
// It is designed to ensure cpu cache locality during field lookup and value retrieval.
type FieldEntry struct {

	// Name is the field name.
	// Unlike the ADIF specification, the field name MUST be in UPPERCASE for use in this library.
	// The UPPERCASE only rule allows for faster lookup and retrieval of field values.
	Name adifield.Field

	// Data is the field value
	Data string
}
