package adif

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
)

//go:embed testdata/*.adi
var testFileFS embed.FS

func TestVerifyRecordCount(t *testing.T) {
	tests := map[string]int{
		"ADIF_315_test_QSOs_2024_11_28.adi": 6156,
		"Log4OM.adi":                        122,
		"N3FJP-AClogAdif.adi":               438,
		"lotwreport.adi":                    438,
		"qrz.adi":                           931,
		"skcc-logger.adi":                   15,
	}

	for filename, expectedCount := range tests {
		t.Run(filename, func(t *testing.T) {
			reader, err := testFileFS.Open("testdata/" + filename)
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()

			p := NewADIReader(reader, true)
			count := 0
			for {
				qso, _, _, err := p.Next()
				if err == io.EOF {
					break
				}
				if qso == nil {
					t.Fatal("Expected non-nil QSO")
				}
				if err != nil {
					t.Fatal(err)
				}
				count++
			}

			if count != expectedCount {
				t.Errorf("Record count mismatch: got %d, want %d", count, expectedCount)
			}
		})
	}
}

func TestParseBasicFunctionality(t *testing.T) {
	tests := []struct {
		hasHeader   bool
		recordCount int
		name        string
		data        string
	}{
		{false, 0, "Empty String", ""},

		{false, 0, "EOR", "<eOR>"},
		{false, 0, "EOH", "<Eoh>"},
		{false, 0, "EOR EOR", "<EOr><eoR>"},
		{false, 0, "EOH EOR", "<EoH><eOr>"},
		{false, 0, "EOH with leading space", " <EOh>"},
		{false, 0, "EOR with leading space", " <EOr>"},
		{false, 0, "EOH with spaces", " <EOh> "},
		{false, 0, "EOR with spaces", " <EOr> "},
		{false, 0, "EOH with trailing space", "<EOh> "},
		{false, 0, "EOR with trailing space", "<EOr> "},
		{false, 0, "EOH with brackets", "><EOh>>"},
		{false, 0, "EOR with brackets", "><EOr>>"},

		{false, 1, "Valid Record", "<CaLL:5>W9PVA<EOR>"},
		{false, 1, "Leading space", " <CaLL:5>W9PVA<EOR>"},
		{false, 1, "Extra character", "<Call:5>W9PVAn<EOR>"},
		{false, 1, "Extra characters around EOR", "<Call:5>W9PVAa<EOR>b"},

		{true, 2, "With header", "<PROGRAMID:4>TEST<EOH><Call:5>W9PVAa<EOR>b"},
		{true, 2, "With header and extra chars", "<PROGRAMID:4>TESTing<EOH><Call:5>W9PVAa<EOR>b"},
		{true, 2, "With header, header preamble and extra chars", "preamble\n<PROGRAMID:4>TESTing<EOH><Call:5>W9PVAa<EOR>b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewADIReader(strings.NewReader(tt.data), false)

			records := make([]Record, 0, 10000)
			for {
				record, _, _, err := p.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				records = append(records, record)
			}

			if len(records) != tt.recordCount {
				t.Errorf("Record count mismatch: got %d, want %d", len(records), tt.recordCount)
			}

			if tt.recordCount == 0 {
				return
			}

			var index = 0
			if tt.hasHeader {
				if records[0][adifield.PROGRAMID] != "TEST" {
					t.Errorf("Expected header record to have PROGRAMID 'TEST', got %s", records[0][adifield.PROGRAMID])
				}
				index++
			}

			if records[index][adifield.CALL] != "W9PVA" {
				t.Errorf("Expected record to have CALL 'W9PVA', got %s", records[index][adifield.CALL])
			}
		})
	}
}

func TestParseWithMissingEOR(t *testing.T) {
	raw := "<CaLL:5>W9PVA"
	p := NewADIReader(strings.NewReader(raw), false)

	qso, _, _, err := p.Next()

	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}

	expectedFields := 1
	if len(qso) != expectedFields {
		t.Errorf("Expected %d fields, got %d", expectedFields, len(qso))
	}
	if qso[adifield.CALL] != "W9PVA" {
		t.Errorf("Expected CALL 'W9PVA', got %s", qso[adifield.CALL])
	}
}

func TestParseWithMissingEOH(t *testing.T) {
	raw := "<ADIF_VER:5>3.1.5"
	p := NewADIReader(strings.NewReader(raw), false)

	qso, _, _, err := p.Next()

	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}
	expectedFields := 1
	if len(qso) != expectedFields {
		t.Errorf("Expected %d fields, got %d", expectedFields, len(qso))
	}
	if qso[adifield.ADIF_VER] != "3.1.5" {
		t.Errorf("Expected ADIF_VER '3.1.5', got %s", qso[adifield.ADIF_VER])
	}
}

