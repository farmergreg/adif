package adif

import "errors"

var (
	// ErrAdiReaderMalformedADI is returned when the ADI formatted data does not conform to the ADIF specification.
	ErrAdiReaderMalformedADI = errors.New("malformed ADI data")

	// ErrAdiReaderTooManyUniqueFields is returned when the number of unique field names exceeds the internal limit.
	// This prevents denial of service attacks from malformed ADI files with unlimited unique field names.
	ErrAdiReaderTooManyUniqueFields = errors.New("too many unique field names")

	// ErrDocumentUnexpectedHeader is returned when a header record is encountered after QSO records have already been read, or after a header has already been processed.
	ErrDocumentUnexpectedHeader = errors.New("unexpected header record")

	// ErrWriterHeaderAlreadyWritten is returned when attempting to write more than one header record.
	ErrWriterHeaderAlreadyWritten = errors.New("header record already written")
)
