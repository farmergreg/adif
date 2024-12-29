package adif

import (
	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

// Document represents a complete ADIF file.
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
	Name adifield.Field // Name is the field name
	Data string         // Data is the field value
}
