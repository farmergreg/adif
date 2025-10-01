package adif

import (
	"io"
	"strings"
	"testing"
)

func TestADIJReaderEmpty(t *testing.T) {
	jsonInput := `{}`

	reader, err := NewJSONDocumentReader(strings.NewReader(jsonInput), false)
	if err != nil {
		t.Fatalf("Failed to create ADIJ reader: %v", err)
	}

	// Should immediately return EOF
	_, _, err = reader.Next()
	if err != io.EOF {
		t.Error("Expected EOF error when no records")
	}
}

func TestADIJReaderInvalidJSON(t *testing.T) {
	// Invalid JSON should cause NewADIJReader to return an error
	invalidJSON := `{"invalid": json syntax`

	_, err := NewJSONDocumentReader(strings.NewReader(invalidJSON), false)
	if err == nil {
		t.Error("Expected error when parsing invalid JSON, but got nil")
	}
}
