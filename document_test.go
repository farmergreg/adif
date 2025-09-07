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
	if firstDoc.header != nil {
		if firstDoc.header.Count() != secondDoc.header.Count() {
			t.Errorf("header length mismatch: got %d, want %d", secondDoc.header.Count(), firstDoc.header.Count())
		}
		for _, field := range firstDoc.header.Fields() {
			if firstDoc.header.Get(field) != secondDoc.header.Get(field) {
				t.Errorf("header field %q mismatch: got %v, want %v", field, secondDoc.header.Get(field), firstDoc.header.Get(field))
			}
		}
	}

	if len(firstDoc.records) != len(secondDoc.records) {
		t.Errorf("records length mismatch: got %d, want %d", len(secondDoc.records), len(firstDoc.records))
	}
	for r := range firstDoc.records {
		if firstDoc.records[r].Count() != secondDoc.records[r].Count() {
			t.Errorf("record %d fields length mismatch: got %d, want %d", r, secondDoc.records[r].Count(), firstDoc.records[r].Count())
		}
		for _, field := range firstDoc.records[r].Fields() {
			if firstDoc.records[r].Get(field) != secondDoc.records[r].Get(field) {
				t.Errorf("record %d field %q mismatch: got %v, want %v", r, field, secondDoc.records[r].Get(field), firstDoc.records[r].Get(field))
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
