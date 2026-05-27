package adif

import "errors"

var (
	// ErrAdiReaderMalformedADI is returned when the ADI formatted data does not conform to the ADIF specification.
	ErrAdiReaderMalformedADI = errors.New("adi reader: data is malformed")

	// ErrHeaderAlreadyWritten is returned when attempting to write more than one header record.
	ErrHeaderAlreadyWritten = errors.New("header record already written")

	// ErrAdiReaderTooManyUniqueFields is returned when the number of unique field names exceeds the internal limit.
	ErrAdiReaderTooManyUniqueFields = errors.New("adi reader: too many unique field names")
)
