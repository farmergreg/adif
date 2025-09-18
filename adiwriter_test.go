package adif

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/hamradiolog-net/spec/v6/adifield"
	"github.com/hamradiolog-net/spec/v6/spec"
)

func TestADIRecordWriterWrite(t *testing.T) {
	hdr := NewRecord()
	hdr.SetIsHeader(true)
	hdr.Set(adifield.PROGRAMID, "HamRadioLog.Net")
	hdr.Set(adifield.PROGRAMVERSION, "1.0.0")
	hdr.Set(adifield.ADIF_VER, "3.1.4")

	qso := NewRecord()
	qso.Set(adifield.CALL, "K9CTS")

	qso1 := NewRecord()

	sb := &strings.Builder{}
	w := NewADIRecordWriterWithPreamble(sb, "")

	records := []Record{hdr, qso, qso1}
	for _, r := range records {
		err := w.Write(r)
		if err != nil {
			t.Fatal(err)
		}
	}

	expectedADIF := "\n<adif_ver:5>3.1.4<programid:15>HamRadioLog.Net<programversion:5>1.0.0<eoh>\n<call:5>K9CTS<eor>\n"
	if sb.String() != expectedADIF {
		t.Errorf("Expected '%s', got '%s'", expectedADIF, sb.String())
	}
}

func TestADIRecordWriterWrite_BigRecord(t *testing.T) {
	// force re-allocation of the internal write buffer.
	var size = rand.Intn(10000000) + (1024 * 50)
	qso := NewRecord()
	qso.Set(adifield.COMMENT, strings.Repeat("1", size))

	sb := &strings.Builder{}
	w := NewADIRecordWriter(sb)
	err := w.Write(qso)
	if err != nil {
		t.Fatal(err)
	}
}

func TestADIRecordWriterWriteError(t *testing.T) {
	expectedBytes := 20

	qso1 := NewRecord()
	qso1.Set(adifield.CALL, "K9CTS")

	qso2 := NewRecord()
	qso2.Set(adifield.CALL, "W1AW")

	fw := &mockFailWriter{maxBytes: expectedBytes}
	w := NewADIRecordWriter(fw)

	err := w.Write(qso1)
	if err != nil {
		t.Fatalf("Expected nil error but got %v", err)
	}

	err = w.Write(qso2)
	if err == nil {
		t.Fatal("Expected error but got none")
	}

	t.Logf("Error: %v", err)
}

func TestAppendADIFRecordAsADIPreCalculate(t *testing.T) {
	var size = rand.Intn(10000000) + (1024 * 50)
	qso := NewRecord()
	qso.SetIsHeader(true)
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

	buf = appendAsADI(qso, buf)
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
