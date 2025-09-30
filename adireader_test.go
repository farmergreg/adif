package adif

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
)

//go:embed testdata/*.adi
var testFileFS embed.FS

func TestADIRecordReaderVerifyRecordCount(t *testing.T) {
	tests := map[string]int{
		"ADIF_316_test_QSOs_2025_08_27.adi": 6191,
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

			p := NewADIRecordReader(reader, true)
			count := 0
			for {
				_, _, err := p.Next()
				if err == io.EOF {
					break
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

func TestADIRecordReaderParseBasicFunctionality(t *testing.T) {
	tests := []struct {
		hasHeader   bool
		recordCount int
		name        string
		data        string
	}{
		{false, 0, "Empty String", ""},

		{false, 1, "Valid Record", "<CaLL:5>W9PVA<EOr>"},
		{false, 1, "Leading space", " <CaLL:5>W9PVA<eor>"},
		{false, 1, "Extra character", "<Call:5>W9PVAn<EOR>"},
		{false, 1, "Extra characters around EOR", "<Call:5>W9PVAa<EoR>b"},

		{true, 2, "With header", "<PROGRAMID:4>TEST<EOH><Call:5>W9PVAa<EOR>b"},
		{true, 2, "With header and extra chars", "<PROGRAMID:4>TESTing<EOH><Call:5>W9PVAa<EOR>b"},
		{true, 2, "With header, header preamble and extra chars", "preamble\n<PROGRAMID:4>TESTing<EOH><Call:5>W9PVAa<EOR>b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewADIRecordReader(strings.NewReader(tt.data), false)

			records := make([]Record, 0, 10000)
			foundHeader := false
			for {
				record, isHeader, err := p.Next()
				foundHeader = foundHeader || isHeader
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
				if !foundHeader {
					t.Errorf("Expected first record to be a header")
				}
				if records[0].Get(adifield.PROGRAMID) != "TEST" {
					t.Errorf("Expected header record to have PROGRAMID 'TEST', got %s", records[0].Get(adifield.PROGRAMID))
				}
				index++
			}

			if records[index].Get(adifield.CALL) != "W9PVA" {
				t.Errorf("Expected record to have CALL 'W9PVA', got %s", records[index].Get(adifield.CALL))
			}
		})
	}
}

func TestADIRecordReaderParseEOREOH(t *testing.T) {
	tests := []struct {
		expected  int
		hasHeader bool
		name      string
		data      string
	}{
		{1, false, "EOR", "<eOR>"},
		{1, true, "EOH", "<Eoh>"},
		{2, false, "EOR EOR", "<EOr><eoR>"},
		{2, true, "EOH EOR", "<EoH><eOr>"},
		{1, true, "EOH with leading space", " <EOh>"},
		{1, false, "EOR with leading space", " <EOr>"},
		{1, true, "EOH with spaces", " <EOh> "},
		{1, false, "EOR with spaces", " <EOr> "},
		{1, true, "EOH with trailing space", "<EOh> "},
		{1, false, "EOR with trailing space", "<EOr> "},
		{1, true, "EOH with brackets", "><EOh>>"},
		{1, false, "EOR with brackets", "><EOr>>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewADIRecordReader(strings.NewReader(tt.data), false)

			records := make([]Record, 0, 10000)
			for {
				record, _, err := p.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					t.Fatal(err)
				}
				records = append(records, record)
			}

			if len(records) != tt.expected {
				t.Errorf("Record count mismatch: got %d, want %d", len(records), tt.expected)
			}
		})
	}
}

func TestADIRecordReaderParseLoTWEOF(t *testing.T) {
	raw := "<" + string(adifield.APP_LOTW_EOF) + ">"
	p := NewADIRecordReader(strings.NewReader(raw), false)

	qso, _, err := p.Next()
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}

	if qso != nil {
		t.Errorf("Expected nil record, got %v", qso)
	}
}

func TestADIRecordReaderParseWithMissingEOH(t *testing.T) {
	raw := "<ADIF_VER:5>3.1.5<eor>"
	p := NewADIRecordReader(strings.NewReader(raw), false)

	qso, _, err := p.Next()
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	expectedFields := 1
	if qso.Count() != expectedFields {
		t.Errorf("Expected %d fields, got %d", expectedFields, qso.Count())
	}
	if qso.Get(adifield.ADIF_VER) != "3.1.5" {
		t.Errorf("Expected ADIF_VER '3.1.5', got %s", qso.Get(adifield.ADIF_VER))
	}

	_, _, err = p.Next()
	if err != io.EOF {
		t.Errorf("Expected EOF error, got %v", err)
	}
}

