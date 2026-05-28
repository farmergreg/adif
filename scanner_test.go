package adif

import (
	"bufio"
	"embed"
	"fmt"
	"strings"
	"testing"

	"github.com/farmergreg/spec/v6/adifield"
)

//go:embed testdata/*.adi
var testFileFS embed.FS

func TestScannerVerifyRecordCount(t *testing.T) {
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
			f, err := testFileFS.Open("testdata/" + filename)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			count := 0
			s := NewScanner(f)
			for s.Scan() {
				if !s.IsHeader() {
					count++
				}
			}
			if err := s.Err(); err != nil {
				t.Fatal(err)
			}

			if count != expectedCount {
				t.Errorf("record count: got %d, want %d", count, expectedCount)
			}
		})
	}
}

func TestScannerParseBasicFunctionality(t *testing.T) {
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
			s := NewScanner(strings.NewReader(tt.data))
			records := make([]Record, 0)
			foundHeader := false
			for s.Scan() {
				foundHeader = foundHeader || s.IsHeader()
				records = append(records, s.Record())
			}
			if err := s.Err(); err != nil {
				t.Fatal(err)
			}

			if len(records) != tt.recordCount {
				t.Errorf("record count: got %d, want %d", len(records), tt.recordCount)
			}
			if tt.recordCount == 0 {
				return
			}

			index := 0
			if tt.hasHeader {
				if !foundHeader {
					t.Error("expected a header record")
				}
				if records[0][adifield.PROGRAMID] != "TEST" {
					t.Errorf("header PROGRAMID: got %q, want %q", records[0][adifield.PROGRAMID], "TEST")
				}
				index++
			}
			if records[index][adifield.CALL] != "W9PVA" {
				t.Errorf("CALL: got %q, want %q", records[index][adifield.CALL], "W9PVA")
			}
		})
	}
}

func TestScannerParseEOREOH(t *testing.T) {
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
			s := NewScanner(strings.NewReader(tt.data))
			count := 0
			for s.Scan() {
				count++
			}
			if err := s.Err(); err != nil {
				t.Fatal(err)
			}
			if count != tt.expected {
				t.Errorf("record count: got %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestScannerParseLoTWEOF(t *testing.T) {
	raw := "<" + string(adifield.APP_LOTW_EOF) + ">"
	s := NewScanner(strings.NewReader(raw))
	if s.Scan() {
		t.Error("expected Scan to return false")
	}
	if s.Err() != nil {
		t.Errorf("expected nil error (EOF), got %v", s.Err())
	}
	if s.Record() != nil {
		t.Errorf("expected nil record, got %v", s.Record())
	}
}

func TestScannerParseWithMissingEOH(t *testing.T) {
	raw := "<ADIF_VER:5>3.1.5<eor>"
	s := NewScanner(strings.NewReader(raw))
	if !s.Scan() {
		t.Fatal("expected a record")
	}
	if s.Record()[adifield.ADIF_VER] != "3.1.5" {
		t.Errorf("ADIF_VER: got %q, want %q", s.Record()[adifield.ADIF_VER], "3.1.5")
	}
	if s.Scan() {
		t.Error("expected no second record")
	}
	if s.Err() != nil {
		t.Errorf("expected nil error (EOF), got %v", s.Err())
	}
}

func TestScannerParseWithNumbersInFieldName(t *testing.T) {
	raw := "<APP_LoTW_2xQSL:1>Y<EOR>"
	s := NewScanner(strings.NewReader(raw))
	if !s.Scan() {
		t.Fatal("expected a record")
	}
	val := s.Record()[adifield.New("app_lotw_2xqsl")]
	if val != "Y" {
		t.Errorf("got %q, want %q", val, "Y")
	}
}

func TestScannerParseNoRecords(t *testing.T) {
	tests := []struct {
		name        string
		data        string
		expectedErr error // nil means normal EOF; non-nil means a parse error
	}{
		{"Invalid Length", "<APP_WAAT:fake>", ErrAdiReaderMalformedADI},
		{"Empty string", "", nil},
		{"Single space", " ", nil},
		{"Single colon", ":", nil},
		{"Double colon", "::", nil},
		{"Plain text", "no adif here...", nil},
		{"tag close", ">", nil},
		{"tag open", "<", ErrAdiReaderMalformedADI},
		{"Random text with tag", "< some random text", ErrAdiReaderMalformedADI},
		{"Math expression 1", " 3 < 4 ", ErrAdiReaderMalformedADI},
		{"Math expression 2", " 3 > 4 ", nil},
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
			s := NewScanner(strings.NewReader(tt.data))
			if s.Scan() {
				t.Fatal("expected no records")
			}
			if s.Err() != tt.expectedErr {
				t.Errorf("Err(): got %v, want %v", s.Err(), tt.expectedErr)
			}
			if s.Record() != nil {
				t.Errorf("expected nil record, got %v", s.Record())
			}
		})
	}
}

