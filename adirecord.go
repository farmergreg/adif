package adif

import (
	"slices"
	"strings"
	"unsafe"

	"github.com/hamradiolog-net/spec/v6/adifield"
)

var _ ADIFRecord = (*adiRecord)(nil)

// adiRecord represents a single ADI record.
type adiRecord struct {
	fields   map[adifield.ADIField][]byte // map of all fields and their values
	allCache []adifield.ADIField          // all sorts the keys prior to iterating, this caches that work
	isHeader bool                         // true if this record is a header
}

// NewADIRecord creates a new adiRecord with the default initial capacity.
func NewADIRecord() *adiRecord {
	return NewADIRecordWithCapacity(-1)
}

// NewADIRecordWithCapacity creates a new adiRecord with a specific initial capacity.
func NewADIRecordWithCapacity(initialCapacity int) *adiRecord {
	if initialCapacity < 1 {
		initialCapacity = 7
	}
	r := adiRecord{
		fields: make(map[adifield.ADIField][]byte, initialCapacity),
	}
	return &r
}

// reset clears the record for reuse.
func (r *adiRecord) reset() {
	clear(r.fields)
	r.allCache = r.allCache[:0]
	r.isHeader = false
}

// IsHeader implements ADIFRecord.IsHeader
func (r *adiRecord) IsHeader() bool {
	return r.isHeader
}

// SetIsHeader implements ADIFRecord.SetIsHeader
func (r *adiRecord) SetIsHeader(isHeader bool) {
	r.isHeader = isHeader
}

// Get implements ADIFRecord.Get
func (r *adiRecord) Get(field adifield.ADIField) string {
	field = adifield.ADIField(strings.ToUpper(string(field)))
	data := r.fields[field]
	if data == nil {
		return ""
	}
	return unsafe.String(&data[0], len(data)) // avoid allocation if possible
}

// Set implements ADIFRecord.Set
func (r *adiRecord) Set(field adifield.ADIField, value string) {
	// this avoids unnecessary allocations and string interning in the common case (already upper cased field name string).
	// if the field is not found, we assume it is a new field, clear the cache and normalize it.
	if _, ok := r.fields[field]; !ok {
		r.allCache = nil
		field = adifield.ADIField(strings.ToUpper(string(field)))
	}

	r.setInternal(field, []byte(value))
}

// setInternal sets the value for a field without modifying the field name or clearing the cache.
// It is used by the parser to avoid unnecessary allocations.
// It assumes the field name is already normalized (UPPERCASE).
func (r *adiRecord) setInternal(field adifield.ADIField, value []byte) {
	if len(value) == 0 {
		delete(r.fields, field)
	} else {
		r.fields[field] = value
	}
}

// All implements ADIFRecord.All
func (r *adiRecord) All() func(func(adifield.ADIField, string) bool) {
	if r.allCache == nil {
		r.allCache = make([]adifield.ADIField, 0, len(r.fields))
		for field := range r.fields {
			r.allCache = append(r.allCache, field)
		}
		slices.Sort(r.allCache)
	}

	return func(yield func(adifield.ADIField, string) bool) {
		for _, field := range r.allCache {
			data := r.fields[field]
			if !yield(field, unsafe.String(&data[0], len(data))) {
				return
			}
		}
	}
}

// Count implements ADIFRecord.Count
func (r *adiRecord) Count() int {
	return len(r.fields)
}
