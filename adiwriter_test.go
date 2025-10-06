package adif

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
	"github.com/farmergreg/spec/v6/spec"
)

func TestADIDocumentWriterWrite(t *testing.T) {
	hdr := NewRecord()
	hdr.Set(adifield.PROGRAMID, "HamRadioLog.Net")
	hdr.Set(adifield.PROGRAMVERSION, "1.0.0")
	hdr.Set(adifield.ADIF_VER, "3.1.4")

	qso := NewRecord()
	qso.Set(adifield.CALL, "K9CTS")

	qso1 := NewRecord()

	sb := &strings.Builder{}
	w := NewADIDocumentWriterWithPreamble(sb, "")

	err := w.WriteHeader(hdr)
	if err != nil {
		t.Fatal(err)
	}

	records := []Record{qso, qso1}
	for _, r := range records {
		err := w.WriteRecord(r)
		if err != nil {
			t.Fatal(err)
		}
	}

	if err = w.Flush(); err != nil {
		t.Fatal(err)
	}

	expectedADIF := "\n<ADIF_VER:5>3.1.4<PROGRAMID:15>HamRadioLog.Net<PROGRAMVERSION:5>1.0.0<EOH>\n<CALL:5>K9CTS<EOR>\n"
	if sb.String() != expectedADIF {
		t.Errorf("Expected '%s', got '%s'", expectedADIF, sb.String())
	}
}

func TestADIDocumentWriterWrite_BigRecord(t *testing.T) {
	// force re-allocation of the internal write buffer.
	var size = rand.Intn(10000000) + (1024 * 50)
	qso := NewRecord()
	qso.Set(adifield.COMMENT, strings.Repeat("1", size))

	sb := &strings.Builder{}
	w := NewADIDocumentWriter(sb)
	err := w.WriteRecord(qso)
	if err != nil {
		t.Fatal(err)
	}

	if err = w.Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestADIDocumentWriterWriteError(t *testing.T) {
	expectedBytes := 20

	qso1 := NewRecord()
	qso1.Set(adifield.CALL, "K9CTS")

	qso2 := NewRecord()
	qso2.Set(adifield.CALL, "W1AW")

	fw := &mockFailWriter{maxBytes: expectedBytes}
	w := NewADIDocumentWriter(fw)

	err := w.WriteRecord(qso1)
	if err != nil {
		t.Fatalf("Expected nil error but got %v", err)
	}

	err = w.WriteRecord(qso2)
	if err == nil {
		t.Fatal("Expected error but got none")
	}

	t.Logf("Error: %v", err)

	if err = w.Flush(); err != nil {
		t.Fatal(err)
	}
}

func TestAppendADIFRecordAsADIPreCalculate(t *testing.T) {
	var size = rand.Intn(10000000) + (1024 * 50)
	qso := NewRecord()
	qso.Set(adifield.PROGRAMID, "HamRadioLog.Net")
	qso.Set(adifield.PROGRAMVERSION, strings.Repeat("1", size))
	qso.Set(adifield.ADIF_VER, spec.ADIF_VER)
	qso.Set("APP_9", strings.Repeat("1", 9))
	qso.Set("APP_99", strings.Repeat("1", 99))
	qso.Set("APP_999", strings.Repeat("1", 999))
	qso.Set("APP_9999", strings.Repeat("1", 9999))
	qso.Set("APP_99999", strings.Repeat("1", 99999))

	preCalculateLength := appendADIFRecordAsADIPreCalculate(qso)
	buf := make([]byte, 0, preCalculateLength)

	buf = appendAsADI(qso, 'R', buf)
	actualLength := len(buf)

	if preCalculateLength != actualLength {
		t.Errorf("Expected %d, got %d", actualLength, preCalculateLength)
	}
}

func TestAppendAsAdifNoLength(t *testing.T) {
	qso := &mockADIFRecord{}
	len := appendADIFRecordAsADIPreCalculate(qso)
	if len != 6 {
		// <EOR>\n = 6 bytes
		t.Errorf("Expected 0, got %d", len)
	}
}

func TestADIDocumentWriterTwoHeaders(t *testing.T) {
	hdr := NewRecord()
	hdr.Set(adifield.PROGRAMVERSION, "1.0.0")

	sb := &strings.Builder{}
	w := NewADIDocumentWriterWithPreamble(sb, "")

	err := w.WriteHeader(hdr)
	if err != nil {
		t.Fatal(err)
	}

	err = w.WriteHeader(hdr)
	if err != ErrHeaderAlreadyWritten {
		t.Fatalf("Expected ErrHeaderAlreadyWritten but got %v", err)
	}

	if err = w.Flush(); err != nil {
		t.Fatal(err)
	}
}
