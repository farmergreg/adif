package adif

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/hamradiolog-net/adif-spec/v6/adifield"
	"github.com/hamradiolog-net/adif-spec/v6/spec"
)

var qsoWithLou = func() *Record {
	qso := NewRecordWithCapacity(5)
	qso.Set(adifield.CALL, "W9PVA")
	qso.Set(adifield.RST_RCVD, "58")
	qso.Set(adifield.RST_SENT, "59")
	qso.Set(adifield.COMMENT, "Eyeball QSO 👀")
	qso.Set(adifield.QSO_DATE, "") // empty on purpose to test zero length field
	return &qso
}()

// This QSO has exactly 32 fields.
var testADIFSingleRecord = `<APP_LoTW_OWNCALL:5>K9CTS
<STATION_CALLSIGN:5>K9CTS
<MY_DXCC:3>291
<MY_COUNTRY:24>UNITED STATES OF AMERICA
<APP_LoTW_MY_DXCC_ENTITY_STATUS:7>Current
<MY_GRIDSQUARE:6>EN34QU
<MY_STATE:2>WI // Wisconsin
<MY_CNTY:9>WI,PIERCE // Pierce
<MY_CQ_ZONE:2>04
<MY_ITU_ZONE:2>07
<CALL:5>K1ARR
<BAND:3>20M
<FREQ:8>14.06100
<MODE:2>CW
<APP_LoTW_MODEGROUP:2>CW
<QSO_DATE:8>20220122
<APP_LoTW_RXQSO:19>2022-01-22 19:09:09 // QSO record inserted/modified at LoTW
<TIME_ON:6>185309
<APP_LoTW_QSO_TIMESTAMP:20>2022-01-22T18:53:09Z // QSO Date & Time; ISO-8601
<QSL_RCVD:1>Y
<QSLRDATE:8>20220123
<APP_LoTW_RXQSL:19>2022-01-23 00:13:18 // QSL record matched/modified at LoTW
<DXCC:3>291
<COUNTRY:24>UNITED STATES OF AMERICA
<APP_LoTW_DXCC_ENTITY_STATUS:7>Current
<PFX:2>K1
<APP_LoTW_2xQSL:1>Y
<GRIDSQUARE:6>FN34SD
<STATE:2>VT // Vermont
<CNTY:13>VT,WASHINGTON // Washington
<CQZ:2>05
<ITUZ:2>08
<eor>
`

func TestAppendAsADI(t *testing.T) {
	// Arrange
	// the order may be different, but the fields and values must be identical
	expected := "<CALL:5>W9PVA<RST_RCVD:2>58<RST_SENT:2>59<COMMENT:16>Eyeball QSO 👀<QSO_DATE:0>"

	// Act
	adiLength := qsoWithLou.appendAsADIPreCalculate()
	buf := make([]byte, 0, adiLength)
	buf = qsoWithLou.appendAsADI(buf)
	actual := string(buf)

	// Assert
	var expectedLength = len(expected)
	if len(actual) != expectedLength {
		t.Errorf("Expected length %d, got length %d", expectedLength, len(actual))
	}

	expected = "<CALL:5>W9PVA"
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %s to appear in %s", expected, actual)
	}

	expected = "<RST_RCVD:2>58"
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %s to appear in %s", expected, actual)
	}

	expected = "<RST_SENT:2>59"
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %s to appear in %s", expected, actual)
	}

	expected = "<COMMENT:16>Eyeball QSO 👀"
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %s to appear in %s", expected, actual)
	}

	expected = "<QSO_DATE:0>"
	if !strings.Contains(actual, expected) {
		t.Errorf("Expected %s to appear in %s", expected, actual)
	}

	if strings.Contains(actual, TagEOR) {
		t.Errorf("Expected %s to not appear in %s", TagEOR, actual)
	}

	if strings.Contains(actual, TagEOH) {
		t.Errorf("Expected %s to not appear in %s", TagEOH, actual)
	}
}

func TestAppendAsADIPreCalculate(t *testing.T) {
	// Arrange
	var size = rand.Intn(10000000) + (1024 * 50)
	var qso = NewRecord()
	qso.Set(adifield.PROGRAMID, "HamRadioLog.Net")
	qso.Set(adifield.PROGRAMVERSION, strings.Repeat("1", size))
	qso.Set(adifield.ADIF_VER, spec.ADIF_VER)

	adiLength := qso.appendAsADIPreCalculate()
	buf := make([]byte, 0, adiLength)
	buf = qso.appendAsADI(buf)
	expectedLength := len(buf)

	// Act
	length := qso.appendAsADIPreCalculate()

	if length != expectedLength {
		t.Errorf("Expected %d, got %d", expectedLength, length)
	}
}

func TestWriteTo(t *testing.T) {
	// Arrange
	var builder strings.Builder

	// Act
	qsoWithLou.WriteTo(&builder)
	builder.WriteString(TagEOR)

	// Assert
	qso := NewRecord()
	n, err := qso.ReadFrom(strings.NewReader(builder.String()))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	expectedFields := 5
	if qso.Count() != expectedFields {
		t.Errorf("Expected %d fields, got %d", expectedFields, qso.Count())
	}

	expectedLength := 86
	if builder.Len() != expectedLength {
		t.Errorf("Expected length %d, got %d", expectedLength, builder.Len())
	}
	if n != int64(expectedLength) {
		t.Errorf("Expected %d bytes read, got %d", expectedLength, n)
	}
}

func TestReadFrom(t *testing.T) {
	qso := NewRecord()
	qso.ReadFrom(strings.NewReader(testADIFSingleRecord))

	expected := 32
	if qso.Count() != expected {
		t.Errorf("Expected %d fields, got %d", expected, qso.Count())
	}
}
