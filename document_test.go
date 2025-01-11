package adif

import (
	"io"
	"strings"
	"testing"
)

func TestDocumentString(t *testing.T) {
	// Arrange
	doc := NewDocument()

	doc.ReadFrom(strings.NewReader("<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR>"))

	// Act
	s := doc.String()

	// Assert
	want := adiHeaderPreamble + "<PROGRAMID:7>MonoLog<EOH>\n<CALL:5>W9PVA<EOR>\n"
	if s != want {
		t.Errorf("got %q, want %q", s, want)
	}
}

func TestParseExportParseVerifySimple(t *testing.T) {
	// Arrange
	tests := []struct {
		name       string
		adifSource string
	}{
		{"Simple ADI", "<APP_MY_APP_0:0>\n<APP_MY_APP:1>x<EOR>"},
		{"Simple ADI with header", "<PROGRAMID:7>MonoLog<EOH>\n<APP_MY_APP_0:0><APP_MY_APP:1>x<EOR>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parseExportParseVerifyHelper(t, tt.adifSource)
		})
	}
}

func TestParseExportParseVerifyFiles(t *testing.T) {
	t.Parallel()

	// Arrange
	files, err := testFileFS.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		file := file
		t.Run(file.Name(), func(t *testing.T) {
			t.Parallel()

			// Act
			reader, err := testFileFS.Open("testdata/" + file.Name())
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()

			content, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}

			parseExportParseVerifyHelper(t, string(content))
		})
	}
}

func parseExportParseVerifyHelper(t *testing.T, adif string) {
	t.Helper()

	// Parse
	firstDoc := NewDocument()
	firstDoc.ReadFrom(strings.NewReader(adif))

	// Export
	buf := &strings.Builder{}
	firstDoc.WriteTo(buf)

	// Parse
	secondDoc := NewDocument()
	secondDoc.ReadFrom(strings.NewReader(buf.String()))

	// Verify
	if firstDoc.Header != nil {
		if len(firstDoc.Header.fields) != len(secondDoc.Header.fields) {
			t.Errorf("header length mismatch: got %d, want %d", len(secondDoc.Header.fields), len(firstDoc.Header.fields))
		}
		for field := range firstDoc.Header.fields {
			if firstDoc.Header.fields[field] != secondDoc.Header.fields[field] {
				t.Errorf("header field %q mismatch: got %v, want %v", field, secondDoc.Header.fields[field], firstDoc.Header.fields[field])
			}
		}
	}

	if len(firstDoc.Records) != len(secondDoc.Records) {
		t.Errorf("records length mismatch: got %d, want %d", len(secondDoc.Records), len(firstDoc.Records))
	}
	for r := range firstDoc.Records {
		if len(firstDoc.Records[r].fields) != len(secondDoc.Records[r].fields) {
			t.Errorf("record %d fields length mismatch: got %d, want %d", r, len(secondDoc.Records[r].fields), len(firstDoc.Records[r].fields))
		}
		for field := range firstDoc.Records[r].fields {
			if firstDoc.Records[r].fields[field] != secondDoc.Records[r].fields[field] {
				t.Errorf("record %d field %q mismatch: got %v, want %v", r, field, secondDoc.Records[r].fields[field], firstDoc.Records[r].fields[field])
			}
		}
	}
}

func TestDocumentStringNilReceiver(t *testing.T) {
	// Arrange
	var doc *Document

	// Act
	s := doc.String()

	// Assert
	if s != "" {
		t.Errorf("got %q, want empty string", s)
	}
}
