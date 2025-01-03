package adif

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/hamradiolog-net/adif-spec/src/pkg/adifield"
	"github.com/stretchr/testify/assert"
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
			assert.Nil(t, err)
			defer reader.Close()

			content, err := io.ReadAll(reader)
			assert.Nil(t, err)

			p := NewADIParser(strings.NewReader(string(content)), true)

			count := 0
			for {
				qso, _, err := p.Parse()
				if err == io.EOF {
					break
				}
				assert.NotNil(t, qso)
				assert.Nil(t, err)
				count++
			}

			assert.Equal(t, expectedCount, count, "Record count mismatch")
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

		{false, 1, "Record Missing EOR", "<CaLL:5>W9PVA"},
		{false, 1, "Complete Record with EOR", "<CaLL:5>W9PVA<EOR>"},
		{false, 1, "Leading space", " <CaLL:5>W9PVA<EOR>"},
		{false, 1, "Extra character", "<Call:5>W9PVAn<EOR>"},
		{false, 1, "Extra characters around EOR", "<Call:5>W9PVAa<EOR>b"},

		{true, 2, "With header", "<PROGRAMID:4>TEST<EOH><Call:5>W9PVAa<EOR>b"},
		{true, 2, "With header and extra chars", "<PROGRAMID:4>TESTing<EOH><Call:5>W9PVAa<EOR>b"},
		{true, 2, "With header, header preamble and extra chars", "preamble\n<PROGRAMID:4>TESTing<EOH><Call:5>W9PVAa<EOR>b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewADIParser(strings.NewReader(tt.data), false)

			records := make([]*Record, 0, 10000)
			for {
				record, _, err := p.Parse()
				if err == io.EOF {
					break
				}
				assert.Nil(t, err)
				records = append(records, record)
			}

			assert.Equal(t, tt.recordCount, len(records))

			if tt.recordCount == 0 {
				return
			}

			var index = 0
			if tt.hasHeader {
				assert.Equal(t, "TEST", records[0].Get(adifield.PROGRAMID))
				index++
			}

			assert.Equal(t, "W9PVA", records[index].Get(adifield.CALL))
		})
	}
}

func TestParseWithNumbersInFieldName(t *testing.T) {
	// Arrange
	raw := "<APP_LoTW_2xQSL:1>Y"
	p := NewADIParser(strings.NewReader(raw), false)

	// Act
	qso, bytesRead, err := p.Parse()

	// Assert
	assert.Nil(t, err)
	val := qso.Get("APP_LOTW_2XQSL")
	assert.Equal(t, "Y", val)
	assert.Equal(t, int64(len(raw)), bytesRead)
}

func TestParseWithMissingLengthField(t *testing.T) {
	// Arrange
	raw := "<APP_LoTW_EOF>Y" // n.b. 'Y' is NOT part of the data. It is a comment...
	p := NewADIParser(strings.NewReader(raw), false)

	// Act
	qso, bytesRead, err := p.Parse()

	// Assert
	assert.Equal(t, io.EOF, err)
	val := qso.Get("APP_LOTW_EOF")
	assert.Equal(t, "", val)
	assert.Equal(t, int64(len(raw)), bytesRead)
}

