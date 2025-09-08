package adif

import (
	"io"
	"strings"
	"testing"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
)

func TestADIJReader(t *testing.T) {
	jsonInput := `{
  "RECORDS": [
    {
      "BAND": "20M",
      "CALL": "K9CTS",
      "MODE": "SSB"
    }
  ]
}`

	reader, err := NewADIJReader(strings.NewReader(jsonInput), false)
	if err != nil {
		t.Fatalf("Failed to create ADIJ reader: %v", err)
	}

	// Read the first (and only) record
	record, err := reader.Next()
	if err != nil {
		t.Fatalf("Failed to read record: %v", err)
	}

	// Verify the record is not a header
	if record.IsHeader() {
		t.Error("Expected record to not be a header")
	}

	// Verify the field values
	if record.Get(adifield.BAND) != "20M" {
		t.Errorf("Expected BAND to be '20M', got '%s'", record.Get(adifield.BAND))
	}
	if record.Get(adifield.CALL) != "K9CTS" {
		t.Errorf("Expected CALL to be 'K9CTS', got '%s'", record.Get(adifield.CALL))
	}
	if record.Get(adifield.MODE) != "SSB" {
		t.Errorf("Expected MODE to be 'SSB', got '%s'", record.Get(adifield.MODE))
	}

	// Verify EOF is returned when no more records
	_, err = reader.Next()
	if err != io.EOF {
		t.Error("Expected EOF error when no more records")
	}
}

func TestADIJReaderWithHeader(t *testing.T) {
	jsonInput := `{
  "HEADER": {
    "ADIF_VER": "3.1.4",
    "PROGRAMID": "Test Program"
  },
  "RECORDS": [
    {
      "BAND": "20M",
      "CALL": "K9CTS",
      "MODE": "SSB"
    }
  ]
}`

	reader, err := NewADIJReader(strings.NewReader(jsonInput), false)
	if err != nil {
		t.Fatalf("Failed to create ADIJ reader: %v", err)
	}

	// Read the header record first
	headerRecord, err := reader.Next()
	if err != nil {
		t.Fatalf("Failed to read header record: %v", err)
	}

	// Verify the header record
	if !headerRecord.IsHeader() {
		t.Error("Expected record to be a header")
	}
	if headerRecord.Get(adifield.ADIF_VER) != "3.1.4" {
		t.Errorf("Expected ADIF_VER to be '3.1.4', got '%s'", headerRecord.Get(adifield.ADIF_VER))
	}
	if headerRecord.Get(adifield.PROGRAMID) != "Test Program" {
		t.Errorf("Expected PROGRAMID to be 'Test Program', got '%s'", headerRecord.Get(adifield.PROGRAMID))
	}

	// Read the QSO record
	qsoRecord, err := reader.Next()
	if err != nil {
		t.Fatalf("Failed to read QSO record: %v", err)
	}

	// Verify the QSO record is not a header
	if qsoRecord.IsHeader() {
		t.Error("Expected QSO record to not be a header")
	}
	if qsoRecord.Get(adifield.CALL) != "K9CTS" {
		t.Errorf("Expected CALL to be 'K9CTS', got '%s'", qsoRecord.Get(adifield.CALL))
	}

	// Verify EOF is returned when no more records
	_, err = reader.Next()
	if err != io.EOF {
		t.Error("Expected EOF error when no more records")
	}
}

func TestADIJReaderSkipHeader(t *testing.T) {
	jsonInput := `{
  "HEADER": {
    "ADIF_VER": "3.1.4",
    "PROGRAMID": "Test Program"
  },
  "RECORDS": [
    {
      "BAND": "20M",
      "CALL": "K9CTS",
      "MODE": "SSB"
    }
  ]
}`

	reader, err := NewADIJReader(strings.NewReader(jsonInput), true)
	if err != nil {
		t.Fatalf("Failed to create ADIJ reader: %v", err)
	}

	// Read the first record (should be QSO, not header)
	record, err := reader.Next()
	if err != nil {
		t.Fatalf("Failed to read record: %v", err)
	}

	// Verify the record is not a header (header should be skipped)
	if record.IsHeader() {
		t.Error("Expected record to not be a header (header should be skipped)")
	}
	if record.Get(adifield.CALL) != "K9CTS" {
		t.Errorf("Expected CALL to be 'K9CTS', got '%s'", record.Get(adifield.CALL))
	}

	// Verify EOF is returned when no more records
	_, err = reader.Next()
	if err != io.EOF {
		t.Error("Expected EOF error when no more records")
	}
}

func TestADIJReaderEmpty(t *testing.T) {
	jsonInput := `{}`

	reader, err := NewADIJReader(strings.NewReader(jsonInput), false)
	if err != nil {
		t.Fatalf("Failed to create ADIJ reader: %v", err)
	}

	// Should immediately return EOF
	_, err = reader.Next()
	if err != io.EOF {
		t.Error("Expected EOF error when no records")
	}
}
