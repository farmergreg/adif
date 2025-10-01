package adif

import (
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
	"github.com/farmergreg/spec/v6/enum/band"
	"github.com/farmergreg/spec/v6/enum/mode"
	"github.com/farmergreg/spec/v6/spec"
)

func TestJSONWriterEncodeFail(t *testing.T) {
	// Create a writer that will cause json.Encoder.Encode to fail
	mockW := &mockAlwaysErrorWriter{}
	writer := &jsonWriter{w: mockW, doc: &jsonDocument{}}

	qso := NewRecord()
	qso.Set(adifield.CALL, "K9CTS")
	qso.Set(adifield.BAND, band.BAND_20M.String())
	qso.Set(adifield.MODE, mode.SSB.String())

	// Write should return the error from the encoder
	err := writer.WriteRecord(qso)
	if err != nil {
		t.Error("Unexpected error from Write:", err)
	}

	err = writer.Close()
	if err == nil {
		t.Error("Expected error from Close, but got nil")
	}

	if err.Error() != errMockWrite.Error() {
		t.Errorf("Expected error '%s', but got '%s'", errMockWrite.Error(), err.Error())
	}
}

func TestJSONWriterDuplicateHeader(t *testing.T) {
	// Create a writer that will cause json.Encoder.Encode to fail
	sb := &strings.Builder{}
	writer := NewJSONDocumentWriter(sb, "  ")

	qso := NewRecord()
	qso.Set(adifield.ADIF_VER, spec.ADIF_VER)

	err := writer.WriteHeader(qso)
	if err != nil {
		t.Error("Unexpected error from Write:", err)
	}

	err = writer.WriteHeader(qso)
	if err != ErrHeaderAlreadyWritten {
		t.Error("Unexpected error from Write:", err)
	}
}