func TestParseWithNumbersInFieldName(t *testing.T) {
	raw := "<APP_LoTW_2xQSL:1>Y<EOR>"
	p := NewADIReader(strings.NewReader(raw), false)

	qso, _, bytesRead, err := p.Next()

	if err != nil {
		t.Fatal(err)
	}
	val := qso["APP_LOTW_2XQSL"]
	if val != "Y" {
		t.Errorf("got %q, want %q", val, "Y")
	}
	expectedBytesRead := int64(len(raw))
	if bytesRead != expectedBytesRead {
		t.Errorf("got %d bytes read, want %d", bytesRead, expectedBytesRead)
	}
}

func TestParseWithMissingLengthField(t *testing.T) {
	// Arrange
	raw := "<APP_LoTW_EOF>Y" // n.b. 'Y' is NOT part of the data. It is a comment...
	p := NewADIReader(strings.NewReader(raw), false)

	// Act
	qso, _, bytesRead, err := p.Next()

	// Assert
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}
	val := qso["APP_LOTW_EOF"]
	if val != "" {
		t.Errorf("Expected empty string, got %s", val)
	}
	expectedBytesRead := int64(len(raw))
	if bytesRead != expectedBytesRead {
		t.Errorf("got %d bytes read, want %d", bytesRead, expectedBytesRead)
	}
}

func TestParseNoRecords(t *testing.T) {
	// Arrange
	tests := []struct {
		name                string
		data                string
		isNonEOFErrExpected bool // EOF means success, non-EOF means the adi reader rejected the input as malformed.
	}{

		{"Invalid app field", "<APP_WAAT:fake>", true},
		{"Empty string", "", false},
		{"Single space", " ", false},
		{"Single colon", ":", false},
		{"Double colon", "::", false},
		{"Plain text", "no adif here...", false},
		{"tag close", ">", false},
		{"tag open", "<", true},
		{"Random text with tag", "< some random text", true},
		{"Math expression 1", " 3 < 4 ", true},
		{"Math expression 2", " 3 > 4 ", false},
		{"Incomplete tag 1", "<this is not adif", true},
		{"Incomplete tag 2", "<something random", true},
		{"Incomplete tag with colon and >", "<something random:>", true},
		{"Incomplete tag with colon and space >", "<something random: >", true},
		{"Incomplete tag with colon and space", "<something random: ", true},
		{"Incomplete tag with colon", "<something random:", true},
		{"Incomplete tag with number", "<something random:8", true},
		{"Incomplete tag with n", "<something random:n", true},
		{"Incomplete tag with type 1", "<something random:8:", true},
		{"Incomplete tag with type 2", "<something random:n:", true},
		{"Incomplete tag with type 3", "<something random:8:x", true},
		{"Incomplete tag with type 4", "<something random:n:x", true},
		{"Incomplete data field", "<APP_TEST:1>", true},
		{"Incomplete data field, with type", "<APP_TEST:1:x>", true},
		{"Empty tag", "<>", true},
		{"Empty tag with text", "<>fake", true},
		{"Empty tag with colon", "<:>fake", true},
		{"Empty tag with double colon", "<::>fake", true},
		{"Empty tag with triple colon", "<:::>fake", true},
		{"Empty tag with quad colon", "<::::>fake", true},
		{"tag open and close", "<>", true},
		{"tag open and close with colon", "<:>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			p := NewADIReader(strings.NewReader(tt.data), false)

			// Act
			qso, _, _, err := p.Next()

			// Assert
			if len(qso) != 0 {
				t.Errorf("Expected empty QSO, got %v", qso)
			}

			// Assert
			if tt.isNonEOFErrExpected {
				if err == nil {
					t.Error("Expected non-EOF error, got nil")
				}
				if err == io.EOF {
					t.Error("Expected non-EOF error, got EOF")
				}
			} else if err != io.EOF {
				t.Errorf("Expected EOF error, got %v", err)
			}
		})
	}
}