func TestADIRecordReaderParseWithNumbersInFieldName(t *testing.T) {
	raw := "<APP_LoTW_2xQSL:1>Y<EOR>"
	p := NewADIRecordReader(strings.NewReader(raw), false)

	qso, _, err := p.Next()

	if err != nil {
		t.Fatal(err)
	}
	val := qso.Get(adifield.New("app_lotw_2xqsl"))
	if val != "Y" {
		t.Errorf("got %q, want %q", val, "Y")
	}
}

func TestADIRecordReaderParseNoRecords(t *testing.T) {
	// Arrange
	tests := []struct {
		name        string
		data        string
		expectedErr error // EOF means success, non-EOF means the adi reader rejected the input as malformed.
	}{

		{"Invalid Length", "<APP_WAAT:fake>", ErrAdiReaderMalformedADI},
		{"Empty string", "", io.EOF},
		{"Single space", " ", io.EOF},
		{"Single colon", ":", io.EOF},
		{"Double colon", "::", io.EOF},
		{"Plain text", "no adif here...", io.EOF},
		{"tag close", ">", io.EOF},
		{"tag open", "<", ErrAdiReaderMalformedADI},
		{"Random text with tag", "< some random text", ErrAdiReaderMalformedADI},
		{"Math expression 1", " 3 < 4 ", ErrAdiReaderMalformedADI},
		{"Math expression 2", " 3 > 4 ", io.EOF},
		{"Incomplete tag 1", "<this is not adif", ErrAdiReaderMalformedADI},
		{"Incomplete tag with colon and >", "<something random:>", ErrAdiReaderMalformedADI},
		{"Incomplete tag with colon and space >", "<something random: >", ErrAdiReaderMalformedADI},
		{"Incomplete tag with colon and space", "<something random: ", ErrAdiReaderMalformedADI},
		{"Incomplete tag with colon", "<something random:", ErrAdiReaderMalformedADI},
		{"Incomplete tag with number", "<something random:8", ErrAdiReaderMalformedADI},
		{"Incomplete tag with n", "<something random:n", ErrAdiReaderMalformedADI},
		{"Incomplete tag with type 1", "<something random:8:", ErrAdiReaderMalformedADI},
		{"Incomplete tag with type 2", "<something random:n:", ErrAdiReaderMalformedADI},
		{"Incomplete tag with type 3", "<something random:8:x", ErrAdiReaderMalformedADI},
		{"Incomplete tag with type 4", "<something random:n:x", ErrAdiReaderMalformedADI},
		{"Incomplete data field", "<APP_TEST:1>", ErrAdiReaderMalformedADI},
		{"Incomplete data field, with type", "<APP_TEST:1:x>", ErrAdiReaderMalformedADI},
		{"Empty tag", "<>", ErrAdiReaderMalformedADI},
		{"Empty tag with text", "<>fake", ErrAdiReaderMalformedADI},
		{"Empty tag with colon", "<:>fake", ErrAdiReaderMalformedADI},
		{"Empty tag with double colon", "<::>fake", ErrAdiReaderMalformedADI},
		{"Empty tag with triple colon", "<:::>fake", ErrAdiReaderMalformedADI},
		{"Empty tag with quad colon", "<::::>fake", ErrAdiReaderMalformedADI},
		{"tag open and close", "<>", ErrAdiReaderMalformedADI},
		{"tag open and close with colon", "<:>", ErrAdiReaderMalformedADI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			p := NewADIRecordReader(strings.NewReader(tt.data), false)

			// Act
			qso, _, err := p.Next()
			if tt.expectedErr != err {
				t.Error("Expected non-EOF error, got EOF")
			}

			if qso != nil {
				t.Errorf("Expected nil record, got %v", qso)
			}
		})
	}
}

