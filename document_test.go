package adif

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
)

var errMockFlush = errors.New("mock flush error")

// mockFlushErrorWriter succeeds on Write but fails on Flush.
type mockFlushErrorWriter struct{}

func (m *mockFlushErrorWriter) Write(p []byte) (int, error) { return len(p), nil }
func (m *mockFlushErrorWriter) Flush() error                { return errMockFlush }

func TestNewDocument(t *testing.T) {
	d := NewDocument()
	if d == nil {
		t.Fatal("expected non-nil document")
	}
	if d.Header != nil {
		t.Error("expected nil header")
	}
	if len(d.Records) != 0 {
		t.Error("expected empty records slice")
	}
}

func TestDocument_ReadFrom_QSOsOnly(t *testing.T) {
	adi := "<CALL:5>W9PVA<BAND:3>20m<EOR><CALL:5>K9CTS<BAND:3>40m<EOR>"
	d := NewDocument()
	n, err := d.ReadFrom(strings.NewReader(adi))
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes read")
	}
	if d.Header != nil {
		t.Error("expected nil header")
	}
	if len(d.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(d.Records))
	}
	if d.Records[0][adifield.CALL] != "W9PVA" {
		t.Errorf("Records[0] CALL: got %q, want %q", d.Records[0][adifield.CALL], "W9PVA")
	}
	if d.Records[1][adifield.CALL] != "K9CTS" {
		t.Errorf("Records[1] CALL: got %q, want %q", d.Records[1][adifield.CALL], "K9CTS")
	}
}

func TestDocument_ReadFrom_WithHeader(t *testing.T) {
	adi := "<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<EOR>"
	d := NewDocument()
	_, err := d.ReadFrom(strings.NewReader(adi))
	if err != nil {
		t.Fatal(err)
	}
	if d.Header == nil {
		t.Fatal("expected header record")
	}
	if d.Header[adifield.PROGRAMID] != "MonoLog" {
		t.Errorf("header PROGRAMID: got %q, want %q", d.Header[adifield.PROGRAMID], "MonoLog")
	}
	if len(d.Records) != 1 {
		t.Fatalf("expected 1 QSO record, got %d", len(d.Records))
	}
	if d.Records[0][adifield.CALL] != "W9PVA" {
		t.Errorf("Records[0] CALL: got %q, want %q", d.Records[0][adifield.CALL], "W9PVA")
	}
}

func TestDocument_ReadFrom_MultipleHeaders(t *testing.T) {
	adi := "<PROGRAMID:7>MonoLog<EOH><PROGRAMID:5>Other<EOH><CALL:5>W9PVA<EOR>"
	d := NewDocument()
	_, err := d.ReadFrom(strings.NewReader(adi))
	if !errors.Is(err, ErrUnexpectedHeader) {
		t.Fatalf("got %v, want ErrDocumentMultipleHeaders", err)
	}
}

func TestDocument_WriteTo_RoundTrip(t *testing.T) {
	adi := "<PROGRAMID:7>MonoLog<EOH><CALL:5>W9PVA<BAND:3>20m<EOR>"
	d := NewDocument()
	if _, err := d.ReadFrom(strings.NewReader(adi)); err != nil {
		t.Fatal(err)
	}

	var sb strings.Builder
	n, err := d.WriteTo(&sb)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Error("expected non-zero bytes written")
	}

	// Re-read the serialized document and verify the data survived the round trip.
	d2 := NewDocument()
	if _, err := d2.ReadFrom(strings.NewReader(sb.String())); err != nil {
		t.Fatal(err)
	}
	if d2.Header == nil {
		t.Fatal("expected header after round trip")
	}
	if d2.Header[adifield.PROGRAMID] != "MonoLog" {
		t.Errorf("header PROGRAMID after round trip: got %q, want %q", d2.Header[adifield.PROGRAMID], "MonoLog")
	}
	if len(d2.Records) != 1 {
		t.Fatalf("expected 1 record after round trip, got %d", len(d2.Records))
	}
	if d2.Records[0][adifield.CALL] != "W9PVA" {
		t.Errorf("CALL after round trip: got %q, want %q", d2.Records[0][adifield.CALL], "W9PVA")
	}
}

func TestDocument_String_Empty(t *testing.T) {
	d := NewDocument()
	if d.String() != "" {
		t.Errorf("expected empty string, got %q", d.String())
	}
}

func TestDocument_String_NilDocument(t *testing.T) {
	var d *Document
	if d.String() != "" {
		t.Errorf("expected empty string for nil document")
	}
}

func TestDocument_String_WithContent(t *testing.T) {
	d := NewDocument()
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	d.Records = append(d.Records, r)

	s := d.String()
	if !strings.Contains(s, "K9CTS") {
		t.Errorf("expected K9CTS in String() output, got %q", s)
	}
}

func TestDocument_WriteTo_WriteError(t *testing.T) {
	d := NewDocument()
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	d.Records = append(d.Records, r)

	_, err := d.WriteTo(&mockFailWriter{maxBytes: 5})
	if err == nil {
		t.Fatal("expected a write error")
	}
}

func TestDocument_WriteTo_HeaderWriteError(t *testing.T) {
	d := NewDocument()
	d.Header = NewRecord()
	d.Header[adifield.PROGRAMID] = "Test"

	_, err := d.WriteTo(&mockAlwaysErrorWriter{})
	if err != errMockWrite {
		t.Fatalf("got %v, want errMockWrite", err)
	}
}

func TestDocument_WriteTo_FlushError(t *testing.T) {
	// An empty document skips WriteHeader and Write entirely, going straight to Flush.
	d := NewDocument()
	_, err := d.WriteTo(&mockFlushErrorWriter{})
	if err != errMockFlush {
		t.Fatalf("got %v, want errMockFlush", err)
	}
}

func TestDocument_JSON_RoundTrip(t *testing.T) {
	d := NewDocument()
	d.Header = Record{adifield.PROGRAMID: "MonoLog"}
	d.Records = append(d.Records, Record{adifield.CALL: "W9PVA", adifield.BAND: "20m"})
	d.Records = append(d.Records, Record{adifield.CALL: "K9CTS", adifield.BAND: "40m"})

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}

	var d2 Document
	if err := json.Unmarshal(data, &d2); err != nil {
		t.Fatal(err)
	}

	if d2.Header[adifield.PROGRAMID] != "MonoLog" {
		t.Errorf("PROGRAMID: got %q, want %q", d2.Header[adifield.PROGRAMID], "MonoLog")
	}
	if len(d2.Records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(d2.Records))
	}
	if d2.Records[0][adifield.CALL] != "W9PVA" {
		t.Errorf("Records[0] CALL: got %q, want %q", d2.Records[0][adifield.CALL], "W9PVA")
	}
	if d2.Records[1][adifield.CALL] != "K9CTS" {
		t.Errorf("Records[1] CALL: got %q, want %q", d2.Records[1][adifield.CALL], "K9CTS")
	}
}
