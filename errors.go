package adif

import "errors"

var (
	// ErrMalformedADI is returned when the ADI formatted data does not conform to the ADIF specification.
	ErrMalformedADI = errors.New("adi reader: data is malformed")

	// ErrInvalidFieldLength is returned when the field length is invalid, or would cause a large memory allocation.
	ErrInvalidFieldLength = errors.New("adi reader: invalid field length")
)

// ErrDocumentTooLarge is returned when the ADI document is too large.
// For large documents, consider using the ADI reader to stream the records.
//
// DocumentMaxSizeInBytes controls the maximum size of data read into an Document struct in bytes.
// This is to prevent memory exhaustion attacks.
// You can change this value to suit your needs.
// The default is 256MB.
var ErrDocumentTooLarge = errors.New("adi document reader: document is too large")