func TestADIRecordReaderParseSingleRecord(t *testing.T) {
	tests := []struct {
		name           string
		adifSource     string
		fieldName      string
		fieldData      string
		isHeaderRecord bool
	}{
		{"Header record", "<progRamid:4>MonoLog<EOH>", "PROGRAMID", "Mono", true},
		{"Header record", "<progRamid:4>MonoLog<EoH>", "PROGRAMID", "Mono", true},
		{"Short Record", "<WeB:1>X<Eor>", "WEB", "X", false},
		{"Short Record with type", "<WeB:1:s>X<Eor>", "WEB", "X", false},
		{"Zero length data", "<APP_MY_APP:0>\r\n<EOR>", "APP_MY_APP", "", false},
		{"Single char data", "<APP_MY_APP:1>x <EOR>", "APP_MY_APP", "x", false},
		{"Basic TIME_ON", "<TIME_ON:6>161819<EOR>", "TIME_ON", "161819", false},
		{"TIME_ON with type", "<TIME_ON:6:T>161819<EOR>", "TIME_ON", "161819", false},
		{"Mixed case TIME_ON", "<TiMe_ON:6>161819<EOR>", "TIME_ON", "161819", false},
		{"TIME_ON with type and space", "<TIME_ON:6:T>161819 <EOR>", "TIME_ON", "161819", false},
		{"Mixed case with space", "<TiMe_ON:6>161819 <EOR>", "TIME_ON", "161819", false},
		{"Leading space with type", " <TIME_ON:6:T>161819 <EOR>", "TIME_ON", "161819", false},
		{"Leading space mixed case", " <TiMe_ON:6>161819 <EOR>", "TIME_ON", "161819", false},
		{"Leading space with type no end space", " <TIME_ON:6:T>161819<EOR>", "TIME_ON", "161819", false},
		{"Leading space mixed case no end space", " <TiME_ON:6>161819<EOR>", "TIME_ON", "161819", false},
		{"Extra brackets", "><TiME_ON:6>161819<EOR>", "TIME_ON", "161819", false},
		{"Long Field Name", "<App_K9CTS_012345678901234567890:5>W9PVA<EOR>", "APP_K9CTS_012345678901234567890", "W9PVA", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use small buffer to test ReadSlice / ErrBufferFull handling
			br := bufio.NewReaderSize(strings.NewReader(tt.adifSource), 16)
			p := NewADIRecordReader(br, false)

			qso, isHeader, err := p.Next()
			if err != nil {
				t.Fatal(err)
			}

			if qso.Get(adifield.New(tt.fieldName)) != tt.fieldData {
				t.Errorf("Expected %s field to be %s, got %s", tt.fieldName, tt.fieldData, qso.Get(adifield.New(tt.fieldName)))
			}

			if isHeader != tt.isHeaderRecord {
				t.Errorf("Expected header record status %v, got %v", tt.isHeaderRecord, isHeader)
			}
		})
	}
}

func TestADIRecordReaderParseSkipHeader(t *testing.T) {
	// Arrange
	adif := "<PROGRAMID:7>MonoLog<EOH>\n<COMMENT:4>GOOD<EOR>"
	p := NewADIRecordReader(strings.NewReader(adif), true)

	// Act & Assert
	record, _, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if record.Get(adifield.PROGRAMID) != "" {
		t.Errorf("Expected empty PROGRAMID, got %s", record.Get(adifield.PROGRAMID))
	}
	if record.Get(adifield.COMMENT) != "GOOD" {
		t.Errorf("Expected COMMENT 'GOOD', got %s", record.Get(adifield.COMMENT))
	}

	recordTwo, _, errTwo := p.Next()
	if errTwo != io.EOF {
		t.Errorf("Expected EOF error, got %v", errTwo)
	}
	if recordTwo != nil {
		t.Errorf("Expected nil record, got %v", recordTwo)
	}
}

func TestADIRecordReaderParseLongFieldName(t *testing.T) {
	const len int = 2000
	fieldName := "app_k9cts_" + strings.Repeat("X", len)

	// Arrange
	adif := fmt.Sprintf("<%s:4>TEST<eor>", fieldName)
	p := NewADIRecordReader(strings.NewReader(adif), false)

	// Act
	record, _, err := p.Next()
	_, _, _ = p.Next()

	// Assert
	if err != nil {
		t.Fatal(err)
	}
	if record.Get(adifield.New(fieldName)) != "TEST" {
		t.Errorf("Expected %s field to be TEST, got %s", fieldName, record.Get(adifield.New(fieldName)))
	}
}

func TestADIRecordReaderParseLargeData(t *testing.T) {
	// Arrange
	p := NewADIRecordReader(strings.NewReader("<COMMENT:1000002>0"+strings.Repeat("1", 1_000_000)+"01<EOR>"), false)

	// Act
	record, _, err := p.Next() // Force the buffer to be resized to accommodate the large value
	if err != nil {
		t.Fatal(err)
	}
	if record.Get(adifield.COMMENT) != "0"+strings.Repeat("1", 1_000_000)+"0" {
		t.Errorf("Expected %s, got %s", strings.Repeat("1", 1_000_000), record.Get(adifield.COMMENT))
	}
}

func TestADIRecordReaderReadDataSpecifierVolatileRetunsError(t *testing.T) {
	mockReader := &mockFailReader{
		maxBytes:    5,
		backingData: []byte("<COMMENT:10>" + strings.Repeat("1", 10) + "<EOR>"),
	}
	rdr := NewADIRecordReader(mockReader, false)
	_, _, err := rdr.Next()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
