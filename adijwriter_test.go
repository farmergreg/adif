package adif

import (
	"testing"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

func TestADIJWriterWriteError(t *testing.T) {
	// Create a writer that will cause json.Encoder.Encode to fail
	mockW := &mockAlwaysErrorWriter{}
	writer := &adijWriter{w: mockW}

	qso := NewADIRecord()
	qso.Set(adifield.BAND, "20M")
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.MODE, "SSB")

	// Write should return the error from the encoder
	err := writer.Write([]ADIFRecord{qso})
	if err == nil {
		t.Error("Expected error from Write, but got nil")
	}

	if err.Error() != errMockWrite.Error() {
		t.Errorf("Expected error '%s', but got '%s'", errMockWrite.Error(), err.Error())
	}
}
