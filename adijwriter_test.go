package adif

import (
	"testing"

	"github.com/hamradiolog-net/spec/v6/adifield"
	"github.com/hamradiolog-net/spec/v6/enum/band"
	"github.com/hamradiolog-net/spec/v6/enum/mode"
)

func TestADIJWriterWriteError(t *testing.T) {
	// Create a writer that will cause json.Encoder.Encode to fail
	mockW := &mockAlwaysErrorWriter{}
	writer := &adijWriter{w: mockW}

	qso := NewADIRecord()
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.BAND, band.Band20m.String())
	qso.Set(adifield.MODE, mode.SSB.String())

	// Write should return the error from the encoder
	err := writer.Write([]ADIFRecord{qso})
	if err == nil {
		t.Error("Expected error from Write, but got nil")
	}

	if err.Error() != errMockWrite.Error() {
		t.Errorf("Expected error '%s', but got '%s'", errMockWrite.Error(), err.Error())
	}
}