func TestParseNoRecords(t *testing.T) {
	// Arrange
	tests := []struct {
		name                string
		data                string
		isNonEOFErrExpected bool // EOF means success, non-EOF means the parser rejected the input as malformed.
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
		{"EOR after EOH", "<EOR><EOH>", true},
		{"Duplicate EOH", "<EOH><EOH>", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			p := NewADIParser(strings.NewReader(tt.data), false)

			// Act
			qso, _, err := p.Parse()
			assert.Equal(t, 0, len(qso.Fields))

			// Assert
			if tt.isNonEOFErrExpected {
				assert.NotNil(t, err)
				assert.NotEqual(t, io.EOF, err)
			} else {
				assert.Equal(t, io.EOF, err)
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
		{"Header record", "<progRamid:4>MonoLog", "PROGRAMID", "Mono", true, false},
		{"Zero length data", "<APP_MY_APP:0>\r\n", "APP_MY_APP", "", false, true},
		{"Single char data", "<APP_MY_APP:1>x ", "APP_MY_APP", "x", false, false},
		{"Basic TIME_ON", "<TIME_ON:6>161819", "TIME_ON", "161819", false, false},
		{"TIME_ON with type", "<TIME_ON:6:Time>161819", "TIME_ON", "161819", false, false},
		{"Mixed case TIME_ON", "<TiMe_ON:6>161819", "TIME_ON", "161819", false, false},
		{"TIME_ON with type and space", "<TIME_ON:6:Time>161819 ", "TIME_ON", "161819", false, false},
		{"Mixed case with space", "<TiMe_ON:6>161819 ", "TIME_ON", "161819", false, false},
		{"Leading space with type", " <TIME_ON:6:Time>161819 ", "TIME_ON", "161819", false, false},
		{"Leading space mixed case", " <TiMe_ON:6>161819 ", "TIME_ON", "161819", false, false},
		{"Leading space with type no end space", " <TIME_ON:6:Time>161819", "TIME_ON", "161819", false, false},
		{"Leading space mixed case no end space", " <TiME_ON:6>161819", "TIME_ON", "161819", false, false},
		{"Extra brackets", "><TiME_ON:6>161819>", "TIME_ON", "161819", false, false},
		{"Long Field Name", "<App_K9CTS_012345678901234567890:5>W9PVA", "APP_K9CTS_012345678901234567890", "W9PVA", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use small buffer to test ReadSlice / ErrBufferFull handling
			br := bufio.NewReaderSize(strings.NewReader(tt.adifSource), 16)
			p := NewADIParser(br, false)

			qso, bytesRead, err := p.Parse()
			if tt.isExpectEOF {
				assert.Equal(t, io.EOF, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.fieldData, qso.Get(adifield.Field(tt.fieldName)))

			isHeader, _ := qso.isHeaderRecord()
			assert.Equal(t, tt.isHeaderRecord, isHeader)
			assert.Equal(t, int64(len(tt.adifSource)), bytesRead)
		})
	}
}

func TestParseSkipHeader(t *testing.T) {
	// Arrange
	adif := "<PROGRAMID:7>MonoLog<EOH>\n<COMMENT:4>GOOD<EOR>"
	p := NewADIParser(strings.NewReader(adif), true)

	// Act & Assert
	record, _, err := p.Parse()
	assert.Nil(t, err)
	assert.Equal(t, "", record.Get("PROGRAMID"))
	assert.Equal(t, "GOOD", record.Get("COMMENT"))

	recordTwo, _, errTwo := p.Parse()
	assert.Equal(t, io.EOF, errTwo)
	assert.Equal(t, 0, len(recordTwo.Fields))
}

func TestParseLongFieldName(t *testing.T) {
	const len int = 2000
	fieldName := "APP_K9CTS_" + strings.Repeat("X", len)

	// Arrange
	adif := fmt.Sprintf("<%s:4>TEST", fieldName)
	p := NewADIParser(strings.NewReader(adif), false)

	// Act
	record, _, err := p.Parse()
	_, _, _ = p.Parse()

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, "TEST", record.Get(adifield.Field(fieldName)))
}

func TestParseLargeData(t *testing.T) {
	// Arrange
	p := NewADIParser(strings.NewReader("<COMMENT:1000002>0"+strings.Repeat("1", 1_000_000)+"01"), false)

	// Act
	record, _, err := p.Parse() // Force the buffer to be resized to accommodate the large value
	_, _, _ = p.Parse()         // Force the buffer to be resized back to "normal"

	// Assert
	assert.Nil(t, err)
	if record.Get(adifield.COMMENT) != "0"+strings.Repeat("1", 1_000_000)+"0" {
		t.Errorf("Expected %s, got %s", strings.Repeat("1", 1_000_000), record.Get(adifield.COMMENT))
	}
}

func TestParseLargeDataTooBigShouldReturnErr(t *testing.T) {
	// Arrange
	p := NewADIParser(strings.NewReader("<COMMENT:10000002>0"+strings.Repeat("1", 10_000_000)+"01"), false)

	// Act
	record, _, err := p.Parse() // Force the buffer to be resized to accommodate the large value
	_, _, _ = p.Parse()         // Force the buffer to be resized back to "normal"

	// Assert
	assert.Equal(t, ErrInvalidFieldLength, err)
	assert.Empty(t, record.Get(adifield.COMMENT))
}
