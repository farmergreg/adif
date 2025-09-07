package adif

import (
	"io"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

// Document represents a complete ADIF document.
//
// Future Work:
// This type intentionally resembles the ADX XML structure even though XML is not currently supported by this library.
type Document struct {
	// header is nil when there is no header.
	// Otherwise it will be a Record with header fields inside.
	header RecordThing `json:"HEADER,omitempty"`

	// records is a slice of Record.
	// It contains the QSO records.
	records []RecordThing `json:"RECORDS"`

	headerPreamble string
}

func (d Document) Header() RecordThing {
	return d.header
}

func (d Document) Records() []RecordThing {
	return d.records
}

// Record is a map of ADIF fields to their values, representing EITHER a Header record or a QSO record.
// The field keys must be UPPERCASE strings of type adifield.Field.
type Record struct {
	r map[adifield.ADIField]string
}

func (r Record) Get(field adifield.ADIField) string {
	if r.r == nil {
		return ""
	}
	return r.r[field]
}

func (r Record) Set(field adifield.ADIField, value string) {
	if r.r == nil {
		r.r = make(map[adifield.ADIField]string)
	}
	r.r[field] = value
}

func (r Record) Count() int {
	if r.r == nil {
		return 0
	}
	return len(r.r)
}

func (r Record) Fields() []adifield.ADIField {
	if r.r == nil {
		return nil
	}
	fields := make([]adifield.ADIField, 0, len(r.r))
	for field := range r.r {
		fields = append(fields, field)
	}
	return fields
}

type DocumentThing interface {
	Header() RecordThing
	Records() []RecordThing
}

type RecordThing interface {
	Get(field adifield.ADIField) string
	Set(field adifield.ADIField, value string)
	Count() int
	Fields() []adifield.ADIField
	WriteTo(dest io.Writer) (int64, error)
}
