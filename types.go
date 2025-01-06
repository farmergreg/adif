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
	Header Record `json:"header,omitempty"`

	// Records is a slice of Record.
	Records []Record `json:"records"`
}

// Record represents one ADIF record which may be a Header or a QSO.
type Record map[adifield.Field]string
