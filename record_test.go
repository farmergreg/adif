package adif

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
	"github.com/hamradiolog-net/adif-spec/src/pkg/spec"
	"github.com/stretchr/testify/assert"
)

var qsoWithLou = func() *Record {
	qso := NewRecordWithCapacity(5)
	qso.fields[adifield.CALL] = "W9PVA"
	qso.fields[adifield.RST_RCVD] = "58"
	qso.fields[adifield.RST_SENT] = "59"
	qso.fields[adifield.COMMENT] = "Eyeball QSO 👀"
	qso.fields[adifield.QSO_DATE] = "" // empty on purpose to test zero length field
	return qso
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
	var expectedLength = len(expected) - len("<QSO_DATE:0>")
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
	if strings.Contains(actual, expected) {
		t.Errorf("Expected %s to NOT appear in %s", expected, actual)
	}
}

func TestAppendAsADIPreCalculate(t *testing.T) {
	// Arrange
	var size = rand.Intn(10000000) + (1024 * 50)
	var qso = NewRecord()
	qso.fields[adifield.PROGRAMID] = "HamRadioLog.Net"
	qso.fields[adifield.PROGRAMVERSION] = strings.Repeat("1", size)
	qso.fields[adifield.ADIF_VER] = spec.ADIFVersion

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

func TestQSOClean(t *testing.T) {
	// Arrange
	qso := NewRecord()
	qso.fields[adifield.CALL] = "W9PVA "
	qso.fields[adifield.COMMENT] = " COMMENT "
	qso.fields[adifield.QSO_DATE] = ""

	// Act
	qso.Clean()

	// Assert
	if len(qso.fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(qso.fields))
	}

	if qso.fields[adifield.COMMENT] != "COMMENT" {
		t.Errorf("Expected \"COMMENT\", got \"%s\"", qso.fields[adifield.COMMENT])
	}

	if qso.fields[adifield.CALL] != "W9PVA" {
		t.Errorf("Expected \"W9PVA\", got \"%s\"", qso.fields[adifield.CALL])
	}
}

func TestWriteTo(t *testing.T) {
	// Arrange
	var builder strings.Builder

	// Act
	qsoWithLou.WriteTo(&builder)

	// Assert
	qso := NewRecord()
	n, err := qso.ReadFrom(strings.NewReader(builder.String()))
	assert.Nil(t, err)
	assert.Equal(t, 4, len(qso.fields))
	assert.Equal(t, 69, builder.Len())
	assert.Equal(t, int64(69), n)
}

func TestReadFrom(t *testing.T) {
	qso := NewRecord()
	qso.ReadFrom(strings.NewReader(testADIFSingleRecord))

	if len(qso.fields) != 32 {
		t.Errorf("Expected 32 fields, got %d", len(qso.fields))
	}
}
