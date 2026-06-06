package adif

import (
	"bufio"
	"math/rand"
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
)

func TestWriter_Write(t *testing.T) {
	hdr := NewRecord()
	hdr[adifield.PROGRAMID] = "HamRadioLog.Net"
	hdr[adifield.PROGRAMVERSION] = "1.0.0"
	hdr[adifield.ADIF_VER] = "3.1.4"

	qso := NewRecord()
	qso[adifield.CALL] = "K9CTS"

	emptyQSO := NewRecord()

	var sb strings.Builder
	w := NewWriterWithPreamble(&sb, "")

	if err := w.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	for _, r := range []Record{qso, emptyQSO} {
		if err := w.Write(r); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	// Non-priority fields (ADIF_VER, PROGRAMID, PROGRAMVERSION) are written in map
	// iteration order, so we verify values by round-tripping through the Scanner.
	d := NewDocument()
	if _, err := d.ReadFrom(strings.NewReader(sb.String())); err != nil {
		t.Fatalf("round-trip read failed: %v", err)
	}
	if d.Header == nil {
		t.Fatal("expected header record")
	}
	if d.Header[adifield.ADIF_VER] != "3.1.4" {
		t.Errorf("ADIF_VER: got %q, want %q", d.Header[adifield.ADIF_VER], "3.1.4")
	}
	if d.Header[adifield.PROGRAMID] != "HamRadioLog.Net" {
		t.Errorf("PROGRAMID: got %q, want %q", d.Header[adifield.PROGRAMID], "HamRadioLog.Net")
	}
	if d.Header[adifield.PROGRAMVERSION] != "1.0.0" {
		t.Errorf("PROGRAMVERSION: got %q, want %q", d.Header[adifield.PROGRAMVERSION], "1.0.0")
	}
	if len(d.Records) != 1 { // emptyQSO produces no output
		t.Fatalf("expected 1 QSO record, got %d", len(d.Records))
	}
	if d.Records[0][adifield.CALL] != "K9CTS" {
		t.Errorf("CALL: got %q, want %q", d.Records[0][adifield.CALL], "K9CTS")
	}
}

func TestWriter_Write_EmptyRecordsProduceNoOutput(t *testing.T) {
	emptyMap := NewRecord()

	allBlankValues := NewRecord()
	allBlankValues[adifield.CALL] = ""
	allBlankValues[adifield.BAND] = ""

	for _, r := range []Record{emptyMap, allBlankValues} {
		var sb strings.Builder
		w := NewWriterWithPreamble(&sb, "")
		if err := w.Write(r); err != nil {
			t.Fatal(err)
		}
		if got := sb.String(); got != "" {
			t.Errorf("expected no output for empty record, got %q", got)
		}
	}
}

func TestWriter_Write_BigRecord(t *testing.T) {
	// Force the internal write buffer to reallocate.
	size := rand.Intn(10_000_000) + (1024 * 50)
	r := NewRecord()
	r[adifield.COMMENT] = strings.Repeat("1", size)

	var sb strings.Builder
	w := NewWriter(&sb)
	if err := w.Write(r); err != nil {
		t.Fatal(err)
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestWriter_Write_Error(t *testing.T) {
	qso1 := NewRecord()
	qso1[adifield.CALL] = "K9CTS"

	qso2 := NewRecord()
	qso2[adifield.CALL] = "W1AW"

	fw := &mockFailWriter{maxBytes: 20}
	w := NewWriter(fw)

	if err := w.Write(qso1); err != nil {
		t.Fatalf("expected nil error on first write, got %v", err)
	}
	if err := w.Write(qso2); err == nil {
		t.Fatal("expected error on second write")
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestWriter_WriteHeader_PreambleError(t *testing.T) {
	hdr := NewRecord()
	hdr[adifield.PROGRAMID] = "Test"

	w := NewWriterWithPreamble(&mockAlwaysErrorWriter{}, "preamble")
	if err := w.WriteHeader(hdr); err != errMockWrite {
		t.Fatalf("got %v, want errMockWrite", err)
	}
}

func TestWriter_WriteHeader_Twice(t *testing.T) {
	hdr := NewRecord()
	hdr[adifield.PROGRAMVERSION] = "1.0.0"

	var sb strings.Builder
	w := NewWriterWithPreamble(&sb, "")

	if err := w.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteHeader(hdr); err != ErrWriterHeaderAlreadyWritten {
		t.Fatalf("got %v, want ErrWriterHeaderAlreadyWritten", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestWriter_Write_FastMode(t *testing.T) {
	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	r[adifield.BAND] = "20M"

	var sb strings.Builder
	w := NewWriterWithPreamble(&sb, "").SetWriteMode(ADIWriteModeFast)
	if err := w.Write(r); err != nil {
		t.Fatal(err)
	}

	got := sb.String()
	if !strings.Contains(got, "<CALL:5>K9CTS") {
		t.Errorf("CALL missing from Fast mode output: %q", got)
	}
	if !strings.Contains(got, "<BAND:3>20M") {
		t.Errorf("BAND missing from Fast mode output: %q", got)
	}
}

func TestWriter_Flush_Buffered(t *testing.T) {
	// When the underlying writer is a bufio.Writer, Flush must delegate to it.
	var sb strings.Builder
	bw := bufio.NewWriter(&sb)
	w := NewWriter(bw)

	r := NewRecord()
	r[adifield.CALL] = "K9CTS"
	if err := w.Write(r); err != nil {
		t.Fatal(err)
	}

	// Data lives in bufio's buffer — not yet in sb.
	if strings.Contains(sb.String(), "K9CTS") {
		t.Error("expected data to be buffered before Flush")
	}

	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(sb.String(), "K9CTS") {
		t.Error("expected K9CTS in output after Flush")
	}
}
