package adif

import (
	"io"
	"strings"
	"testing"
)

func TestJSONRecordReaderEmpty(t *testing.T) {
	jsonInput := `{}`

	reader, err := NewJSONRecordReader(strings.NewReader(jsonInput), false)
	if err != nil {
		t.Fatalf("Failed to create ADIJ reader: %v", err)
	}

	// Should immediately return EOF
	_, err = reader.Next()
	if err != io.EOF {
		t.Error("Expected EOF error when no records")
	}
}

func TestJSONRecordReaderInvalidJSON(t *testing.T) {
	// Invalid JSON should cause NewJSONRecordReader to return an error
	invalidJSON := `{"invalid": json syntax`

	_, err := NewJSONRecordReader(strings.NewReader(invalidJSON), false)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON, but got nil")
	}
}
