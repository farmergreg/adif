package adif

import (
	"io"
	"strings"
	"testing"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
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
		if len(firstDoc.Header) != len(secondDoc.Header) {
			t.Errorf("header length mismatch: got %d, want %d", len(secondDoc.Header), len(firstDoc.Header))
		}
		for field := range firstDoc.Header {
			if (firstDoc.Header)[field] != (secondDoc.Header)[field] {
				t.Errorf("header field %q mismatch: got %v, want %v", field, (secondDoc.Header)[field], (firstDoc.Header)[field])
			}
		}
	}

	if len(firstDoc.Records) != len(secondDoc.Records) {
		t.Errorf("records length mismatch: got %d, want %d", len(secondDoc.Records), len(firstDoc.Records))
	}
	for r := range firstDoc.Records {
		if len(firstDoc.Records[r]) != len(secondDoc.Records[r]) {
			t.Errorf("record %d fields length mismatch: got %d, want %d", r, len(secondDoc.Records[r]), len(firstDoc.Records[r]))
		}
		for field := range firstDoc.Records[r] {
			if firstDoc.Records[r][field] != secondDoc.Records[r][field] {
				t.Errorf("record %d field %q mismatch: got %v, want %v", r, field, secondDoc.Records[r][field], firstDoc.Records[r][field])
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

func TestDocumentWriteEmptyHeader(t *testing.T) {
	// Arrange
	sb := &strings.Builder{}
	doc := NewDocument()
	doc.Header = Record{
		adifield.PROGRAMID: "",
	}

	// Act
	doc.WriteTo(sb)

	// Assert
	if sb.String() != adiHeaderPreamble {
		t.Errorf("got %q, want empty string", sb.String())
	}
}

func TestDocumentWriteEmptyRecord(t *testing.T) {
	// Arrange
	sb := &strings.Builder{}
	doc := NewDocument()
	doc.Records = []Record{{
		adifield.CALL: "",
	}}

	// Act
	doc.WriteTo(sb)

	// Assert
	if sb.String() != "" {
		t.Errorf("got %q, want empty string", sb.String())
	}
}
