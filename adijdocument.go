package adif

import "github.com/hamradiolog-net/adif-spec/v6/adifield"

// ADIJDocument represents an ADIJ (ADIF as JSON) document.
// This may be used directly with the encoding/json package to marshal or unmarshal ADIJ data.
type ADIJDocument struct {
	// Header is nil when there is no header.
	// Otherwise it is a Record with header fields inside.
	Header map[adifield.ADIField]string `json:"HEADER,omitempty"`

	// Records is a slice of Record.
	// It contains zero or more QSO records.
	Records []map[adifield.ADIField]string `json:"RECORDS"`
}
