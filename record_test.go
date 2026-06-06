package adif

import (
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
	"github.com/farmergreg/spec/v6/enum/band"
)

func TestNewRecord(t *testing.T) {
	r := NewRecord()
	if r == nil {
		t.Fatal("expected non-nil record")
	}
}

func TestRecord_MapAccess(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	if r[adifield.CALL] != "K9CTS" {
		t.Errorf("got %q, want %q", r[adifield.CALL], "K9CTS")
	}
}

func TestRecord_Delete(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	delete(r, adifield.CALL)
	if r[adifield.CALL] != "" {
		t.Errorf("expected empty after delete, got %q", r[adifield.CALL])
	}
}

func TestRecord_String_Empty(t *testing.T) {
	r := NewRecord()
	if r.String() != "" {
		t.Errorf("expected empty string, got %q", r.String())
	}
}

func TestRecord_String_PriorityOrder(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	r[adifield.BAND] = band.BAND_20M.String()
	got := r.String()
	// BAND (priority index 2) must appear before CALL (priority index 5).
	want := "<BAND:3>20M<CALL:5>K9CTS"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRecord_String_EmptyValuesOmitted(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	r[adifield.COMMENT] = "" // empty values must not appear in output
	got := r.String()
	want := "<CALL:5>K9CTS"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRecord_WriteTo_Pretty(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	var sb strings.Builder
	n, err := r.WriteTo(&sb)
	if err != nil {
		t.Fatal(err)
	}
	want := "<CALL:5>K9CTS"
	if sb.String() != want {
		t.Errorf("got %q, want %q", sb.String(), want)
	}
	if n != int64(len(want)) {
		t.Errorf("n=%d, want %d", n, len(want))
	}
}

func TestRecord_WriteToMode_Fast(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	var sb strings.Builder
	n, err := r.WriteToMode(&sb, WriteModeFast)
	if err != nil {
		t.Fatal(err)
	}
	want := "<CALL:5>K9CTS"
	if sb.String() != want {
		t.Errorf("got %q, want %q", sb.String(), want)
	}
	if n != int64(len(want)) {
		t.Errorf("n=%d, want %d", n, len(want))
	}
}

func TestRecord_WriteTo_Empty(t *testing.T) {
	r := NewRecord()
	var sb strings.Builder
	n, err := r.WriteTo(&sb)
	if err != nil {
		t.Fatal(err)
	}
	if sb.String() != "" {
		t.Errorf("got %q, want empty", sb.String())
	}
	if n != 0 {
		t.Errorf("n=%d, want 0", n)
	}
}
