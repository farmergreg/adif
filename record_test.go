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
	qso := NewRecord(1)
	qso.Set(adifield.CALL, "W9PVA")
	qso.Set(adifield.RST_RCVD, "58")
	qso.Set(adifield.RST_SENT, "59")
	qso.Set(adifield.COMMENT, "Eyeball QSO 👀")
	qso.Set(adifield.QSO_DATE, "") // empty on purpose to test zero length field
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
	expected := "<CALL:5>W9PVA<RST_RCVD:2>58<RST_SENT:2>59<COMMENT:16>Eyeball QSO 👀<QSO_DATE:0><EOR>"

	// Act
	adiLength, isHeader := qsoWithLou.appendAsADIPreCalculate()
	buf := make([]byte, 0, adiLength)
	buf = qsoWithLou.appendAsADI(buf, isHeader)
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

	if !strings.HasSuffix(actual, TagEOR) {
		t.Errorf("Expected %s to end with %s", actual, TagEOR)
	}
}

func TestAppendAsADIPreCalculate(t *testing.T) {
	// Arrange
	var size = rand.Intn(10000000) + (1024 * 50)
	var qso = *NewRecord(1)
	qso.Set(adifield.PROGRAMID, "HamRadioLog.Net")
	qso.Set(adifield.PROGRAMVERSION, strings.Repeat("1", size))
	qso.Set(adifield.ADIF_VER, spec.ADIFVersion)

	adiLength, isHeader := qso.appendAsADIPreCalculate()
	buf := make([]byte, 0, adiLength)
	buf = qso.appendAsADI(buf, isHeader)
	expectedLength := len(buf)

	// Act
	length, _ := qso.appendAsADIPreCalculate()

	if length != expectedLength {
		t.Errorf("Expected %d, got %d", expectedLength, length)
	}
}

func TestQSOClean(t *testing.T) {
	// Arrange
	qso := *NewRecord(1)
	qso.Set(adifield.CALL, "W9PVA ")
	qso.Set(adifield.COMMENT, " COMMENT ")
	qso.Set(adifield.QSO_DATE, "")

	// Act
	qso.Clean()

	// Assert
	if len(qso.Fields) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(qso.Fields))
	}

	if qso.Get(adifield.COMMENT) != "COMMENT" {
		t.Errorf("Expected \"COMMENT\", got \"%s\"", qso.Get(adifield.COMMENT))
	}

	if qso.Get(adifield.CALL) != "W9PVA" {
		t.Errorf("Expected \"W9PVA\", got \"%s\"", qso.Get(adifield.CALL))
	}
}

func TestIsHeaderRecord(t *testing.T) {
	tests := []struct {
		name           string
		qso            Record
		want           bool
		wantConclusive bool
	}{
		{
			name:           "Header record",
			qso:            *NewRecord(1).Set(adifield.ADIF_VER, spec.ADIFVersion),
			want:           true,
			wantConclusive: true,
		},
		{
			name:           "QSO record",
			qso:            *NewRecord(1).Set(adifield.CALL, "W9PVA"),
			want:           false,
			wantConclusive: true,
		},
		{
			name:           "User defined record",
			qso:            *NewRecord(1).Set(adifield.USERDEF1, "Concertina"),
			want:           true,
			wantConclusive: true,
		},
		{
			name:           "User defined record",
			qso:            *NewRecord(1).Set("APP_UNKNOWN", "Concertina"),
			want:           false,
			wantConclusive: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotConclusive := tt.qso.isHeaderRecord()
			if got != tt.want {
				t.Errorf("QSO.IsHeaderRecord() = %v, want %v", got, tt.want)
			}
			if gotConclusive != tt.wantConclusive {
				t.Errorf("QSO.IsHeaderRecord() = %v, want %v", gotConclusive, tt.wantConclusive)
			}
		})
	}
}

func TestWriteTo(t *testing.T) {
	// Arrange
	var builder strings.Builder

	// Act
	qsoWithLou.WriteTo(&builder)

	// Assert
	qso := Record{}
	n, err := qso.ReadFrom(strings.NewReader(builder.String()))
	assert.Nil(t, err)
	assert.Equal(t, 4, len(qso.Fields))
	assert.Equal(t, 74, builder.Len())
	assert.Equal(t, int64(74), n)
}

func TestReadFrom(t *testing.T) {
	qso := Record{}
	qso.ReadFrom(strings.NewReader(testADIFSingleRecord))

	if len(qso.Fields) != 32 {
		t.Errorf("Expected 32 fields, got %d", len(qso.Fields))
	}
}