func TestParseSingleRecord(t *testing.T) {
	tests := []struct {
		name           string
		adifSource     string
		fieldName      string
		fieldData      string
		isHeaderRecord bool
		isExpectEOF    bool
	}{
		{"Header record", "<progRamid:4>MonoLog<EOH>", "PROGRAMID", "Mono", true, false},
		{"Zero length data", "<APP_MY_APP:0>\r\n<EOR>", "APP_MY_APP", "", false, false},
		{"Single char data", "<APP_MY_APP:1>x <EOR>", "APP_MY_APP", "x", false, false},
		{"Basic TIME_ON", "<TIME_ON:6>161819<EOR>", "TIME_ON", "161819", false, false},
		{"TIME_ON with type", "<TIME_ON:6:Time>161819<EOR>", "TIME_ON", "161819", false, false},
		{"Mixed case TIME_ON", "<TiMe_ON:6>161819<EOR>", "TIME_ON", "161819", false, false},
		{"TIME_ON with type and space", "<TIME_ON:6:Time>161819 <EOR>", "TIME_ON", "161819", false, false},
		{"Mixed case with space", "<TiMe_ON:6>161819 <EOR>", "TIME_ON", "161819", false, false},
		{"Leading space with type", " <TIME_ON:6:Time>161819 <EOR>", "TIME_ON", "161819", false, false},
		{"Leading space mixed case", " <TiMe_ON:6>161819 <EOR>", "TIME_ON", "161819", false, false},
		{"Leading space with type no end space", " <TIME_ON:6:Time>161819<EOR>", "TIME_ON", "161819", false, false},
		{"Leading space mixed case no end space", " <TiME_ON:6>161819<EOR>", "TIME_ON", "161819", false, false},
		{"Extra brackets", "><TiME_ON:6>161819<EOR>", "TIME_ON", "161819", false, false},
		{"Long Field Name", "<App_K9CTS_012345678901234567890:5>W9PVA<EOR>", "APP_K9CTS_012345678901234567890", "W9PVA", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use small buffer to test ReadSlice / ErrBufferFull handling
			br := bufio.NewReaderSize(strings.NewReader(tt.adifSource), 16)
			p := NewADIReader(br, false)

			qso, isHeader, bytesRead, err := p.Next()
			if tt.isExpectEOF {
				if err != io.EOF {
					t.Errorf("Expected EOF error, got %v", err)
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
			}

			if qso[adifield.Field(tt.fieldName)] != tt.fieldData {
				t.Errorf("Expected %s field to be %s, got %s", tt.fieldName, tt.fieldData, qso[adifield.Field(tt.fieldName)])
			}

			if isHeader != tt.isHeaderRecord {
				t.Errorf("Expected header record status %v, got %v", tt.isHeaderRecord, isHeader)
			}
			expectedBytesRead := int64(len(tt.adifSource))
			if bytesRead != expectedBytesRead {
				t.Errorf("Expected %d bytes read, got %d", expectedBytesRead, bytesRead)
			}
		})
	}
}

func TestParseSkipHeader(t *testing.T) {
	// Arrange
	adif := "<PROGRAMID:7>MonoLog<EOH>\n<COMMENT:4>GOOD<EOR>"
	p := NewADIReader(strings.NewReader(adif), true)

	// Act & Assert
	record, _, _, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if record[adifield.PROGRAMID] != "" {
		t.Errorf("Expected empty PROGRAMID, got %s", record[adifield.PROGRAMID])
	}
	if record[adifield.COMMENT] != "GOOD" {
		t.Errorf("Expected COMMENT 'GOOD', got %s", record[adifield.COMMENT])
	}

	recordTwo, _, _, errTwo := p.Next()
	if errTwo != io.EOF {
		t.Errorf("Expected EOF error, got %v", errTwo)
	}
	if len(recordTwo) != 0 {
		t.Errorf("Expected empty record, got %v", recordTwo)
	}
}

func TestParseLongFieldName(t *testing.T) {
	const len int = 2000
	fieldName := "APP_K9CTS_" + strings.Repeat("X", len)

	// Arrange
	adif := fmt.Sprintf("<%s:4>TEST<EOR>", fieldName)
	p := NewADIReader(strings.NewReader(adif), false)

	// Act
	record, _, _, err := p.Next()
	_, _, _, _ = p.Next()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if record[adifield.Field(fieldName)] != "TEST" {
		t.Errorf("Expected %s field to be TEST, got %s", fieldName, record[adifield.Field(fieldName)])
	}
}

func TestParseLargeData(t *testing.T) {
	// Arrange
	p := NewADIReader(strings.NewReader("<COMMENT:1000002>0"+strings.Repeat("1", 1_000_000)+"01<EOR>"), false)

	// Act
	record, _, _, err := p.Next() // Force the buffer to be resized to accommodate the large value
	_, _, _, _ = p.Next()         // Force the buffer to be resized back to "normal"

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if record[adifield.COMMENT] != "0"+strings.Repeat("1", 1_000_000)+"0" {
		t.Errorf("Expected %s, got %s", strings.Repeat("1", 1_000_000), record[adifield.COMMENT])
	}
}

func TestParseLargeDataTooBigShouldReturnErr(t *testing.T) {
	// Arrange
	p := NewADIReader(strings.NewReader("<COMMENT:10000002>0"+strings.Repeat("1", 10_000_000)+"01"), false)

	// Act
	record, _, _, err := p.Next() // Force the buffer to be resized to accommodate the large value
	_, _, _, _ = p.Next()         // Force the buffer to be resized back to "normal"

	// Assert
	if err != ErrInvalidFieldLength {
		t.Errorf("Expected ErrInvalidFieldLength error, got %v", err)
	}
	if record[adifield.COMMENT] != "" {
		t.Errorf("Expected empty string, got %s", record[adifield.COMMENT])
	}
}