func TestScannerParseSingleRecord(t *testing.T) {
	tests := []struct {
		name           string
		adifSource     string
		fieldName      string
		fieldData      string
		isHeaderRecord bool
	}{
		{"Header record", "<progRamid:4>MonoLog<EOH>", "PROGRAMID", "Mono", true},
		{"Header record lower EOH", "<progRamid:4>MonoLog<EoH>", "PROGRAMID", "Mono", true},
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
			// Use a small buffer to exercise ReadSlice / ErrBufferFull handling.
			br := bufio.NewReaderSize(strings.NewReader(tt.adifSource), 16)
			s := NewScanner(br)
			if !s.Scan() {
				t.Fatalf("expected a record; Err=%v", s.Err())
			}

			got := s.Record()[adifield.New(tt.fieldName)]
			if got != tt.fieldData {
				t.Errorf("%s: got %q, want %q", tt.fieldName, got, tt.fieldData)
			}
			if s.IsHeader() != tt.isHeaderRecord {
				t.Errorf("IsHeader: got %v, want %v", s.IsHeader(), tt.isHeaderRecord)
			}
		})
	}
}

func TestScannerSkipHeader(t *testing.T) {
	adi := "<PROGRAMID:7>MonoLog<EOH>\n<COMMENT:4>GOOD<EOR>"
	s := NewScanner(strings.NewReader(adi))

	// Callers skip headers by checking IsHeader and continuing.
	foundQSO := false
	for s.Scan() {
		if s.IsHeader() {
			continue
		}
		foundQSO = true
		if s.Record()[adifield.PROGRAMID] != "" {
			t.Errorf("expected empty PROGRAMID on QSO record")
		}
		if s.Record()[adifield.COMMENT] != "GOOD" {
			t.Errorf("COMMENT: got %q, want %q", s.Record()[adifield.COMMENT], "GOOD")
		}
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	if !foundQSO {
		t.Error("expected at least one QSO record")
	}
}

func TestScannerParseLongFieldName(t *testing.T) {
	const nameLen = 2000
	fieldName := "app_k9cts_" + strings.Repeat("X", nameLen)
	adi := fmt.Sprintf("<%s:4>TEST<eor>", fieldName)

	s := NewScanner(strings.NewReader(adi))
	if !s.Scan() {
		t.Fatalf("expected a record; Err=%v", s.Err())
	}
	got := s.Record()[adifield.New(fieldName)]
	if got != "TEST" {
		t.Errorf("got %q, want %q", got, "TEST")
	}
}

func TestScannerParseLargeData(t *testing.T) {
	large := "0" + strings.Repeat("1", 1_000_000) + "0"
	adi := fmt.Sprintf("<COMMENT:%d>%s<EOR>", len(large), large)

	s := NewScanner(strings.NewReader(adi))
	if !s.Scan() {
		t.Fatalf("expected a record; Err=%v", s.Err())
	}
	got := s.Record()[adifield.COMMENT]
	if got != large {
		t.Errorf("large value mismatch (length got=%d, want=%d)", len(got), len(large))
	}
}

func TestScannerTooManyUniqueFields(t *testing.T) {
	var sb strings.Builder
	for i := range 1026 {
		fmt.Fprintf(&sb, "<APP_TEST_%04d:1>X", i)
	}
	sb.WriteString("<EOR>")

	s := NewScanner(strings.NewReader(sb.String()))
	if s.Scan() {
		t.Error("expected Scan to return false")
	}
	if s.Err() != ErrAdiReaderTooManyUniqueFields {
		t.Errorf("Err(): got %v, want %v", s.Err(), ErrAdiReaderTooManyUniqueFields)
	}
	if s.Record() != nil {
		t.Errorf("expected nil record, got %v", s.Record())
	}
}

func TestScannerReadDataSpecifierVolatile_IOError(t *testing.T) {
	mock := &mockFailReader{
		maxBytes:    5,
		backingData: []byte("<COMMENT:10>" + strings.Repeat("1", 10) + "<EOR>"),
	}
	s := NewScanner(mock)
	if s.Scan() {
		t.Error("expected Scan to return false")
	}
	if s.Err() == nil {
		t.Error("expected an error, got nil")
	}
}
