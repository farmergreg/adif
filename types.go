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

	// headerPreamble is the preamble that will be written when the document is written to an io.Writer.
	// it has NOTHING to do with ADIFReader/ADIReader.
	headerPreamble string
}

// Record represents one ADIF record which may be a Header or a QSO.
// n.b. Some software, like Log4OM like to incorrectly place header only fields (i.e. PROGRAMID) into QSO records...
type Record struct {
	// Fields is a slice of FieldEntry.
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
