package adif

import "errors"

var (
	// ErrAdiReaderMalformedADI is returned when the ADI formatted data does not conform to the ADIF specification.
	ErrAdiReaderMalformedADI = errors.New("adi reader: data is malformed")

	// ErrHeaderAlreadyWritten is returned when attempting to write more than one header record.
	ErrHeaderAlreadyWritten = errors.New("header record already written")
)
